package config

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"os"
	"reflect"

	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/tests/util"
	_ "github.com/go-sql-driver/mysql"
)

type (
	ExecCase struct {
		Description    string           `json:"description"`    // 可选的，为测试场景提供描述
		SetUp          []EnvironmentSQL `json:"setUp"`          // 测试环境的设置 默认都是根据slice名下的主库执行
		GaeaActions    []GaeaAction     `json:"gaeaActions"`    // 中间件上执行的操作
		MasterCheckSQL []DBAction       `json:"masterCheckSQL"` // 主库执行的操作
		SlaveCheckSQL  []DBAction       `json:"slaveCheckSQL"`  // 从库执行的操作
		TearDown       []EnvironmentSQL `json:"tearDown"`       // 测试环境的清理 默认都是根据slice名下的主库执行
	}

	EnvironmentSQL struct {
		Description string `json:"description"` // 可选的，为测试场景提供描述
		Slice       string `json:"slice"`
		SQL         string `json:"sql"`
	}

	GaeaAction struct {
		SQL      string      `json:"sql"`
		ExecType string      `json:"execType"`
		Expect   interface{} `json:"expect"`
	}

	DBAction struct {
		Slice    string      `json:"slice"`
		DB       string      `json:"db"`
		SQL      string      `json:"sql"`
		ExecType string      `json:"execType"`
		Expect   interface{} `json:"expect"`
	}

	MysqlCluster struct {
		Slices map[string]*Slice `json:"slices"`
	}

	Slice struct {
		Master *DB   `json:"master"`
		Slaves []*DB `json:"slaves"`
	}
	DB struct {
		User     string `yaml:"user" json:"user"`
		Password string `yaml:"password" json:"password"`
		Url      string `yaml:"url" json:"url"`
	}

	Plan struct {
		ExecCases []ExecCase `json:"execCases"`
	}
)

type (
	SliceConn struct {
		Master *sql.DB
		Slaves []*sql.DB
	}
)

func InitMysqlClusterConn(slices []*models.Slice) (map[string]*SliceConn, error) {
	conn := make(map[string]*SliceConn)
	for _, slice := range slices {
		master, err := initDB(slice.UserName, slice.Password, slice.Master)
		if err != nil {
			return nil, err
		}
		sliceConn := &SliceConn{
			Master: master,
			Slaves: []*sql.DB{},
		}

		for _, slave := range slice.Slaves {
			slaveDB, err := initDB(slice.UserName, slice.Password, slave)
			if err != nil {
				return nil, err
			}
			sliceConn.Slaves = append(sliceConn.Slaves, slaveDB)
		}
		conn[slice.Name] = sliceConn
	}
	return conn, nil
}

func InitGaeaConn(user *models.User, host string) (*sql.DB, error) {

	gaeaDB, err := initDB(user.UserName, user.Password, host)
	if err != nil {
		return nil, err
	}
	return gaeaDB, nil
}

func initDB(username, password, host string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("error on initializing database connection: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database connection error:: %s", err.Error())
	}
	return db, nil
}

func (e ExecCase) GaeaRun(gaeaDB *sql.DB) error {
	// 执行在Gaea中的操作
	for _, action := range e.GaeaActions {
		isSuccess, err := DBExecAndCompare(gaeaDB, action.ExecType, action.SQL, action.Expect)
		if err != nil || !isSuccess {
			return fmt.Errorf("[gaeaAction] failed to execute SQL statement %s: %v", action.SQL, err)
		}
	}
	return nil
}

func DBExecAndCompare(db *sql.DB, expectType string, execSql string, expect interface{}) (bool, error) {
	switch expectType {
	case "Query":
		if !util.IsSlice(expect) {
			return false, fmt.Errorf("expect is not slice")
		}
		if v, ok := expect.([]interface{}); ok {
			// Convert []interface{} to [][]string
			var values [][]string
			for _, row := range v {
				var rowValues []string
				if rowItems, ok := row.([]interface{}); ok {
					for _, item := range rowItems {
						if strValue, ok := item.(string); ok {
							rowValues = append(rowValues, strValue)
						}
					}
				}
				values = append(values, rowValues)
			}
			rows, err := db.Query(execSql)
			if err != nil {
				if err == sql.ErrNoRows && len(v) == 0 {
					return true, nil
				}
				return false, fmt.Errorf("db Exec Error %v", err)
			}
			res, err := util.GetDataFromRows(rows)
			if err != nil {
				return false, fmt.Errorf("get data from rows error:%v", err)
			}
			// res为空代表没有查到数据
			if (len(res) == 0 || res == nil) && len(v) == 0 {
				return true, nil
			}
			if !reflect.DeepEqual(values, res) {

				return false, fmt.Errorf("mismatch. Actual: %v, Expect: %v", res, values)
			}
		} else {
			return false, fmt.Errorf("expect data assertion failed")
		}
		return true, nil // 所有数据都匹配
	case "QueryRow":
		if !util.IsSlice(expect) {
			return false, fmt.Errorf("expect is not slice")
		}
		if v, ok := expect.([]interface{}); ok {
			// Convert []interface{} to []string
			values := []string{}
			for _, item := range v {
				if strValue, ok := item.(string); ok {
					values = append(values, strValue)
				}
			}
			row := db.QueryRow(execSql)
			res, err := util.GetDataFromRow(row, len(v))
			if err != nil {
				if err == sql.ErrNoRows && len(v) == 0 {
					return true, nil
				}
				return false, fmt.Errorf("get data from row error:%v", err)
			}
			if !reflect.DeepEqual(values, res) {
				return false, fmt.Errorf("mismatch. Actual: %v, Expect: %v", res, values)
			}
		} else {
			return false, fmt.Errorf("expect data assertion failed")
		}
		return true, nil // 所有数据都匹配
	case "Exec":
		result, err := db.Exec(execSql)
		if err != nil {
			return false, err
		}
		// Assuming expect is a map with possible keys "lastInsertId" and "rowsAffected"
		expectedResults, ok := expect.(map[string]int64)
		if !ok {
			return false, errors.New("expect format for Exec type is incorrect")
		}

		// Check if "lastInsertId" is in expect and compare it with the result
		if lastInsertId, exists := expectedResults["lastInsertId"]; exists {
			actualLastInsertId, err := result.LastInsertId()
			if err != nil {
				return false, err
			}
			if lastInsertId != actualLastInsertId {
				return false, nil
			}
		}
		// Check if "rowsAffected" is in expect and compare it with the result
		if rowsAffected, exists := expectedResults["rowsAffected"]; exists {
			actualRowsAffected, err := result.RowsAffected()
			if err != nil {
				return false, err
			}
			if rowsAffected != actualRowsAffected {
				return false, nil
			}
		}
		return true, nil
	case "Default":
		_, err := db.Exec(execSql)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, errors.New("unsupported expect type")
	}
}

func (e *ExecCase) GetSetUpSQL() []EnvironmentSQL {
	return e.SetUp
}

func (e *ExecCase) GetTearDownSQL() []EnvironmentSQL {
	return e.TearDown
}

func (e EnvironmentSQL) MasterRun(cluster map[string]*SliceConn) error {
	// 执行
	var masterDB *sql.DB
	if slice, ok := cluster[e.Slice]; ok {
		masterDB = slice.Master
	} else {
		return fmt.Errorf("failed to get master database connection")
	}
	if _, err := masterDB.Exec(e.SQL); err != nil {
		return fmt.Errorf("[environment set action] failed to execute sql: %s,error :%s", e.SQL, err)
	}
	return nil
}

func (e ExecCase) MasterRunAndCheck(cluster map[string]*SliceConn) error {
	// 对主库执行的操作进行检查
	for _, check := range e.MasterCheckSQL {
		var masterDB *sql.DB
		if slice, ok := cluster[check.Slice]; ok {
			masterDB = slice.Master
		} else {
			return fmt.Errorf("failed to get master database connection")
		}
		// 这里不能关闭
		// defer masterDB.Close()
		if len(check.DB) != 0 {
			res, err := masterDB.Exec(fmt.Sprintf("USE %s", check.DB))
			if err != nil {
				return fmt.Errorf("[checkAction] failed to use DB statement %s: %v", res, err)
			}
		}
		isSuccess, err := DBExecAndCompare(masterDB, check.ExecType, check.SQL, check.Expect)
		if err != nil {
			return fmt.Errorf("[checkAction] failed to execute SQL statement %s: %v", check.SQL, err)
		}
		if !isSuccess {
			return fmt.Errorf("[checkAction] SQL execution %s did not meet the expected result: %v", check.SQL, check.Expect)
		}
	}
	return nil
}

func NewPlan(path string) (*Plan, error) {
	var p Plan
	err := LoadJsonConfig(path, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Plan) GetExecCases() []ExecCase {
	return p.ExecCases
}

func LoadJsonConfig(path string, val interface{}) error {
	if err := validInputIsPtr(val); err != nil {
		return fmt.Errorf("valid conf failed:%v", err)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = json.NewDecoder(bufio.NewReader(f)).Decode(val); err != nil {
		return err
	}
	return nil
}

func validInputIsPtr(conf interface{}) error {
	tp := reflect.TypeOf(conf)
	if tp.Kind() != reflect.Ptr {
		return errors.New("conf should be pointer")
	}
	return nil
}

type PlanManager struct {
	PlanPath         string
	Plan             *Plan
	MysqlClusterConn map[string]*SliceConn
	GeaaDB           *sql.DB
}

type PlanManagerOption func(*PlanManager)

func NewPlanManager(o ...PlanManagerOption) *PlanManager {
	m := &PlanManager{}
	for _, v := range o {
		v(m)
	}
	return m
}

func WithPlanPath(path string) PlanManagerOption {
	return func(m *PlanManager) {
		m.PlanPath = path
	}
}

func WithGaeaConn(conn *sql.DB) PlanManagerOption {
	return func(m *PlanManager) {
		m.GeaaDB = conn
	}
}

func WithMysqlClusterConn(conn map[string]*SliceConn) PlanManagerOption {
	return func(m *PlanManager) {
		m.MysqlClusterConn = conn
	}
}

func (m *PlanManager) Init() error {
	err := m.LoadPlan()
	if err != nil {
		return err
	}
	return nil
}

func (m *PlanManager) GetMasterConnByName(sliceName string) (*sql.DB, error) {
	var master *sql.DB
	if slice, ok := m.MysqlClusterConn[sliceName]; ok {
		master = slice.Master
	} else {
		return nil, fmt.Errorf("failed to get master database connection")
	}
	return master, nil
}

func (m *PlanManager) GetSlaveConnByName(sliceName string) ([]*sql.DB, error) {
	var master []*sql.DB
	if slice, ok := m.MysqlClusterConn[sliceName]; ok {
		master = slice.Slaves
	} else {
		return []*sql.DB{}, fmt.Errorf("failed to get master database connection")
	}
	return master, nil
}

func (m *PlanManager) Run() error {
	if m.Plan == nil {
		return errors.New("plan is not loaded")
	}
	// Run set up actions
	for _, execCase := range m.Plan.GetExecCases() {
		for _, initCase := range execCase.GetSetUpSQL() {
			// Run Master actions
			if err := initCase.MasterRun(m.MysqlClusterConn); err != nil {
				return fmt.Errorf("master init action failed '%s': %v", initCase.Description, err)
			}
		}
		// Run Gaea actions
		if err := execCase.GaeaRun(m.GeaaDB); err != nil {
			return fmt.Errorf("gaea action failed for test '%s': %v", execCase.Description, err)
		}

		// Run master checks
		if err := execCase.MasterRunAndCheck(m.MysqlClusterConn); err != nil {
			return fmt.Errorf("master check failed for test '%s': %v", execCase.Description, err)
		}
		// Run tear down action
		for _, tearDownCase := range execCase.GetTearDownSQL() {
			// Run Master actions
			if err := tearDownCase.MasterRun(m.MysqlClusterConn); err != nil {
				return fmt.Errorf("master init action failed '%s': %v", tearDownCase.Description, err)
			}
		}

	}
	return nil
}

func (m *PlanManager) MysqlClusterConnClose() {
	for _, slice := range m.MysqlClusterConn {
		if slice.Master != nil {
			slice.Master.Close()
		}
		// 关闭 Slaves 连接
		for _, slave := range slice.Slaves {
			if slave != nil {
				slave.Close()
			}
		}
	}
}

func (i *PlanManager) LoadPlan() error {
	var plan *Plan
	err := LoadJsonConfig(i.PlanPath, &plan)
	if err != nil {
		return err
	}
	i.Plan = plan
	return nil
}

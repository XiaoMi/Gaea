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
		SQL string `json:"sql"`
	}

	DBAction struct {
		Slice  string `json:"slice"`
		SQL    string `json:"sql"`
		Expect string `json:"expect"`
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
	MysqlClusterConn struct {
		SlicesConn map[string]*SliceConn `json:"slices"`
	}

	SliceConn struct {
		Master *sql.DB
		Slaves []*sql.DB
	}
)

func InitMysqlClusterConn(slices []*models.Slice) (*MysqlClusterConn, error) {
	conn := &MysqlClusterConn{
		SlicesConn: make(map[string]*SliceConn),
	}

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
		conn.SlicesConn[slice.Name] = sliceConn
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

		if _, err := gaeaDB.Exec(action.SQL); err != nil {
			return fmt.Errorf("failed to execute GaeaAction: %s", err)
		}
	}
	return nil
}

func (e *ExecCase) GetSetUpSQL() []EnvironmentSQL {
	return e.SetUp
}

func (e *ExecCase) GetTearDownSQL() []EnvironmentSQL {
	return e.TearDown
}

func (e EnvironmentSQL) MasterRun(cluster *MysqlClusterConn) error {
	// 执行
	var masterDB *sql.DB
	if slice, ok := cluster.SlicesConn[e.Slice]; ok {
		masterDB = slice.Master
	} else {
		return fmt.Errorf("failed to get master database connection")
	}

	if _, err := masterDB.Exec(e.SQL); err != nil {
		return fmt.Errorf("failed to execute environment sql : %s", err)
	}
	return nil
}

func (e ExecCase) MasterRunAndCheck(cluster *MysqlClusterConn) error {
	// 对主库执行的操作进行检查
	for _, check := range e.MasterCheckSQL {
		var masterDB *sql.DB
		if slice, ok := cluster.SlicesConn[check.Slice]; ok {
			masterDB = slice.Master
		} else {
			return fmt.Errorf("failed to get master database connection")
		}
		// 这里不能关闭
		// defer masterDB.Close()
		row := masterDB.QueryRow(check.SQL)
		var result string
		if err := row.Scan(&result); err != nil {
			if err == sql.ErrNoRows && check.Expect == "" {
				continue
			} else {
				return fmt.Errorf("failed to query master: %s", err)
			}
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
	MysqlClusterConn *MysqlClusterConn
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

func WithMysqlClusterConn(conn *MysqlClusterConn) PlanManagerOption {
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
	if slice, ok := m.MysqlClusterConn.SlicesConn[sliceName]; ok {
		master = slice.Master
	} else {
		return nil, fmt.Errorf("failed to get master database connection")
	}
	return master, nil
}

func (m *PlanManager) GetSlaveConnByName(sliceName string) ([]*sql.DB, error) {
	var master []*sql.DB
	if slice, ok := m.MysqlClusterConn.SlicesConn[sliceName]; ok {
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
	for _, slice := range m.MysqlClusterConn.SlicesConn {
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

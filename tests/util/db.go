package util

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func GetDataFromRows(rows *sql.Rows) ([][]string, error) {
	pr := func(t interface{}) (r string) {
		r = "\\N"
		switch v := t.(type) {
		case *sql.NullBool:
			if v.Valid {
				r = strconv.FormatBool(v.Bool)
			}
		case *sql.NullString:
			if v.Valid {
				r = v.String
			}
		case *sql.NullInt64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Int64)
			}
		case *sql.NullFloat64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Float64)
			}
		case *time.Time:
			if v.Year() > 1900 {
				r = v.Format("2006-01-02 15:04:05")
			}
		default:
			r = fmt.Sprintf("%#v", t)
		}
		return
	}

	c, _ := rows.Columns()
	n := len(c)
	field := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		field = append(field, new(sql.NullString))
	}

	var converts [][]string

	for rows.Next() {
		if err := rows.Scan(field...); err != nil {
			return nil, err
		}
		row := make([]string, 0, n)
		for i := 0; i < n; i++ {
			col := pr(field[i])
			row = append(row, col)
		}
		converts = append(converts, row)
	}
	return converts, nil
}

func GetDataFromRow(row *sql.Row, numCols int) ([]string, error) {
	pr := func(t interface{}) (r string) {
		r = "\\N"
		switch v := t.(type) {
		case *sql.NullBool:
			if v.Valid {
				r = strconv.FormatBool(v.Bool)
			}
		case *sql.NullString:
			if v.Valid {
				r = v.String
			}
		case *sql.NullInt64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Int64)
			}
		case *sql.NullFloat64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Float64)
			}
		case *time.Time:
			if v.Year() > 1900 {
				r = v.Format("2006-01-02 15:04:05")
			}
		default:
			r = fmt.Sprintf("%#v", t)
		}
		return
	}

	// Note: The sql.Row does not have a method to get columns directly
	// You'll need to know the number of columns (numCols) beforehand
	field := make([]interface{}, 0, numCols)
	for i := 0; i < numCols; i++ {
		field = append(field, new(sql.NullString))
	}

	if err := row.Scan(field...); err != nil {
		if err == sql.ErrNoRows {
			// Return nil or []string{} as per your requirement
			return []string{}, err
		}
		return nil, err
	}
	result := make([]string, 0, numCols)
	for i := 0; i < numCols; i++ {
		col := pr(field[i])
		result = append(result, col)
	}

	// As the sql.Row doesn't have a .Columns() method like sql.Rows, you'll need to know the column names from elsewhere
	// In this example, I've just returned an empty slice for column names
	return result, nil
}

func IsSlice(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Slice
}

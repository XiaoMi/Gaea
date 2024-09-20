package server

import (
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 检测 continueConn.MoreRowsExist Fields 是否正常
func TestWriteOkResultStreamMoreRowsExistFields(t *testing.T) {
	mockey.PatchConvey("test", t, func() {
		testMoreRowExistStatus := []bool{true, true, false}
		testMoreRowExistIndex := 0
		mockey.Mock((*ClientConn).writeOK).Return(nil).Build()
		mockey.Mock((*ClientConn).writeOKResult).To(func(status uint16, moreRows bool, r *mysql.Result) error {
			r.Free()
			return nil
		}).Build()
		mockey.Mock((*ClientConn).writeRowsWithEOF).To(func(r *mysql.Result, moreRowsExists bool, status uint16) error {
			r.Free()
			return nil
		}).Build()
		mockey.Mock((*mysql.Result).BuildBinaryResultSet).Return(nil).Build()
		mockey.Mock((*backend.MockPooledConnect).MoreRowsExist).To(func() bool {
			if testMoreRowExistIndex >= len(testMoreRowExistStatus) {
				return false
			}
			ret := testMoreRowExistStatus[testMoreRowExistIndex]
			testMoreRowExistIndex += 1
			return ret
		}).Build()
		mockey.Mock((*backend.MockPooledConnect).MoreResultsExist).Return(false).Build()
		mockey.Mock((*backend.MockPooledConnect).FetchMoreRows).To(func(result *mysql.Result, maxRows int) error {
			// 检测 fields 长度是否正常
			assert.Equal(t, 1, len(result.Fields))
			return nil
		}).Build()

		rs := mysql.ResultPool.Get()
		rs.Resultset = &mysql.Resultset{
			Fields: []*mysql.Field{
				{},
			},
		}
		backendConn := &backend.MockPooledConnect{}
		c := ClientConn{}
		c.writeOKResultStream(0, rs, backendConn, 0, true)
	})
}

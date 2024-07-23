package mysql

import "testing"
import "github.com/stretchr/testify/assert"

// test result set
func TestGetResultSet(t *testing.T) {
	rs := ResultPool.Get()
	assert.Equal(t, rs.Resultset == nil, false)
	rs.Free()
	rs = ResultPool.GetWithoutResultSet()
	assert.Equal(t, rs.Resultset == nil, true)
	rs.Free()
	rss := make([]*Result, 0)
	for i := 0; i < 10; i++ {
		rs = ResultPool.Get()
		assert.Equal(t, rs.Resultset == nil, false)
		rss = append(rss, rs)
	}
	for i := 0; i < 10; i++ {
		rs = ResultPool.GetWithoutResultSet()
		assert.Equal(t, rs.Resultset == nil, true)
		rss = append(rss, rs)
	}
	for _, rs := range rss {
		rs.Free()
	}
	for i := 0; i < 10; i++ {
		rs = ResultPool.Get()
		assert.Equal(t, rs.Resultset == nil, false)
		rs.Free()
	}
	for i := 0; i < 10; i++ {
		rs = ResultPool.GetWithoutResultSet()
		assert.Equal(t, rs.Resultset == nil, true)
		rs.Free()
	}
}

package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitMarginComments(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantQuery    string
		wantComments MarginComments
	}{
		{
			name:      "test no comment",
			sql:       "select * from t",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: "",
			},
		},
		{
			name:      "test leading comment",
			sql:       "/*master*/ select * from t",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "/*master*/ ",
				Trailing: "",
			},
		},
		{
			name:      "test internal comment",
			sql:       "select /*master*/ * from t",
			wantQuery: "select /*master*/ * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: "",
			},
		},
		{
			name:      "test mycat hint comment",
			sql:       "select * from t /* !mycat:sql=select 1 from xm_order where order_id = 4230329282206593 */",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: " /* !mycat:sql=select 1 from xm_order where order_id = 4230329282206593 */",
			},
		},
		{
			name:      "test trailing comment",
			sql:       "select * from t /*master*/",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: " /*master*/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotComments := SplitMarginComments(tt.sql)
			assert.Equalf(t, tt.wantQuery, gotQuery, "SplitMarginComments(%v)", tt.sql)
			assert.Equalf(t, tt.wantComments, gotComments, "SplitMarginComments(%v)", tt.sql)
		})
	}
}

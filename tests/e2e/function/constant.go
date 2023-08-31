package function

import (
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/tests/config"
)

var (
	db    = "db_e2e_test"
	slice = []*models.Slice{
		{
			Name:     "slice-0",
			UserName: "superroot",
			Password: "superroot",
			Master:   "127.0.0.1:3319",
			Slaves: []string{
				"127.0.0.1:3329",
				"127.0.0.1:3339",
			},
			StatisticSlaves: []string{},
			Capacity:        12,
			MaxCapacity:     24,
			IdleTimeout:     60,
			InitConnect:     "",
		},
	}

	allowedDBS = map[string]bool{
		"db_e2e_test": true,
	}
	users = []*models.User{
		{
			UserName:      "superroot",
			Password:      "superroot",
			Namespace:     "test_namespace_e2e_test",
			RWFlag:        2,
			RWSplit:       1,
			OtherProperty: 0,
		},
	}
	defaultNamespace = config.NewNamespace(
		config.WithIsEncrypt(true),
		config.WithName("test_namespace_e2e_test"),
		config.WithAllowedDBS(allowedDBS),
		config.WithSlices(slice),
		config.WithUsers(users),
		config.WithDefaultSlice("slice-0"),
	)
)

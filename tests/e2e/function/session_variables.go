// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package function

import (
	"database/sql"
	"fmt"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

// This Go testing file is part of an end-to-end testing suite designed to validate the functionality of setting variables in a Gaea database environment.
// The primary focus is to ensure that session-specific settings such as transaction isolation levels, maximum execution time, and unique checks behave as expected within a database session,
// and to validate their correct implementation in the Gaea database middleware.
// Key Objectives of the Testing File:
// Session Variable Configuration: The tests aim to verify that various session variables can be set and retrieved correctly.
// These variables include settings for transaction isolation, maximum execution time, and checks for data uniqueness.
var _ = ginkgo.Describe("Test Gaea SET SESSION Variables", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	table := config.DefaultE2eTable
	slice := e2eMgr.NsSlices[config.SliceSingleTestMaster]
	var needCleanup bool
	var allowedSessionVariables map[string]string

	// This is run once before each test (`It` block) to ensure fresh setup
	ginkgo.BeforeEach(func() {
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)
	})

	ginkgo.Context("When session variables are correctly configured", func() {
		ginkgo.BeforeEach(func() {
			allowedSessionVariables = map[string]string{
				"transaction_isolation": "string",
				"max_execution_time":    "int",
				"unique_checks":         "bool",
				"max_heap_table_size":   "int",
			}

			// Namespace preparation and modification per test to ensure isolation
			initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
			util.ExpectNoError(err, "parse namespace template")
			initNs.AllowedSessionVariables = allowedSessionVariables
			err = e2eMgr.ModifyNamespace(initNs)
			util.ExpectNoError(err)
		})

		ginkgo.It("should validate correctly configured session variables that are present in variableVerifyFuncMap", func() {
			// Sub-scenario A: Variables are in the verification map
			scenarioATests := []struct {
				VariableName      string
				Type              string
				AlternativeValues []interface{}
			}{
				// session variable in variableVerifyFuncMap and Namespace
				{"transaction_isolation", "string", []interface{}{"READ-UNCOMMITTED", "REPEATABLE-READ"}},
				{"unique_checks", "bool", []interface{}{1, 0}},
				{"max_execution_time", "int", []interface{}{5000, 10000}},
			}

			for _, test := range scenarioATests {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)

				// Check current value to avoid testing with the default value
				getCurrentValueSQL := fmt.Sprintf("SELECT @@SESSION.%s", test.VariableName)
				row := gaeaConn.QueryRow(getCurrentValueSQL)
				defaultValue, err := scanVariableValue(row, test.Type)
				util.ExpectNoError(err)

				for _, setValue := range test.AlternativeValues {
					if setValue != defaultValue {
						// Attempt to set the New value for the session variable
						setSQL := generateSetSessionSQL(test.VariableName, test.Type, setValue)
						_, err = gaeaConn.Exec(setSQL)
						util.ExpectNoError(err)

						//Verify that the session variable's value has changed
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						actualValue, err := scanVariableValue(row, test.Type)
						util.ExpectNoError(err)
						util.ExpectEqual(actualValue, setValue)

						// Restore the session variable to its original default value
						restoreSQL := generateSetSessionSQL(test.VariableName, test.Type, defaultValue)
						_, err = gaeaConn.Exec(restoreSQL)
						util.ExpectNoError(err)

						// Verify that the session variable has been restored to its default value
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						restoredValue, err := scanVariableValue(row, test.Type)
						util.ExpectNoError(err)
						util.ExpectEqual(restoredValue, defaultValue)
						// Test next variable after successful change and restoration
						break
					}
				}
			}
		})

		ginkgo.It("should handle correctly configured session variables not present in variableVerifyFuncMap", func() {
			// Test handling when variables are correct but not in the verification map
			needCleanup = true // Set this flag only if the test modifies the namespace
			tests := []struct {
				VariableName      string
				Type              string
				AlternativeValues []interface{}
			}{
				// session variable not in variableVerifyFuncMap and Namespace
				{"max_heap_table_size", "int", []interface{}{16384, 16777216}},
			}

			for _, test := range tests {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)

				// Check current value to avoid testing with the default value
				getCurrentValueSQL := fmt.Sprintf("SELECT @@SESSION.%s", test.VariableName)
				row := gaeaConn.QueryRow(getCurrentValueSQL)
				defaultValue, err := scanVariableValue(row, test.Type)
				util.ExpectNoError(err)

				for _, setValue := range test.AlternativeValues {
					if setValue != defaultValue {
						// Attempt to set the New value for the session variable
						setSQL := generateSetSessionSQL(test.VariableName, test.Type, setValue)
						_, err = gaeaConn.Exec(setSQL)
						util.ExpectNoError(err)

						//Verify that the session variable's value has changed
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						actualValue, err := scanVariableValue(row, test.Type)
						util.ExpectNoError(err)
						util.ExpectEqual(actualValue, setValue)

						// Restore the session variable to its original default value
						restoreSQL := generateSetSessionSQL(test.VariableName, test.Type, defaultValue)
						_, err = gaeaConn.Exec(restoreSQL)
						util.ExpectNoError(err)

						// Verify that the session variable has been restored to its default value
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						restoredValue, err := scanVariableValue(row, test.Type)
						util.ExpectNoError(err)
						util.ExpectEqual(restoredValue, defaultValue)
						// Test next variable after successful change and restoration
						break
					}
				}
			}
		})
	})

	ginkgo.Context("When session variables are incorrectly configured", func() {
		ginkgo.BeforeEach(func() {
			allowedSessionVariables = map[string]string{
				// wrong type for session variable
				"transaction_isolation": "int",
				// correct type for session variable
				"max_execution_time": "int",
				// correct type for session variable
				"unique_checks": "bool",
				// wrong type for session variable
				"max_heap_table_size": "string",
			}

			// Namespace preparation and modification per test to ensure isolation
			initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
			util.ExpectNoError(err, "parse namespace template")
			initNs.AllowedSessionVariables = allowedSessionVariables
			err = e2eMgr.ModifyNamespace(initNs)
			util.ExpectNoError(err)
		})

		ginkgo.It("should handle incorrectly configured session variables present in variableVerifyFuncMap", func() {
			// Sub-scenario A: Variables are in the verification map but have incorrect types
			scenarioATests := []struct {
				VariableName      string
				Type              string
				AlternativeValues []interface{}
			}{
				// session variable in variableVerifyFuncMap but configured with incorrect types
				{"transaction_isolation", "string", []interface{}{"READ-UNCOMMITTED", "REPEATABLE-READ"}},
			}

			for _, test := range scenarioATests {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				// Check current value to avoid testing with the default value
				getCurrentValueSQL := fmt.Sprintf("SELECT @@SESSION.%s", test.VariableName)
				row := gaeaConn.QueryRow(getCurrentValueSQL)
				defaultValue, err := scanVariableValue(row, test.Type)
				util.ExpectNoError(err)

				for _, setValue := range test.AlternativeValues {
					if setValue != defaultValue {
						// Attempt to set the New value for the session variable
						setSQL := generateSetSessionSQL(test.VariableName, test.Type, setValue)
						_, err = gaeaConn.Exec(setSQL)
						util.ExpectError(err, "Expected error due to incorrect type for variable")

						// Verify that the session variable's value has not changed
						row := gaeaConn.QueryRow(getCurrentValueSQL)
						actualValue, err := scanVariableValue(row, "string")
						util.ExpectNoError(err)
						util.ExpectEqual(actualValue, defaultValue)
						break
					}
				}
			}
		})

		ginkgo.It("should handle incorrectly configured session variables not present in variableVerifyFuncMap", func() {
			// Sub-scenario B: Variables not in the verification map but have incorrect types
			needCleanup = true // Set this flag only if the test modifies the namespace
			tests := []struct {
				VariableName      string
				Type              string
				AlternativeValues []interface{}
			}{
				// session variable not in variableVerifyFuncMap but configured with incorrect types
				{"max_heap_table_size", "int", []interface{}{16384, 16777216}},
			}

			for _, test := range tests {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				// Check current value to avoid testing with the default value
				getCurrentValueSQL := fmt.Sprintf("SELECT @@SESSION.%s", test.VariableName)
				row := gaeaConn.QueryRow(getCurrentValueSQL)
				defaultValue, err := scanVariableValue(row, test.Type)
				util.ExpectNoError(err)

				for _, setValue := range test.AlternativeValues {
					if setValue != defaultValue {
						// Attempt to set the New value for the session variable
						setSQL := generateSetSessionSQL(test.VariableName, test.Type, setValue)
						_, err = gaeaConn.Exec(setSQL)
						util.ExpectNoError(err, "Expected error due to incorrect type for variable")

						// First Select
						row := gaeaConn.QueryRow(getCurrentValueSQL)
						_, err := scanVariableValue(row, "int")
						util.ExpectError(err)

						// Verify that the session variable's value has not changed
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						actualValue, err := scanVariableValue(row, "int")
						util.ExpectNoError(err)
						util.ExpectEqual(actualValue, defaultValue)
						break
					}
				}
			}
		})
	})

	ginkgo.Context("When session variables are not configured", func() {
		ginkgo.BeforeEach(func() {
			// Namespace preparation and modification per test to ensure isolation
			initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
			util.ExpectNoError(err, "parse namespace template")
			err = e2eMgr.ModifyNamespace(initNs)
			util.ExpectNoError(err)
		})

		ginkgo.It("should not allow session variables not configured but present in variableVerifyFuncMap", func() {
			// Sub-scenario A: Variables are in the verification map but not configured in allowedSessionVariables
			scenarioATests := []struct {
				VariableName      string
				Type              string
				AlternativeValues []interface{}
			}{
				// session variable in variableVerifyFuncMap but not configured in allowedSessionVariables
				{"transaction_isolation", "string", []interface{}{"READ-UNCOMMITTED", "REPEATABLE-READ"}},
			}

			for _, test := range scenarioATests {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				// Check current value to avoid testing with the default value
				getCurrentValueSQL := fmt.Sprintf("SELECT @@SESSION.%s", test.VariableName)
				row := gaeaConn.QueryRow(getCurrentValueSQL)
				defaultValue, err := scanVariableValue(row, test.Type)
				util.ExpectNoError(err)

				for _, setValue := range test.AlternativeValues {
					if setValue != defaultValue {
						// Attempt to set the New value for the session variable
						setSQL := generateSetSessionSQL(test.VariableName, test.Type, setValue)
						_, err = gaeaConn.Exec(setSQL)
						util.ExpectNoError(err, "Setting a session variable not configured in allowedSessionVariables should not throw an error")

						// Verify that the session variable's value has not changed
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						actualValue, err := scanVariableValue(row, "string")
						util.ExpectNoError(err)
						util.ExpectEqual(actualValue, defaultValue)
					}
				}
			}
		})

		ginkgo.It("should not allow session variables not configured and not present in variableVerifyFuncMap", func() {
			// Sub-scenario B: Variables not in the verification map and not configured in allowedSessionVariables
			needCleanup = true // Set this flag only if the test modifies the namespace
			tests := []struct {
				VariableName      string
				Type              string
				AlternativeValues []interface{}
			}{
				// session variable not in variableVerifyFuncMap and not configured in allowedSessionVariables
				{"max_heap_table_size", "int", []interface{}{16384, 16777216}},
			}

			for _, test := range tests {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				// Check current value to avoid testing with the default value
				getCurrentValueSQL := fmt.Sprintf("SELECT @@SESSION.%s", test.VariableName)
				row := gaeaConn.QueryRow(getCurrentValueSQL)
				defaultValue, err := scanVariableValue(row, test.Type)
				util.ExpectNoError(err)

				for _, setValue := range test.AlternativeValues {
					if setValue != defaultValue {
						// Attempt to set the New value for the session variable
						setSQL := generateSetSessionSQL(test.VariableName, test.Type, setValue)
						_, err = gaeaConn.Exec(setSQL)
						util.ExpectNoError(err, "Setting a session variable not configured in allowedSessionVariables should not throw an error")

						// Verify that the session variable's value has not changed
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						actualValue, err := scanVariableValue(row, "int")
						util.ExpectNoError(err)
						util.ExpectEqual(actualValue, defaultValue)

					}
				}
			}
		})

		ginkgo.It("should support lock wait timeout", func() {
			tests := []struct {
				VariableName      string
				Type              string
				AlternativeValues []interface{}
			}{
				{"lock_wait_timeout", "int", []interface{}{2, 3}},
			}

			for _, test := range tests {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				// Check current value to avoid testing with the default value
				getCurrentValueSQL := fmt.Sprintf("SELECT @@SESSION.%s", test.VariableName)
				row := gaeaConn.QueryRow(getCurrentValueSQL)
				defaultValue, err := scanVariableValue(row, test.Type)
				util.ExpectNoError(err)
				for _, setValue := range test.AlternativeValues {
					if setValue != defaultValue {
						// Attempt to set the New value for the session variable
						setSQL := generateSetSessionSQL(test.VariableName, test.Type, setValue)
						_, err = gaeaConn.Exec(setSQL)
						util.ExpectNoError(err, "Setting lock_wait_timeout should not throw an error")

						// Verify that the session variable's value has not changed
						row = gaeaConn.QueryRow(getCurrentValueSQL)
						actualValue, err := scanVariableValue(row, "int")
						util.ExpectNoError(err)
						util.ExpectEqual(actualValue, setValue)
					}
				}
			}
		})
	})

	ginkgo.AfterEach(func() {
		if needCleanup {
			e2eMgr.Clean()
			needCleanup = false // Reset flag after cleaning
		}
	})
})

func generateSetSessionSQL(variableName string, variableType string, value interface{}) string {
	if variableType == "string" {
		// Properly quote string values
		return fmt.Sprintf("SET SESSION %s = '%v';", variableName, value)
	}
	// Non-string values do not need quotes
	return fmt.Sprintf("SET SESSION %s = %v;", variableName, value)
}

func scanVariableValue(row *sql.Row, variableType string) (interface{}, error) {
	switch variableType {
	case "int", "bool":
		var intValue int
		err := row.Scan(&intValue)
		if err != nil {
			return nil, err
		}
		return intValue, nil
	case "string":
		var stringValue string
		err := row.Scan(&stringValue)
		if err != nil {
			return nil, err
		}
		return stringValue, nil
	default:
		return nil, fmt.Errorf("unsupported variable type: %s", variableType)
	}
}

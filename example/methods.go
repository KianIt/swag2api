package main

import (
	"errors"
	"net/http"

	s2aErrors "github.com/KianIt/swag2api/errors"
	"github.com/KianIt/swag2api/example/models"
)

type ThisPackageModel struct {
	Field1 string  `json:"field1"`
	Field2 int     `json:"field2"`
	Field3 float64 `json:"field3"`
	Field4 bool    `json:"field4"`
	Field5 []byte  `json:"field5"`
}

// method1 godoc
// @ID method1
// @Param pathString path string true " "
// @Param pathInt path int true " "
// @Param pathFloat64 path float64 true " "
// @Param pathBool path bool true " "
// @Param pathBytes path []byte true " "
// @Router /path-to-method1 [get]
func method1(pathString string, pathInt int, pathFloat64 float64, pathBool bool, pathBytes []byte) (result string, err error) {
	return "success", nil
}

// method2 godoc
// @ID method2
// @Param queryString query string true " "
// @Param queryInt query int true " "
// @Param queryFloat64 query float64 true " "
// @Param queryBool query bool true " "
// @Param queryBytes query []byte true " "
// @Router /path-to-method2 [post]
func method2(queryString string, queryInt int, queryFloat64 float64, queryBool bool, queryBytes []byte) (_ string, err error) {
	return "success", nil
}

// method3 godoc
// @ID method3
// @Param bodyString body string true " "
// @Param bodyInt body int true " "
// @Param bodyFloat64 body float64 true " "
// @Param bodyBool body bool true " "
// @Param bodyBytes body []byte true " "
// @Router /path-to-method3 [put]
func method3(bodyString string, bodyInt int, bodyFloat64 float64, bodyBool bool, bodyBytes []byte) (result string, _ error) {
	return "success", nil
}

// method4 godoc
// @ID method4
// @Param pathString path string true " "
// @Param queryString query string true " "
// @Param bodyString body string true " "
// @Router /path-to-method4 [delete]
func method4(pathString, queryString, bodyString string) (string, error) {
	return "success", nil
}

// method5 godoc
// @ID method5
// @Param bodyMap body map[string][]map[string]int true " "
// @Router /path-to-method5 [options]
func method5(bodyMap map[string][]map[string]int) (result string, Error error) {
	return "success", nil
}

// method6 godoc
// @ID method6
// @Param queryModel query ThisPackageModel true " "
// @Router /path-to-method6 [head]
func method6(field1 string, field2 int, field3 float64, field4 bool, field5 []byte) (code int, err error) {
	return http.StatusOK, nil
}

// method7 godoc
// @ID method7
// @Param bodyModel body ThisPackageModel true " "
// @Param bodyModelList body []ThisPackageModel true " "
// @Param bodyModelMap body map[string]models.AnotherPackageModel true " "
// @Router /path-to-method7 [patch]
func method7(bodyModel ThisPackageModel, bodyModelList []ThisPackageModel, bodyModelMap map[string]models.AnotherPackageModel) (code int, err error) {
	return http.StatusOK, nil
}

// method8 godoc
// @ID method8
// @Router /path-to-method8 [get]
func method8() (result string, err error) {
	return "success", s2aErrors.NotFound(nil)
}

// method9 godoc
// @ID method9
// @Router /path-to-method9 [get]
func method9() (result string, err error) {
	return "failed", s2aErrors.NotFound(errors.New("test error"))
}

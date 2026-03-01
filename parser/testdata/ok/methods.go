package ok

import (
	"net/http"

	models "github.com/KianIt/swag2api/parser/testdata/ok/models"
)

var handler http.HandlerFunc

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
// @Param pathInt query int true " "
// @Param pathFloat64 body float64 true " "
// @Router /path-to-method1 [get]
func method1(pathString string, pathInt int, pathFloat64 float64) (result string, err error) {
	return "success", nil
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

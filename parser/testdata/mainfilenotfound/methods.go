package ok

// method1 godoc
// @ID method1
// @Param pathString path string true " "
// @Param pathInt query int true " "
// @Param pathFloat64 body float64 true " "
// @Router /path-to-method1 [get]
func method1(pathString string, pathInt int, pathFloat64 float64) (result string, err error) {
	return "success", nil
}

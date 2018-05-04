//+build !debug

package onering

const DEBUG = false

func checkGetType(interface{}, string)  {}
func checkPutType(interface{}, string)  {}
func checkFuncType(interface{}, string) {}

func getCallerPath() string { return "" }

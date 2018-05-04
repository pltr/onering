package onering

//func ExamplePutType() {
//	var queue = New{Size:4}.SPSC()
//	queue.Put(5)
//	// Output:
//	//
//}
//
//func TestGetType(t *testing.T) {
//	var queue = New{Size:4}.SPSC()
//	queue.Put(5)
//	var i int
//	for queue.Get(&i) {
//		fmt.Println(i)
//	}
//}
//
//func TestFuncReturnType(t *testing.T) {
//	var queue = New{Size: 4}.SPSC()
//	queue.Close()
//	fn := func(s string, v int) string {
//		return ""
//	}
//	_ = fn
//	queue.Consume(fn)
//}

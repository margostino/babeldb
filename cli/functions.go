package cli

func createSource(params map[Param]interface{}) {
	name := params[sourceName].(string)
	println(name)
	//source := &db.Source{
	//	Name: name,
	//	Url:
	//}

}

package cli

import "github.com/margostino/babeldb/engine"

func createSource(engine *engine.Engine, params map[Param]interface{}) {
	name := params[sourceName].(string)
	url := params[sourceUrl].(string)
	cronExp := params[schedule].(string)
	engine.AddSource(name, url, cronExp)
}

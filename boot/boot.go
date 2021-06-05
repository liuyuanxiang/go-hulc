package boot

func (app *Application) Init() error {
	// 加载对应的配置文件内容
	app.Config.Load("app.yaml")
	return nil
}

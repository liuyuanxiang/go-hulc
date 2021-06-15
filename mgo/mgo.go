package mgo

import (
	"fmt"
	"time"

	"github.com/liuyuanxiang/go-hulc/boot"
	"gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session

// MgoSession 可以根据提供的 App 应用实例，自动获取其中加载的配置信息来返回对应的 Mongo 实例
func MgoSession(app *boot.Application) (*mgo.Session, error) {
	if mgoSession != nil {
		return mgoSession, nil
	}

	var err error
	mgoSession, err = newMgoSessionFromApp(app)
	if err != nil {
		return nil, err
	}
	return mgoSession, nil
}

// ReloadMgoSession 可以重新加载 Mongo 实例并返回
// 如果 App 中的 Mongo 配置信息变更，希望关闭旧的链接并建立新的链接返回时，可以调用该方法
// 该方法调用后，原本的链接将会失效不可用
func ReloadMgoSession(app *boot.Application) (*mgo.Session, error) {
	if mgoSession != nil {
		// 关闭旧链接
		mgoSession.Close()
	}

	var err error
	mgoSession, err = newMgoSessionFromApp(app)
	if err != nil {
		return nil, err
	}
	return mgoSession, nil
}

func newMgoSessionFromApp(app *boot.Application) (*mgo.Session, error) {
	host, _ := app.Config.GetDefault("mongo.host", "127.0.0.1").(string)
	port, _ := app.Config.GetDefault("mongo.port", 27017).(int64)
	username, _ := app.Config.GetDefault("mongo.username", "").(string)
	password, _ := app.Config.GetDefault("mongo.password", "").(string)
	timeout, _ := app.Config.GetDefault("mongo.timeout", 300).(int64)
	if host == "" || port == 0 {
		return nil, fmt.Errorf("MongoDB 所需 host port 配置缺失")
	}

	url := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	m, err := mgo.DialWithTimeout(url, time.Duration(timeout)*time.Microsecond)
	if err != nil {
		return nil, err
	}

	return m, nil
}

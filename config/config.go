package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	// 配置文件默认放置于项目下的一级目录 config 中
	defaultConfigPath = "./config/"
)

type Config struct {
	loadPath string
	isLoad   bool
	v        *viper.Viper
}

// NewConfig 返回一个配置管理实例
func NewConfig() *Config {
	return &Config{
		loadPath: defaultConfigPath,
		isLoad:   false,
		v:        viper.GetViper(),
	}
}

// SetConfigLoadPath 可以设置配置文件的加载路径
// 未通过调用该方法设置配置文件路径时，默认会从 项目/config/ 目录下进行文件读取
func SetConfigLoadPath(path string) {
	defaultConfigPath = path
}

// Load 根据默认或应用指定的路径加载对应的配置文件
func (c *Config) Load(file string) error {
	if !c.isLoad {
		c.v.SetConfigFile(c.loadPath + file)
		if err := c.v.ReadInConfig(); err != nil {
			return fmt.Errorf("配置文件 %s 加载失败 err: %v", file, err)
		}
		c.isLoad = true
	}
	return nil
}

// GetViper 返回底层 viper 的实例
func (c *Config) GetViper() *viper.Viper { return c.v }

// Get 根据提供的配置 Key 返回对应的配置内容
// 如果配置文件未加载或加载失败时，将统一返回 nil
func (c *Config) Get(key string) interface{} {
	if !c.isLoad {
		return nil
	}
	return c.v.Get(key)
}

// GetDefault 在找不到可用的配置内容时，可以自定义返回的默认值
func (c *Config) GetDefault(key string, dv interface{}) interface{} {
	if v := c.Get(key); v != nil {
		return v
	}
	return dv
}

// GetInt return a int
func (c *Config) GetInt(key string) int {
	if !c.isLoad {
		return 0
	}
	return c.v.GetInt(key)
}

func (c *Config) GetInt32(key string) int32 {
	if !c.isLoad {
		return 0
	}
	return c.v.GetInt32(key)
}

// GetInt64 return a int64
func (c *Config) GetInt64(key string) int64 {
	if !c.isLoad {
		return 0
	}
	return c.v.GetInt64(key)
}

// GetString return a string
func (c *Config) GetString(key string) string {
	if !c.isLoad {
		return ""
	}
	return c.v.GetString(key)
}

// GetStringMap return a map
func (c *Config) GetStringMap(key string) map[string]interface{} {
	if !c.isLoad {
		return make(map[string]interface{})
	}
	return c.v.GetStringMap(key)
}

// IsProdEnv 判断当前应用的运行环境是否为生产环境
// 根据配置文件中的 app.env 内容判断
func (c *Config) IsProdEnv() bool {
	return c.GetString("app.env") == "prod"
}

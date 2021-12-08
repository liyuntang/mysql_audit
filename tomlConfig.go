package tomlConfig

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"sync"
)

var (
	conf *AUDIT
	once sync.Once
)

func TomlConfig(configFile string) *AUDIT {
	// 检测配置文件是否存在
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		// 说明配置文件不存在，直接退出程序
		fmt.Println("sorry config file", configFile, "is not exist, err is", err)
		os.Exit(0)
	}
	// 说明文件存在，获取配置文件的绝对路径
	absPath, err1 := filepath.Abs(configFile)
	if err1 != nil {
		// 说明获取配置文件绝对路径失败，报错，退出程序
		fmt.Println("sorry get abs path of config file", configFile, "is bad, err is", err1)
		os.Exit(0)
	}
	// 说明获取配置文件绝对路径成功，使用单例模式读取配置
	once.Do(func() {
		_, err2 := toml.DecodeFile(absPath, &conf)
		if err2 != nil {
			// 说明解析配置文件失败
			fmt.Println("sorry, toml config file", configFile, "is bad, err is", err2)
			os.Exit(0)
		}
	})
	return conf
}
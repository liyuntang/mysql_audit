package tomlConfig

import (
	"time"
)

type AUDIT struct {
	System system
	Database database
}

type system struct {
	Port int	`toml:"port"`
	Retry int 	`toml:"retry"`
	Thread int	`toml:"thread"`
	IntervalTime time.Duration	`toml:"interval_time"`
	DataDir string		`toml:"data_dir"`
	LogFile string		`toml:"log_file"`
}

type database struct {
	User	string	`toml:"user"`
	Passwd	string	`toml:"passwd"`
	Address 	string	`toml:"address"`
	Port 	int		`toml:"port"`
	Charset	string	`toml:"charset"`
	Schema 	string	`toml:"schema"`

}
// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Address          string `config:"address"`
	Port             int    `config:"port"`
	Protocol         string `config:"protocol"`
	MaxMsgSize       int    `config:"max_msg_size"`
	DefaultEsLogType string `config:"default_es_log_type"`
}

var DefaultConfig = Config{
	Address:          "127.0.0.1",
	Port:             5000,
	Protocol:         "udp",
	MaxMsgSize:       4096,
	DefaultEsLogType: "simpleudpbeat",
}

package v2ray_ssrpanel_plugin

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	flag "github.com/spf13/pflag"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/main/confloader"
	json_reader "v2ray.com/ext/encoding/json"
	"v2ray.com/ext/tools/conf"
)

var (
	commandLine = flag.NewFlagSet(os.Args[0]+"-plugin-ssrpanel", flag.ContinueOnError)
	configFile = commandLine.String("config", "", "Config file for V2Ray.")
	test       = commandLine.Bool("test", false, "Test config file only, without launching V2Ray server.")
)

type UserConfig struct {
	InboundTag     string `json:"inboundTag"`
	Level          uint32 `json:"level"`
	AlterID        uint32 `json:"alterId"`
	SecurityStr    string `json:"securityConfig"`
	securityConfig *protocol.SecurityConfig
}

func (c *UserConfig) UnmarshalJSON(data []byte) error {
	type config UserConfig
	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// set default value
	if cfg.SecurityStr == "" {
		cfg.SecurityStr = "AUTO"
	}

	cfg.securityConfig = &protocol.SecurityConfig{
		Type: protocol.SecurityType(protocol.SecurityType_value[strings.ToUpper(cfg.SecurityStr)]),
	}
	*c = UserConfig(cfg)
	return nil
}

type Config struct {
	NodeID      uint         `json:"nodeId"`
	CheckRate   int          `json:"checkRate"`
	TrafficRate float64      `json:"trafficRate"`
	MySQL       *MySQLConfig `json:"mysql"`
	UserConfig  *UserConfig  `json:"user"`
	GRPCAddr    string       `json:"gRPCAddr"`
	v2rayConfig   *conf.Config
}

func getConfig() (*Config, error) {
	type config struct {
		*conf.Config
		SSRPanel *Config `json:"ssrpanel"`
	}

	configFile := getConfigFilePath()
	configInput, err := confloader.LoadConfig(configFile)
	if err != nil {
		return nil, errors.New("failed to load config: ", configFile).Base(err)
	}
	defer configInput.Close()

	cfg := &config{}
	if err = decodeCommentJSON(configInput, cfg); err != nil {
		return nil, err
	}
	if cfg.SSRPanel != nil {
		cfg.SSRPanel.v2rayConfig = cfg.Config
	}

	return cfg.SSRPanel, err
}

func getConfigFilePath() string {
	if len(*configFile) > 0 {
		return *configFile
	}

	if workingDir, err := os.Getwd(); err == nil {
		configFile := filepath.Join(workingDir, "config.json")
		if fileExists(configFile) {
			return configFile
		}
	}

	if configFile := platform.GetConfigurationPath(); fileExists(configFile) {
		return configFile
	}

	return ""
}

func decodeCommentJSON(reader io.Reader, i interface{}) error {
	jsonContent := bytes.NewBuffer(make([]byte, 0, 10240))
	jsonReader := io.TeeReader(&json_reader.Reader{
		Reader: reader,
	}, jsonContent)
	decoder := json.NewDecoder(jsonReader)
	return decoder.Decode(i)
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

package v2ray_ssrpanel_plugin

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"v2ray.com/core/common/errors"

	"v2ray.com/core/common/platform"
	"v2ray.com/core/main/confloader"
	json_reader "v2ray.com/ext/encoding/json"
)

var (
	commandLine = flag.NewFlagSet(os.Args[0]+"-plugin-ssrpanel", flag.ContinueOnError)

	configFile = commandLine.String("config", "", "Config file for V2Ray.")
	_          = commandLine.Bool("version", false, "Show current version of V2Ray.")
	test       = commandLine.Bool("test", false, "Test config file only, without launching V2Ray server.")
	_          = commandLine.String("format", "json", "Format of input file.")
	_          = commandLine.Bool("plugin", false, "True to load plugins.")
)

type Config struct {
	InboundTag  string  `json:"inboundTag"`
	NodeID      int     `json:"nodeId"`
	CheckRate   int     `json:"checkRate"`
	TrafficRate float64 `json:"trafficRate"`
	MySQL       struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"mysql"`
}

func testConfig() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	logConfig(cfg)

	return nil
}

func getConfig() (*Config, error) {
	configFile := getConfigFilePath()
	configInput, err := confloader.LoadConfig(configFile)
	if err != nil {
		return nil, errors.New("failed to load config: ", configFile).Base(err)
	}
	defer configInput.Close()

	plugins := &struct {
		Other struct {
			Plugins map[string]json.RawMessage `json:"plugins"`
		} `json:"other"`
	}{}

	if err = decodeCommentJSON(configInput, plugins); err != nil {
		return nil, err
	}

	rawCfg, ok := plugins.Other.Plugins["ssrpanel"]
	if !ok {
		return nil, errors.New("please add SSR Panel config")
	}

	cfg := &Config{}
	err = json.Unmarshal(rawCfg, cfg)

	return cfg, err
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

func logConfig(cfg *Config) {
	configContent, _ := json.Marshal(cfg)
	newError("got config: ", string(configContent)).AtInfo().WriteToLog()
}

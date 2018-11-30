package v2ray_ssrpanel_plugin

import (
	"fmt"
	"os"
	"time"
	"v2ray.com/core/common/errors"
)

func Run() {
	err := run()
	if err != nil {
		fatal(err)
	}
}

func run() error {
	commandLine.Parse(os.Args[1:])
	if *test {
		return testConfig()
	}

	globalCfg, err := getConfig()
	if err != nil {
		return err
	}

	db, err := NewMySQLConn(globalCfg.myPluginConfig.MySQL)
	if err != nil {
		return err
	}

	go func() {
		gRPCAddr := globalCfg.myPluginConfig.GRPCAddr
		gRPCConn, err := connectGRPC(gRPCAddr, 45*time.Second)
		if err != nil {
			fatal(fmt.Sprintf("connect to gRPC server \"%s\"", gRPCAddr), err)
		}
		newError(fmt.Sprintf("Connected gRPC server \"%s\"", gRPCAddr)).AtWarning().WriteToLog()

		p := NewPanel(gRPCConn, db, globalCfg)
		p.Start()
	}()

	return nil
}

func newError(values ...interface{}) *errors.Error {
	values = append([]interface{}{"SSRPanelPlugin: "}, values...)
	return errors.New(values...)
}

func fatal(values ...interface{}) {
	newError(values...).AtError().WriteToLog()
	// Wait log
	time.Sleep(1 * time.Second)
	os.Exit(-2)
}

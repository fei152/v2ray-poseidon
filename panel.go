package v2ray_ssrpanel_plugin

import (
	"fmt"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

type Panel struct {
	handlerServiceClient *HandlerServiceClient
	db                   *DB
	userModels           []UserModel
	globalConfig         *Config
}

func NewPanel(gRPCConn *grpc.ClientConn, db *DB, globalConfig *Config) *Panel {
	return &Panel{
		db:                   db,
		handlerServiceClient: NewHandlerServiceClient(gRPCConn, globalConfig.myPluginConfig.InboundTag),
		globalConfig:         globalConfig,
	}
}

func (p *Panel) Start() {
	doFunc := func() {
		if err := p.do(); err != nil {
			newError("panel#do").Base(err).AtError().WriteToLog()
		}
	}
	doFunc()

	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %ds", p.globalConfig.myPluginConfig.CheckRate), doFunc)
	c.Start()
	c.Run()
}

func (p *Panel) do() error {
	var addedUserCount, deletedUserCount int
	var uplinkTraffic, downlinkTraffic int64
	newError("start jobs").AtDebug().WriteToLog()
	defer func() {
		// todo
		newError(fmt.Sprintf("jobs info: addded %d users, deleteted %d users, downlink traffic %d KB, uplink traffic %d KB",
			addedUserCount, deletedUserCount, downlinkTraffic, uplinkTraffic)).AtInfo().WriteToLog()
	}()

	userModels, err := p.db.GetAllUsers()
	if err != nil {
		return err
	}

	// Calculate addition users
	addUserModels := make([]UserModel, 0)
	for _, userModel := range userModels {
		if inUserModels(&userModel, p.userModels) {
			continue
		}

		addUserModels = append(addUserModels, userModel)
	}

	// Calculate deletion users
	delUserModels := make([]UserModel, 0)
	for _, userModel := range p.userModels {
		if inUserModels(&userModel, userModels) {
			continue
		}

		delUserModels = append(delUserModels, userModel)
	}

	// Delete
	for _, userModel := range delUserModels {
		if i := findUserModelIndex(&userModel, p.userModels); i != -1 {
			p.userModels = append(p.userModels[:i], p.userModels[i+1:]...)
			p.handlerServiceClient.DelUser(userModel.Email)
		}
	}
	deletedUserCount = len(delUserModels)

	// Add
	p.userModels = append(p.userModels, addUserModels...)
	for _, userModel := range addUserModels {
		p.handlerServiceClient.AddUser(p.convertUser(userModel))
	}
	addedUserCount = len(addUserModels)

	return nil
}

func (p *Panel) convertUser(userModel UserModel) *protocol.User {
	userCfg := p.globalConfig.myPluginConfig.UserConfig
	return &protocol.User{
		Level: userCfg.Level,
		Email: userModel.Email,
		Account: serial.ToTypedMessage(&vmess.Account{
			Id:               userModel.VmessID,
			AlterId:          userCfg.AlterID,
			SecuritySettings: userCfg.securityConfig,
		}),
	}
}

func findUserModelIndex(u *UserModel, userModels []UserModel) int {
	for i, user := range userModels {
		if user.ID == u.ID {
			return i
		}
	}
	return -1
}

func inUserModels(u *UserModel, userModels []UserModel) bool {
	return findUserModelIndex(u, userModels) != -1
}

package v2ray_ssrpanel_plugin

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"v2ray.com/core/app/stats/command"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

type Panel struct {
	handlerServiceClient *HandlerServiceClient
	statsServiceClient *StatsServiceClient
	db                   *DB
	userModels           []UserModel
	globalConfig         *Config
}

func NewPanel(gRPCConn *grpc.ClientConn, db *DB, globalConfig *Config) *Panel {
	return &Panel{
		db:                   db,
		handlerServiceClient: NewHandlerServiceClient(gRPCConn, globalConfig.myPluginConfig.InboundTag),
		statsServiceClient: NewStatsServiceClient(gRPCConn),
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


func (p *Panel) getTraffic() (downlinkTotal uint64, uplinkTotal uint64, err error) {
	var stat *command.Stat
	for _, user := range p.userModels {
		stat, err = p.statsServiceClient.getUserDownlink(user.Email)
		if err != nil {
			return
		}
		downlink := uint64(stat.Value)

		stat, err = p.statsServiceClient.getUserUplink(user.Email)
		if err != nil {
			return
		}
		uplink := uint64(stat.Value)

		log := UserTrafficLog{
			UserID: user.ID,
			Uplink: uplink,
			Downlink:downlink,
			NodeID: p.globalConfig.myPluginConfig.NodeID,
			Rate: p.globalConfig.myPluginConfig.TrafficRate,
			Traffic: bytefmt.ByteSize(uplink + downlink),
		}
		if p.db.CreateUserTrafficLog(&log) == false {
			return 0, 0, errors.New("create user traffic log error")
		}

		downlinkTotal += downlink
		uplinkTotal += uplink
	}

	return
}

func (p *Panel) do() error {
	var addedUserCount, deletedUserCount int
	var uplinkTraffic, downlinkTraffic uint64
	newError("start jobs").AtDebug().WriteToLog()
	defer func() {
		newError(fmt.Sprintf("jobs info: addded %d users, deleteted %d users, downlink traffic %s, uplink traffic %s",
			addedUserCount, deletedUserCount, bytefmt.ByteSize(downlinkTraffic), bytefmt.ByteSize(uplinkTraffic))).AtInfo().WriteToLog()
	}()

	var err error
	if uplinkTraffic, downlinkTraffic, err = p.getTraffic(); err != nil {
		return errors.New("get traffic").Base(err)
	}

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
			if err = p.handlerServiceClient.DelUser(userModel.Email); err != nil {
				return err
			}
		}
	}
	deletedUserCount = len(delUserModels)

	// Add
	p.userModels = append(p.userModels, addUserModels...)
	for _, userModel := range addUserModels {
		if err = p.handlerServiceClient.AddUser(p.convertUser(userModel)); err != nil {
			return err
		}
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

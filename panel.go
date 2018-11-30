package v2ray_ssrpanel_plugin

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

type Panel struct {
	handlerServiceClient *HandlerServiceClient
	statsServiceClient   *StatsServiceClient
	db                   *DB
	userModels           []UserModel
	globalConfig         *Config
}

func NewPanel(gRPCConn *grpc.ClientConn, db *DB, globalConfig *Config) *Panel {
	return &Panel{
		db:                   db,
		handlerServiceClient: NewHandlerServiceClient(gRPCConn, globalConfig.myPluginConfig.UserConfig.InboundTag),
		statsServiceClient:   NewStatsServiceClient(gRPCConn),
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

func (p *Panel) do() (err error) {
	var addedUserCount, deletedUserCount int
	var uplinkTraffic, downlinkTraffic uint64
	newError("start jobs").AtDebug().WriteToLog()
	defer func() {
		newError(fmt.Sprintf("jobs info: addded %d users, deleteted %d users, downlink traffic %s, uplink traffic %s",
			addedUserCount, deletedUserCount, bytefmt.ByteSize(downlinkTraffic), bytefmt.ByteSize(uplinkTraffic))).AtInfo().WriteToLog()
	}()

	uplinkTraffic, downlinkTraffic, err = p.getTraffic()
	if err != nil {
		return
	}

	addedUserCount, deletedUserCount, err = p.syncUser()
	return
}

func (p *Panel) getTraffic() (downlinkTotal uint64, uplinkTotal uint64, err error) {
	var downlink, uplink uint64
	for _, user := range p.userModels {
		downlink, err = p.statsServiceClient.getUserDownlink(user.Email)
		if err != nil {
			return
		}

		uplink, err = p.statsServiceClient.getUserUplink(user.Email)
		if err != nil {
			return
		}

		if uplink+downlink > 0 {
			log := UserTrafficLog{
				UserID:   user.ID,
				Uplink:   uplink,
				Downlink: downlink,
				NodeID:   p.globalConfig.myPluginConfig.NodeID,
				Rate:     p.globalConfig.myPluginConfig.TrafficRate,
				Traffic:  bytefmt.ByteSize(uplink + downlink),
			}

			p.db.DB.Create(&log)
		}

		downlinkTotal += downlink
		uplinkTotal += uplink
	}

	return
}

func (p *Panel) syncUser() (addedUserCount, deletedUserCount int, err error) {
	userModels, err := p.db.GetAllUsers()
	if err != nil {
		return 0, 0, err
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
				return
			}
			deletedUserCount++
		}
	}

	// Add
	for _, userModel := range addUserModels {
		if err = p.handlerServiceClient.AddUser(p.convertUser(userModel)); err != nil {
			return
		}
		p.userModels = append(p.userModels, userModel)
		addedUserCount++
	}

	return
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
		if user == *u {
			return i
		}
	}
	return -1
}

func inUserModels(u *UserModel, userModels []UserModel) bool {
	return findUserModelIndex(u, userModels) != -1
}

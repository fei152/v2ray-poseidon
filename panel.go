package v2ray_ssrpanel_plugin

import (
	"fmt"
	"google.golang.org/grpc"
	"time"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

type Panel struct {
	handlerServiceClient *HandlerServiceClient
	db                   *DB
	userModels           []UserModel
}

func NewPanel(gRPCConn *grpc.ClientConn, db *DB) *Panel {
	return &Panel{
		db:                   db,
		handlerServiceClient: NewHandlerServiceClient(gRPCConn, "proxy"),
	}
}

func (p *Panel) Start() {
	for {
		if err := p.do(); err != nil {
			newError("panel#do").Base(err).AtError().WriteToLog()
		}
		time.Sleep(10 * time.Second)
	}
}

func (p *Panel) do() error {
	newError("start doing ssr panel jobs").AtWarning().WriteToLog()
	defer newError("finished doing ssr panel jobs").AtWarning().WriteToLog()

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

	fmt.Println(addUserModels)
	fmt.Println(delUserModels)
	// Add
	p.userModels = append(p.userModels, addUserModels...)
	for _, userModel := range addUserModels {
		p.handlerServiceClient.AddUser(p.convertUser(userModel))
	}

	return nil
}

func (p *Panel) convertUser(userModel UserModel) *protocol.User {
	return &protocol.User{
		// todo
		Level: 0,
		Email: userModel.Email,
		Account: serial.ToTypedMessage(&vmess.Account{
			Id: userModel.VmessID,
			// todo
			AlterId: 1,
			SecuritySettings: &protocol.SecurityConfig{
				// todo
				Type: protocol.SecurityType_AUTO,
			},
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

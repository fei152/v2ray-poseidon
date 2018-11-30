package v2ray_ssrpanel_plugin

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	statsservice "v2ray.com/core/app/stats/command"
)

type StatsServiceClient struct {
	statsservice.StatsServiceClient
}

func NewStatsServiceClient(client *grpc.ClientConn) *StatsServiceClient {
	return &StatsServiceClient{
		StatsServiceClient: statsservice.NewStatsServiceClient(client),
	}
}

func (s *StatsServiceClient) getUserUplink(email string) (*statsservice.Stat, error) {
	return s.getUserTraffic(fmt.Sprintf("user>>>%s>>>traffic>>>uplink", email), true)
}

func (s *StatsServiceClient) getUserDownlink(email string) (*statsservice.Stat, error) {
	return s.getUserTraffic(fmt.Sprintf("user>>>%s>>>traffic>>>downlink", email), true)
}

func (s *StatsServiceClient) getUserTraffic(name string, reset bool) (*statsservice.Stat, error) {
	req := &statsservice.GetStatsRequest{
		Name:   name,
		Reset_: reset,
	}

	res, err := s.GetStats(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res.Stat, nil
}

package main

import (
	"github.com/ColetteContreras/v2ray-ssrpanel-plugin"
	"v2ray.com/core"
)

func GetPluginMetadata() core.PluginMetadata {
	v2ray_ssrpanel_plugin.Run()

	return core.PluginMetadata{
		Name: "SSR Panel",
	}
}


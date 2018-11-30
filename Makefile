.PHONEY: build
build:
	go build -o /tmp/v2ray_ssrpanel_plugin.so -buildmode=plugin ./main

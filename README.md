# a V2ray plugin for SSRPanel

Only one thing user should do is that setting up the database connection, without doing that user needn't do anything!

### Features

- Sync user from SSRPanel database to v2ray
- Log user traffic

### V2ray Configuration demo

[Demo config](./release/ssrpanel_plugin-demo.json)

### Screenshot

```
$ tree plugins
plugins
└── ssrpanel.so

0 directories, 1 file

$ ./v2ray -plugin
V2Ray 4.8 (Po) Custom
A unified platform for anti-censorship.
[Info] SSRPanelPlugin: Connecting database...
[Info] SSRPanelPlugin: Connected
[Info] v2ray.com/core: plugin (SSR Panel) loaded.
[Warning] v2ray.com/core: V2Ray 4.8 started
[Warning] SSRPanelPlugin: Connected gRPC server "127.0.0.1:10085"
[Info] [316774646] v2ray.com/core/app/dispatcher: taking detour [api] for [tcp:127.0.0.1:0]
[Warning] SSRPanelPlugin: + 3 users, - 0 users, ↓ 0, ↑ 0, online 0
127.0.0.1:59693 accepted tcp:isoredirect.centos.org:80
[Info] [3628777663] v2ray.com/core/proxy/vmess/inbound: received request for tcp:isoredirect.centos.org:80
[Info] [3628777663] v2ray.com/core/app/dispatcher: default route for tcp:isoredirect.centos.org:80
[Info] [3628777663] v2ray.com/core/proxy/freedom: opening connection to tcp:isoredirect.centos.org:80
[Info] [3628777663] v2ray.com/core/transport/internet/tcp: dialing TCP to tcp:isoredirect.centos.org:80
[Warning] SSRPanelPlugin: + 0 users, - 0 users, ↓ 0, ↑ 200B, online 1
[Info] [3628777663] v2ray.com/core/app/proxyman/outbound: failed to process outbound traffic > v2ray.com/core/proxy/freedom: connection ends > context canceled
127.0.0.1:59704 accepted tcp:ftp.kaist.ac.kr:80
[Info] [218896217] v2ray.com/core/proxy/vmess/inbound: received request for tcp:ftp.kaist.ac.kr:80
[Info] [218896217] v2ray.com/core/app/dispatcher: default route for tcp:ftp.kaist.ac.kr:80
[Info] [218896217] v2ray.com/core/proxy/freedom: opening connection to tcp:ftp.kaist.ac.kr:80
[Info] [218896217] v2ray.com/core/transport/internet/tcp: dialing TCP to tcp:ftp.kaist.ac.kr:80
[Warning] SSRPanelPlugin: + 0 users, - 0 users, ↓ 4.3M, ↑ 200B, online 1
[Warning] SSRPanelPlugin: + 0 users, - 0 users, ↓ 13.6M, ↑ 0, online 1
[Warning] SSRPanelPlugin: + 0 users, - 0 users, ↓ 12.7M, ↑ 0, online 1
[Warning] SSRPanelPlugin: + 0 users, - 0 users, ↓ 15.4M, ↑ 0, online 1
```

### References

- [V2ray](https://github.com/v2ray/v2ray-core)
- [SSRPanel](https://github.com/ssrpanel/SSRPanel)
- [Go Plugin Package](https://golang.org/pkg/plugin)

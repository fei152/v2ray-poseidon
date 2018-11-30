# a V2ray plugin for SSRPanel

Only one thing user should do is that setting up the database connection, without doing that user needn't do anything!

### Features

- Sync user from SSRPanel database to v2ray
- Log user traffic

### V2ray Configuration demo

```json
{
  "log": {
    "loglevel": "debug"
  },
  "api": {
    "tag": "api",
    "services": [
      "HandlerService",
      "LoggerService",
      "StatsService"
    ]
  },
  "stats": {},
  "inbounds": [{
    "port": 10086,
    "protocol": "vmess",
    "tag": "proxy"
  },{
    "listen": "127.0.0.1",
    "port": 10085,
    "protocol": "dokodemo-door",
    "settings": {
      "address": "127.0.0.1"
    },
    "tag": "api"
  }],
  "outbounds": [{
    "protocol": "freedom"
  }],
  "routing": {
    "rules": [{
      "type": "field",
      "inboundTag": [ "api" ],
      "outboundTag": "api"
    }],
    "strategy": "rules"
  },
  "policy": {
    "levels": {
      "0": {
        "statsUserUplink": true,
        "statsUserDownlink": true
      }
    },
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true
    }
  },


  "other": {
    "plugins": {

      "ssrpanel": {
        // Node id on your SSR Panel
        "nodeId": 1,
        // every N seconds
        "checkRate": 60,
        // traffic rate
        "trafficRate": 1.0,
	    // gRPC address
	    "gRPCAddr": "127.0.0.1:10085",
	    // user config
	    "user": {
          // inbound tag, which inbound you would like add user to
          "inboundTag": "proxy",
          "level": 0,
          "alterId": 16,
          "security": "none"
	    },
        // db connection
        "mysql": {
          "host": "127.0.0.1",
          "port": 3306,
          "user": "root",
          "password": "ssrpanel",
          "dbname": "ssrpanel"
        }
      }


    }
  }


}
```

### References

- [V2ray](https://github.com/v2ray/v2ray-core)
- [SSRPanel](https://github.com/ssrpanel/SSRPanel)
- [Go Plugin Package](https://golang.org/pkg/plugin)

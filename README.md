# Poseidon -- A buildin V2ray plugin

## Why change the name of the repo

The name have to be changed for these reasons:

1. I've learnt that there are still many panels need to be adapt with V2ray
  - SSPanel-v3
  - WEBAPI SSRPanel（VNetPanel）
2. I'd like to provide a new deployment way
  - Poseidon-Master node
    - Automatic install Poseidon-V2ray node on new nodes which you added to the server list in your panel
    - Automatic config Poseidon-V2ray node
    - Real-time script output of deployment can be seen in Web browsers
    - Any time and any where you can execute shell commands in your poseidon-v2ray node
  - Poseidon-V2ray node
3. To config with ease, a bran-new [v2ray config generator](https://github.com/ColetteContreras/poseidon-v2ray-config-generator) is going to be introduced to Poseidon-master.

## Versions

### Tenet

All features is going to be available for any version, the only one difference is user scale. **They'll be in effect for users which user id is not greater than 50 for `Community` version**. Contact us to get an `Enterprise` version which has no limit of user scale, if needed.

- Community version, which will be released on [GitHub releases](https://github.com/ColetteContreras/v2ray-poseidon/releases)
- Enterprise version, which you are able to get it via [TG group: v2ray_poseidon](https://t.me/v2ray_poseidon)

## Contact

Get in touch via [TG group: v2ray_poseidon](https://t.me/v2ray_poseidon)

## Updates

- v1.0.0

  Breaking Changes:

  - config files structure has been changed, new structure is:

  ```diff
  {
  -  "ssrpanel": {
  +  "poseidon": {
      ... ... 
    }
  }
  ```
  
  - Added IP limit
    - It's a number that how many ip a user can use at the same moment
    - You can set user's `protocol_param` field on the database
  - Added rate limit
    - You should set user's `speed_limit_per_user` and `speed_limit_per_conn` fields on the database
  - Added closing user connections instantly after user has been deleted
    - By default the user's connections will not be disconnected after being deleted, but that is not we wanted.
  - !!!**Warning**: These two features are available if user's id is less or equal 50
  - To support all users, please contact with me via [TG group: v2ray_poseidon](https://t.me/v2ray_poseidon)

=========================

Only one thing user should do is that setting up the database connection, without doing that user needn't do anything!

### Features

- Sync user from SSRPanel database to v2ray
- Log user traffic

### Benefits

- No other requirements
  - It's  able to run if you could launch v2ray core
- Less memory usage
  - It just takes about 5MB to 10MB memories more than v2ray core
  - Small RAM VPS would be joyful
- Simplicity configuration

### Install on Linux

you may want to see docs, all the things as same as the official docs except install command.

[V2ray installation](https://www.v2ray.com/en/welcome/install.html)

```
curl -L -s https://raw.githubusercontent.com/ColetteContreras/v2ray-ssrpanel-plugin/master/install-release.sh | sudo bash
```

#### Uninstall

```
curl -L -s https://raw.githubusercontent.com/ColetteContreras/v2ray-ssrpanel-plugin/master/uninstall.sh | sudo bash
```

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
  "stats": {
    "trackIp": true
  },
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

  "ssrpanel": {
    // Node id on your SSR Panel
    "nodeId": 1,
    // every N seconds
    "checkRate": 60,
    // change this to true if you want to ignore users which has an empty vmess_id
    "ignoreEmptyVmessID": false,
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
```

### Contributing

Read [WiKi](https://github.com/ColetteContreras/v2ray-ssrpanel-plugin/wiki) carefully before submitting issues.

- Test and [report bugs](https://github.com/ColetteContreras/v2ray-ssrpanel-plugin/issues)
- Share your needs/experiences in [github issues](https://github.com/ColetteContreras/v2ray-ssrpanel-plugin/issues)
- Enhance documentation
- Contribute code by sending PR

### References

- [V2ray](https://github.com/v2ray/v2ray-core)
- [SSRPanel](https://github.com/ssrpanel/SSRPanel)

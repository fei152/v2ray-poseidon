# Poseidon -- A buildin V2ray plugin

### Contact

Get in touch via [TG group: v2ray_poseidon](https://t.me/v2ray_poseidon)

### Donation 

If you guys have enjoyed with me, you are able to donate USDT via [MugglePay 麻瓜宝TG支付钱包](https://telegram.me/MugglePayBot?start=8J9V8DCJ "麻瓜宝用户钱包") 

For example, if you have Binance (or any other exchange like Huobi), you can withdraw 1 USDT to MugglePay, and buy me **a cup of coffee**( I wonder that can 1 USDT afford it? ) by sending the message below to [@MugglePayBot](http://t.me/MugglePayBot):

`/pay @ColetteContreras 1 USDT`

##

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
curl -L -s https://raw.githubusercontent.com/ColetteContreras/v2ray-poseidon/master/install-release.sh | sudo bash
```

#### Uninstall

```
curl -L -s https://raw.githubusercontent.com/ColetteContreras/v2ray-poseidon/master/uninstall.sh | sudo bash
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

  "poseidon": {
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

### Acknowledgement

- [V2ray](https://github.com/v2ray/v2ray-core)
- [SSRPanel](https://github.com/ssrpanel/SSRPanel)
- [V2board](https://github.com/v2board/v2board)

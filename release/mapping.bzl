def gen_mappings(os, arch):
  return {
    "v2ray_core/release": "",
    "v2ray_ssrpanel_plugin/plugin/" + os + "/" + arch: "plugins",
    "v2ray_ssrpanel_plugin/release": "",
  }

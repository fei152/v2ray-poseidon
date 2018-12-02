package(default_visibility=["//visibility:public"])

load("@v2ray_ext//bazel:zip.bzl", "pkg_zip")
load("//release:mapping.bzl", "gen_mappings")

filegroup(
  name = "config_json",
  srcs = ["ssrpanel_plugin-demo.json"],
)

pkg_zip(
  name = "ssrpanel_darwin_amd64_package",
  srcs = [
    "@v2ray_core//release:v2ray_darwin_amd64_package",
    "//plugin:ssrpanel_darwin_amd64",
    ":config_json",
  ],
  out = "v2ray-ssrpanel-plugin-macos.zip",
  mappings = gen_mappings("darwin", "amd64"),
)

pkg_zip(
  name = "ssrpanel_linux_amd64_package",
  srcs = [
    "@v2ray_core//release:v2ray_linux_amd64_package",
    "//plugin:ssrpanel_linux_amd64",
    ":config_json",
  ],
  out = "v2ray-ssrpanel-plugin-linux-64.zip",
  mappings = gen_mappings("linux", "amd64"),
)

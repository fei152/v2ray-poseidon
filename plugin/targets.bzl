load("@v2ray_ext//bazel:build.bzl", "foreign_go_binary")
load("@v2ray_ext//bazel:plugin.bzl", "PLUGIN_SUPPORTED_OS")

def gen_targets(matrix):
  pkg = "github.com/ColetteContreras/v2ray-ssrpanel-plugin/plugin"
  output = "ssrpanel.so"
  for (os, arch) in matrix:

    if os not in PLUGIN_SUPPORTED_OS:
      continue

    bin_name = "ssrpanel_" + os + "_" + arch

    foreign_go_binary(
      name = bin_name,
      pkg = pkg,
      output = output,
      os = os,
      arch = arch,
      cgo_enabled = '1',
      buildmode = 'plugin',
    )


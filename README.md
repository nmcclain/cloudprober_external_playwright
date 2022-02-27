# cloudprober_external_playwright: a cloudprober external probe wrapper to run playwright tests

## Usage:

```
  cloudprober_external_playwright [options] <test_directory>
  cloudprober_external_playwright --help
  cloudprober_external_playwright --version

Options:
  --version                    Show version
  -h, --help                   Show this screen
```

This tool is intended to be called from cloudprober, for example:
```
probe {
  name: "playwright"
  type: EXTERNAL
  targets { dummy_targets {} }
  external_probe {
    mode: ONCE
    command: "cloudprober_external_playwright tests"
  }
  interval_msec: 30000  # 30s
  timeout_msec: 10000   # 10s
}
```

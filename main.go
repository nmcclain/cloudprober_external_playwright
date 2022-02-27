package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/docopt/docopt-go"
)

const version = "0.0.2"

var usage = `cloudprober_external_playwright: a cloudprober external probe wrapper to run playwright tests
Usage:
  cloudprober_external_playwright [options] <test_directory>
  cloudprober_external_playwright --help
  cloudprober_external_playwright --version

Options:
  -d, --debug                  Enable debugging output
  --version                    Show version
  -h, --help                   Show this screen
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		return
	}
	testdir := args["<test_directory>"].(string)
	sanitizeMetricRE := regexp.MustCompile(`[\.\\\/ =]`) // "It must match the regex [a-zA-Z_:][a-zA-Z0-9_:]*."
	suiteTitle := sanitizeMetricRE.ReplaceAllString(testdir, "_")

	r, err := runPlaywright(testdir, args["--debug"].(bool))
	if err != nil {
		log.Printf("Error with Playwright execution: %s", err)
		fmt.Printf("all_tests_passing{suite=%s} 0\n", suiteTitle)
		return
	}

	if r.AllTestsSuccessful {
		fmt.Printf("all_tests_passing{suite=%s} 1\n", suiteTitle)
	} else {
		fmt.Printf("all_tests_passing{suite=%s} 0\n", suiteTitle)
	}

	for _, suite := range r.Suites {
		for _, spec := range suite.Specs {
			title := sanitizeMetricRE.ReplaceAllString(spec.Title, "_")
			if spec.Ok {
				fmt.Printf("test_passing{suite=%s,title=%s} 1\n", suiteTitle, title)
			} else {
				fmt.Printf("test_passing{suite=%s,title=%s} 0\n", suiteTitle, title)
			}
			fmt.Printf("test_duration{suite=%s,title=%s} %d\n", suiteTitle, title, spec.TestDuration)
		}
	}
}

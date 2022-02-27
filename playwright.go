package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/buger/jsonparser"
)

const NPXCMD = "/usr/local/bin/npx"

type PlaywrightResults struct {
	AllTestsSuccessful bool
	PlaywrightVersion  string
	Suites             []PlaywrightSuite
}

type PlaywrightSuite struct {
	Title              string
	File               string
	AllTestsSuccessful bool
	Specs              []PlaywrightSpec
}
type PlaywrightSpec struct {
	Title        string
	Ok           bool
	TestDuration int64
}

type PlaywrightTest struct {
	Title string
	File  string
	Suite string
	// Name  string
	// Index string
}

func runPlaywright(suiteDir string, debug bool) (PlaywrightResults, error) {
	res := PlaywrightResults{AllTestsSuccessful: true}
	if len(suiteDir) < 1 {
		return res, fmt.Errorf("Invalid test: %s", suiteDir)
	}
	args := []string{"playwright", "test", "--reporter", "json", suiteDir}
	if debug {
		log.Printf("Running Playwright: %s %s", NPXCMD, strings.Join(args, " "))
	}
	out, _ := exec.Command(NPXCMD, args...).Output() // no err check here - we get JSON even when the test fails

	if debug {
		log.Printf("Playwright output: %s", out)
	}

	var err error
	res.PlaywrightVersion, err = jsonparser.GetString(out, "config", "version")
	if err != nil {
		return PlaywrightResults{AllTestsSuccessful: false}, err
	}

	jsonparser.ArrayEach(out, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		suite := PlaywrightSuite{}
		suite.AllTestsSuccessful = true
		if err != nil {
			log.Printf("Error parsing Playwright JSON: %s", err)
			res.AllTestsSuccessful = false
			return
		}
		suite.Title, err = jsonparser.GetString(value, "title")
		if err != nil {
			log.Printf("Error parsing Playwright JSON: %s", err)
			res.AllTestsSuccessful = false
			return
		}
		suite.File, err = jsonparser.GetString(value, "file")
		if err != nil {
			log.Printf("Error parsing Playwright JSON: %s", err)
			res.AllTestsSuccessful = false
			return
		}

		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			spec := PlaywrightSpec{}
			if err != nil {
				log.Printf("Error parsing Playwright JSON: %s", err)
				res.AllTestsSuccessful = false
				return
			}
			spec.Title, err = jsonparser.GetString(value, "title")
			if err != nil {
				log.Printf("Error parsing Playwright JSON: %s", err)
				res.AllTestsSuccessful = false
				return
			}
			spec.Ok, err = jsonparser.GetBoolean(value, "ok")
			if err != nil {
				log.Printf("Error parsing Playwright JSON: %s", err)
				res.AllTestsSuccessful = false
				return
			}
			if !spec.Ok {
				suite.AllTestsSuccessful = false
				res.AllTestsSuccessful = false
			}

			jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					resultDuration, err := jsonparser.GetInt(value, "duration")
					if err != nil {
						log.Printf("Error parsing Playwright JSON: %s", err)
						res.AllTestsSuccessful = false
						return
					}
					spec.TestDuration += resultDuration
				}, "results")
			}, "tests")

			suite.Specs = append(suite.Specs, spec)
		}, "specs")

		res.Suites = append(res.Suites, suite)
	}, "suites")

	if debug {
		log.Printf("Final result: %+v", res)
	}
	return res, nil
}

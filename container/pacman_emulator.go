package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Dependency struct {
	Install []string `json:"install"`
	Make    []string `json:"make"`
}

func readAll(path string) []byte {
	fos, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	txt, err := ioutil.ReadAll(fos)
	if err != nil {
		panic(err)
	}
	fos.Close()
	return txt
}

func applyRules(argsDeps []string, deps map[string]Dependency, install bool, make bool) (testingDeps []string, testingDepsVersions []string) {
	seps := []string{">=", "<=", "=", ">", ">"}
	for _, argDep := range argsDeps {
		var version string
		for _, sep := range seps {
			isep := strings.Index(argDep, sep)
			if isep != -1 {
				version = strings.TrimSpace(argDep[isep:])
				argDep = strings.TrimSpace(argDep[:isep])
				break
			}
		}
		if ruleDep, ok := deps[argDep]; ok {
			if install {
				for _, rd := range ruleDep.Install {
					testingDeps = append(testingDeps, rd)
					testingDepsVersions = append(testingDepsVersions, version)
				}
			}
			if make {
				for _, rd := range ruleDep.Make {
					testingDeps = append(testingDeps, rd)
					testingDepsVersions = append(testingDepsVersions, version)
				}
			}
		} else {
			testingDeps = append(testingDeps, argDep)
			testingDepsVersions = append(testingDepsVersions, version)
		}
	}
	return
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Pacman emulator: no argument provided")
		os.Exit(0)
	}

	// Read dependency rules
	var deps map[string]Dependency
	err := json.Unmarshal(readAll("/usr/share/pacman/rules.json"), &deps)
	//err := json.Unmarshal(readAll("./rules.json"), &deps)
	if err != nil {
		panic(err)
	}

	if os.Args[1] == "-T" {
		// Convert dependency list from Arch to Debian
		testingDeps, _ := applyRules(os.Args[2:], deps, true, true)
		// Search for missing dependencies
		var missingDeps []string
		for _, td := range testingDeps {
			cmd := exec.Command("dpkg-query", "--show", "--showformat", "${Package}", td)
			if err := cmd.Run(); err != nil {
				missingDeps = append(missingDeps, td)
			}
		}
		// Report
		if len(missingDeps) > 0 {
			fmt.Println(strings.Join(missingDeps, "\n"))
		}
	} else if os.Args[1] == "-S" {
		// Convert dependency list from Arch to Debian
		installDeps, _ := applyRules(os.Args[2:], deps, true, true)
		// Check installed dependencies
		for _, td := range installDeps {
			cmd := exec.Command("apt-get", "-y", "install", td)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				panic(cmd)
			}
		}
	} else if os.Args[1] == "-C" {
		// Convert dependency list from Arch to Debian
		// Not including "make" dependencies
		convDeps, convDepsVersions := applyRules(os.Args[2:], deps, true, false)
		// Report
		reports := make([]string, len(convDeps))
		for i := 0; i < len(convDeps); i++ {
			reports[i] += convDeps[i]
			if len(convDepsVersions[i]) > 0 {
				reports[i] += fmt.Sprintf(" (%s)", convDepsVersions[i])
			}

		}
		if len(convDeps) > 0 {
			fmt.Println(strings.Join(reports, ", "))
		}
	}
}

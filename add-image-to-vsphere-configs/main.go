package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/openshift/ci-tools/pkg/api"
)

func main() {

	err := filepath.Walk("./ci-operator/config", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".yaml") {
			var releaseBuildConfiguration api.ReleaseBuildConfiguration
			//fmt.Println(path)
			fileContents, err := os.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			if err := yaml.Unmarshal(fileContents, &releaseBuildConfiguration); err != nil {
				log.Fatal(err)
			}

			modifyConfig := false

			for _, t := range releaseBuildConfiguration.Tests {
				//godump.Dump(t)
				if t.MultiStageTestConfiguration != nil {
					if t.MultiStageTestConfiguration.Workflow != nil {
						if strings.Contains(*t.MultiStageTestConfiguration.Workflow, "vsphere") {
							modifyConfig = true
							fmt.Println(t.As)
						}
					}
				}
			}

			if modifyConfig {
				if _, ok := releaseBuildConfiguration.InputConfiguration.BaseImages["vsphere-ci-python"]; !ok {

					releaseBuildConfiguration.InputConfiguration.BaseImages["vsphere-ci-python"] = api.ImageStreamTagReference{
						Name:      "vsphere-python",
						Namespace: "ci",
						Tag:       "latest",
					}
					changed, err := yaml.Marshal(releaseBuildConfiguration)
					if err != nil {
						log.Fatal(err)
					}

					if err := os.WriteFile(path, changed, 0644); err != nil {
						log.Fatal(err)
					}

				}
			}
		}
		return nil

	})

	if err != nil {
		log.Fatal(err)
	}

}

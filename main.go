package main

import (
	"fmt"
	"log"

	flag "github.com/spf13/pflag"

	"github.com/codeWithUtkarsh/image-scan-poc/scan"
	"github.com/docker/docker/client"
)

type envVariable []string

func (i *envVariable) String() string {
	return "my string representation"
}

func (i *envVariable) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var envFlags envVariable

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Unable to create docker client")
	}

	commands := []string{"python -m ensurepip --upgrade", "pip3 freeze > requirements.txt", "pip3 install cyclonedx-bom==0.4.3 safety", "cyclonedx-py -j -o /tmp/sbom.json", "safety check -r requirements.txt --json --output /tmp/cve.json || true"} //mandatory input, hardcoded for now
	directoryToSaveGeneratedFiles := "/tmp"

	imagename := flag.String("imagename", "xyz", "Docker Imagename to scan")
	inputEnv := flag.StringArray("env", []string{}, "Environment Variable")

	flag.Parse()
	fmt.Println(*inputEnv)

	err = scan.ImageScanWithCustomCommands(cli, *imagename, commands, directoryToSaveGeneratedFiles, *inputEnv)
	if err != nil {
		log.Println(err)
	}
}

// func main() {
// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	// cli, err := client.NewEnvClient()
// 	if err != nil {
// 		log.Fatalf("Unable to create docker client")
// 	}

// 	imagename := "1645370/ortelius-test:latest"                                                                                                                                                                                                        //mandatory input
// 	commands := []string{"python -m ensurepip --upgrade", "pip3 freeze > requirements.txt", "pip3 install cyclonedx-bom==0.4.3 safety", "cyclonedx-py -j -o /tmp/sbom.json", "safety check -r requirements.txt --json --output /tmp/cve.json || true"} //mandatory input
// 	directoryToSaveGeneratedFiles := "/tmp"

// 	inputEnv := []string{"DB_HOST=192.168.225.51", "DB_PORT=9876"}

// 	err = ImageScanWithCustomCommands(cli, imagename, commands, directoryToSaveGeneratedFiles, inputEnv)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

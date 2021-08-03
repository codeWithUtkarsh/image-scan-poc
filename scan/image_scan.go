package scan

import (
	"fmt"
	"os"
	"log"

	"github.com/codeWithUtkarsh/image-scan-poc/functions"

	"github.com/BurntSushi/toml"
	"github.com/docker/docker/client"
)

type Config struct {
	Port          string
	ContainerName string
	UserName      string
	Password      string
}

func ImageScanWithCustomCommands(client *client.Client, imagename string, commands []string, dirToSave string, inputEnv []string) error {

	//---------- Loading configuration -------------
	path, errr := os.Getwd()
	if errr != nil {
	    log.Println(errr)
	}
	fmt.Println(path)

	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println("Error while reading configuration file")
		return err
	}

	//---------- Pulling image -------------
	err := functions.PullImage(client, config.UserName, config.Password, imagename) //imagename is ${ImageRegistry}:${ImageTag} eg: 1645370/test-imag:latest
	if err != nil {
		fmt.Println("Error while pulling image")
		return err
	}

	//---------- Start container -------------
	containerId, err := functions.RunContainer(client, imagename, config.ContainerName, config.Port, inputEnv)
	if err != nil {
		fmt.Println("Error while running container")
		return err
	}

	//---------- Execute commands inside container -------------
	err = functions.ExecCommand(client, containerId, commands)
	if err != nil {
		fmt.Println("Error while executing commands")

		//stop and remove container and return
		serr := functions.StopAndRemoveContainer(client, containerId)
		if serr != nil {
			fmt.Printf("Error while stoping and removing container; Manually remove container with name = %s", config.ContainerName)
			return serr
		}
		return err
	}

	// ---------- Copy generated files to host directory -------------
	err = functions.CopyGeneratedFile(client, containerId, dirToSave) //This method will also remove the container after task is completed
	if err != nil {
		fmt.Println("Error while copying file commands")

		//stop and remove container and return
		serr := functions.StopAndRemoveContainer(client, containerId)
		if serr != nil {
			fmt.Printf("Error while stoping and removing container; Manually remove container with name = %s", config.ContainerName)
			return serr
		}
		return err
	}

	err = functions.StopAndRemoveContainer(client, containerId)
	if err != nil {
		fmt.Printf("Error while stoping and removing container; Manually remove container with name = %s", config.ContainerName)
		return err
	}

	return nil
}

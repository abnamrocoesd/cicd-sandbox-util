package dockerprobe

import (
	"../model"
	"../util"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"strconv"
	"strings"
)

func ContainerInfoList(containers []types.Container, filterDockerRegister bool, filterDockerRegisterName string, labelPrefix string) []model.ContainerInfo {
	var containerInfoList []model.ContainerInfo
	for _, container := range containers {
		ports := parsePorts(container)
		volumes := parseVolumes(container)
		cleanedImageName := container.Image
		if filterDockerRegister {
			cleanedImageName = cleanImageName(container, filterDockerRegisterName)
		}
		containerInfo := model.ContainerInfo{
			Name:        strings.Replace(container.Labels[labelPrefix+util.LabelName], "\"", "", -1),
			Description: strings.Replace(container.Labels[labelPrefix+util.LabelDescription], "\"", "", -1),
			WebPort:     container.Labels[labelPrefix+util.LabelWebPort],
			WebPath:     container.Labels[labelPrefix+util.LabelWebPath],
			Ports:       ports,
			Volumes:     volumes,
			Created:     container.Created,
			Image:       cleanedImageName,
		}
		if strings.TrimSpace(containerInfo.Name) != "" {
			containerInfoList = append(containerInfoList, containerInfo)
		}
	}
	return containerInfoList
}

func cleanImageName(container types.Container, dockerRegistryPrefix string) string {
	name := container.Image
	return strings.Replace(name, dockerRegistryPrefix, "", 1)
}

func parseVolumes(container types.Container) []string {
	volumes := make([]string, len(container.Mounts))
	for _, mount := range container.Mounts {
		volumes = append(volumes, fmt.Sprintf("%s @%s", mount.Name, mount.Destination))
	}
	return volumes
}

func parsePorts(container types.Container) []string {
	ports := make([]string, len(container.Ports))
	for _, port := range container.Ports {
		// Ugly, but it works: https://stackoverflow.com/questions/41787620/convert-uint64-to-string-in-golang
		ports = append(ports, fmt.Sprintf("%s:%s", strconv.Itoa(int(port.PublicPort)), strconv.Itoa(int(port.PrivatePort))))
	}
	return ports
}

func ContainerList(labelFilter string, dockerHost string) ([]types.Container, error) {
	host := dockerHost
	fmt.Println(" > Probing Host: " + host)
	cli, err := client.NewClientWithOpts(client.WithVersion("1.35"))
	if err != nil {
		panic(err)
	}

	filter := filters.NewArgs()
	filter.Add("label", labelFilter) // TODO: make filter optional / parameter
	return cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filter})
}

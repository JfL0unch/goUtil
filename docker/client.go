package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

var cliForTest *Cli

type Cli struct {
	dockerCli         *dockerClient.Client
	runningContainers *sync.Map
}

func NewClient() (*Cli, error) {
	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
	if err != nil {
		return nil, errors.Wrapf(err, "new docker client")
	}
	return &Cli{
		dockerCli:         cli,
		runningContainers: &sync.Map{},
	}, nil
}

type containerState string

func (c containerState) String() string {
	return string(c)
}

const (
	containerStateRunning containerState = "running"
	containerStateExited  containerState = "exited"
)

func (c *Cli) Clear() {
	containerIds := c.RunningContainers()
	for _, cid := range containerIds {
		err := c.StopAndRemoveContainer(cid)
		if err != nil {
			logrus.Errorf("c.StopAndRemoveContaine(%s) err:%s", cid, err)
		}
	}

	err := c.PruneVolumes()
	if err != nil {
		logrus.Errorf("c.PruneVolumes() err:%s", err)
	}
}
func (c *Cli) StartRedisContainer(name, pwd, portNo string, volumeBinds []string) (string, error) {
	imageName := "redis:3.2"
	envs := []string{
		"TZ=Asia/Shanghai",
	}
	inPort, err := nat.NewPort("tcp", "6379")
	if err != nil {
		return "", err
	}
	portMap := make(nat.PortMap, 0)
	binds := make([]nat.PortBinding, 0, 1)
	binds = append(binds, nat.PortBinding{HostPort: portNo})
	portMap[inPort] = binds

	volumesImage := make(map[string]struct{}, 0)
	for _, v := range volumeBinds {
		splits := strings.Split(v, ":")
		if len(splits) > 0 {
			volumesImage[splits[0]] = struct{}{}
		}
	}

	//cmds := []string{"/bin/sh", "-c", "cat /foo/bar"}
	cmds := []string{"redis-server", "--requirepass", pwd}
	config := &container.Config{
		Image:        imageName,
		ExposedPorts: make(nat.PortSet, 10),
		Env:          envs,
		Cmd:          cmds,
		Volumes:      volumesImage,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portMap,
		Binds:        volumeBinds,
	}
	containerId, err := c.createContainer(name, config, hostConfig)
	if err != nil {
		return "", errors.Wrapf(err, "c.createContainer")
	}

	time.Sleep(time.Second * 10)
	return containerId, c.startContainer(containerId)
}
func (c *Cli) RunningContainers() []string {
	res := make([]string, 0)
	c.runningContainers.Range(func(key, val any) bool {
		res = append(res, key.(string))
		return true
	})
	return res
}
func (c *Cli) PruneVolumes() error {
	var err error
	_, err = c.dockerCli.VolumesPrune(context.Background(), filters.Args{})
	return err
}

func (c *Cli) StartMysqlDbContainer(name, dbName, dbUser, dbPwd, portNo string, volumeBinds []string) (string, error) {
	imageName := "mysql:latest"
	envs := []string{
		"MYSQL_DATABASE=" + dbName,
		"TZ=Asia/Shanghai",
	}
	if dbUser == "root" {
		envs = append(envs, "MYSQL_ROOT_PASSWORD="+dbPwd)
	} else {
		envs = append(envs, "MYSQL_USER="+dbUser)
		envs = append(envs, "MYSQL_PASSWORD="+dbPwd)
		envs = append(envs, "MYSQL_ROOT_PASSWORD=root")
	}
	inPort, err := nat.NewPort("tcp", "3306")
	if err != nil {
		return "", err
	}
	portMap := make(nat.PortMap, 0)
	binds := make([]nat.PortBinding, 0, 1)
	binds = append(binds, nat.PortBinding{HostPort: portNo})
	portMap[inPort] = binds

	volumesImage := make(map[string]struct{}, 0)
	for _, v := range volumeBinds {
		splits := strings.Split(v, ":")
		if len(splits) > 0 {
			volumesImage[splits[0]] = struct{}{}
		}
	}

	cmds := make([]string, 0)
	config := &container.Config{
		Image:        imageName,
		ExposedPorts: make(nat.PortSet, 10),
		Env:          envs,
		Cmd:          cmds,
		Volumes:      volumesImage,
	}
	hostConfig := &container.HostConfig{
		PortBindings: portMap,
		Binds:        volumeBinds,
	}

	containerId, err := c.createContainer(name, config, hostConfig)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Second * 10)
	return containerId, c.startContainer(containerId)
}

// StopAndRemoveContainer  停止后删除
func (c *Cli) StopAndRemoveContainer(containerId string) error {
	err := c.stopContainer(containerId)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	return c.removeContainer(containerId)

}

/*
	参考文档：https://docs.docker.com/engine/api/sdk/examples/
*/
func (c *Cli) imageExist(imageName string) (bool, error) {
	iList, err := c.dockerCli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, err
	}
	for _, v := range iList {
		for _, tag := range v.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}

	return false, nil
}

/*
	参考文档：https://docs.docker.com/engine/api/sdk/examples/
*/
func (c *Cli) createContainer(containerName string, config *container.Config, hostConfig *container.HostConfig) (string, error) {

	imageName := config.Image
	ok, err := c.imageExist(imageName)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New(fmt.Sprintf("image[%s] not exist,please  pull this image by [docker pull %s] ", imageName, imageName))
	}

	// 判断是否有相同容器，如有则停止并删除他
	ct, err := c.queryContainerByName(containerName)
	if err != nil {
		logrus.Errorf("dockerUtils createContainer,c.queryContainerByName err:%s", err)
	}
	if ct != nil {
		containerId := ct.ID
		if ct.State == containerStateRunning.String() {
			err = c.StopAndRemoveContainer(containerId)
			if err != nil {
				logrus.Errorf(" dockerUtils createContainer,c.StopAndRemoveContainer err:%s", err)
			}
		}

		if ct.State == containerStateExited.String() {
			err = c.removeContainer(containerId)
			if err != nil {
				logrus.Errorf(" dockerUtils createContainer,c.removeContainer err:%s", err)
			}
		}
	}
	time.Sleep(time.Second)

	resp, err := c.dockerCli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", err
	}
	c.runningContainers.Store(resp.ID, struct{}{})
	return resp.ID, nil

}
func (c *Cli) queryContainerByName(containerName string) (*types.Container, error) {
	args := filters.NewArgs()
	args.Add("name", containerName)
	cts, err := c.dockerCli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: args,
		All:     true,
		Limit:   1,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "dockerUtils queryContainerByName,c.dockerCli.ContainerList(%s) err:%s", containerName, err)
	}
	if len(cts) > 0 {
		return &cts[0], nil
	} else {
		return nil, nil
	}
}
func (c *Cli) containersList() ([]types.Container, error) {
	cts, err := c.dockerCli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}
	if len(cts) > 0 {
		return cts, nil
	} else {
		return nil, nil
	}
}
func (c *Cli) startContainer(containerId string) error {
	err := c.dockerCli.ContainerStart(context.Background(), containerId, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil

}
func (c *Cli) stopContainer(containerId string) error {
	timeout := 30
	opt := container.StopOptions{
		Timeout: &timeout,
	}
	err := c.dockerCli.ContainerStop(context.Background(), containerId, opt)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cli) removeContainer(containerId string) error {
	err := c.dockerCli.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	c.runningContainers.Delete(containerId)
	return nil
}

func (c *Cli) listContainers() {
	containers, err := c.dockerCli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, cont := range containers {
		fmt.Printf("%s %s\n", cont.ID[:10], cont.Image)
	}
}

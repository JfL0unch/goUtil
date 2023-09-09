package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func InitTest() error {
	var err error
	cliForTest, err = NewClient()
	return err
}

// go test -v -run Test_listContainers
func Test_listContainers(t *testing.T) {
	_ = InitTest()

	Convey("1", t, func() {
		cliForTest.listContainers()
	})

}

// go test -v -run Test_StartMysqlDb -args ../../../config_test.yml
func Test_StartMysqlDb(t *testing.T) {
	_ = InitTest()

	Convey("1", t, func() {

		dbAddr := "127.0.0.1:18306"
		dbName := "testdb"
		dbUser := "root"
		dbPwd := "rootpwd"

		aserverDbPort := "18306"
		aserverSplits := strings.Split(dbAddr, ":")
		if len(aserverSplits) >= 2 {
			aserverDbPort = aserverSplits[1]
		}
		volumesImage := make([]string, 0)
		volumesImage = append(volumesImage, "/tmp/mysql/conf:/etc/mysql/conf.d")
		volumesImage = append(volumesImage, "/tmp/mysql/db:/var/lib/mysql")
		containerId, err := cliForTest.StartMysqlDbContainer("test", dbName, dbUser, dbPwd, aserverDbPort, volumesImage)
		if err != nil {
			t.Error(err)
		}
		fmt.Print(containerId)
	})
}

// go test -v -run Test_StartRedisContainer
func Test_StartRedisContainer(t *testing.T) {
	_ = InitTest()

	Convey("true positive", t, func() {
		pwd := "test:12345"
		volumesImage := make([]string, 0)
		containerId, err := cliForTest.StartRedisContainer("test", pwd, "26379", volumesImage)
		if err != nil {
			t.Error(err)
		}
		fmt.Print(containerId)
	})

}

// go test -v -run Test_StopAndRemoveContainer
func Test_StopAndRemoveContainer(t *testing.T) {
	_ = InitTest()

	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containerId := "e28f6a6c42d2"
			err := cliForTest.StopAndRemoveContainer(containerId)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// go test -v -run Test_createContainer
func Test_createContainer(t *testing.T) {
	_ = InitTest()

	Convey("ok", t, func() {
		imageName := "mysql:latest"
		containerName := "mysql-2059-XVlB"
		envs := []string{
			"MYSQL_DATABASE=db",
			"MYSQL_ROOT_PASSWORD=12345",
		}
		inPort, err := nat.NewPort("tcp", "18306")
		portMap := make(nat.PortMap, 0)
		volumesImage := make(map[string]struct{}, 0)
		volumeBinds := make([]string, 0)
		volumeBinds = append(volumeBinds, "/tmp/mysql/conf:/etc/mysql/conf.d")
		volumeBinds = append(volumeBinds, "/tmp/mysql/db:/var/lib/mysql")
		for _, v := range volumeBinds {
			splits := strings.Split(v, ":")
			if len(splits) > 0 {
				volumesImage[splits[0]] = struct{}{}
			}
		}
		binds := make([]nat.PortBinding, 0, 1)
		binds = append(binds, nat.PortBinding{HostPort: "23306"})
		portMap[inPort] = binds
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
		containerId, err := cliForTest.createContainer(containerName, config, hostConfig)

		So(err, ShouldBeNil)
		So(containerId, ShouldNotEqual, "")
	})

}

// go test -v -run Test_queryContainer
func Test_queryContainer(t *testing.T) {
	err := InitTest()
	if err != nil {
		t.Error(err)
	}
	Convey("running", t, func() {
		ctName := "mysqlmock_"
		gotContainer, err := cliForTest.queryContainerByName(ctName)

		So(err, ShouldBeNil)
		So(gotContainer, ShouldNotBeNil)
		So(gotContainer.ID[:12], ShouldEqual, "ea23b013706f")

		So(gotContainer.State, ShouldEqual, containerStateRunning)
	})

	Convey("all", t, func() {
		gotContainers, err := cliForTest.containersList()

		So(err, ShouldBeNil)
		So(len(gotContainers), ShouldBeGreaterThan, 0)
		for _, v := range gotContainers {
			fmt.Println(v.Names)
		}
	})

	Convey("exited", t, func() {
		ctName := "mysql-2059-XVlB"
		gotContainer, err := cliForTest.queryContainerByName(ctName)

		So(gotContainer, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(gotContainer.ID, ShouldEqual, "408c186edfee97c5c1cb8bc86a92e70c0f827c118e6e4ea24c18c3b3684e461d")

		So(gotContainer.State, ShouldEqual, containerStateExited)
	})
}

// go test -v -run Test_startContainer
func Test_startContainer(t *testing.T) {
	_ = InitTest()

	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containerId := "3e7000bcc961"
			err := cliForTest.startContainer(containerId)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// go test -v -run Test_stopContainer
func Test_stopContainer(t *testing.T) {
	_ = InitTest()

	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containerId := "3e7000bcc961"
			err := cliForTest.stopContainer(containerId)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// go test -v -run Test_removeContainer
func Test_removeContainer(t *testing.T) {
	_ = InitTest()

	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containerId := "3e7000bcc961"
			err := cliForTest.removeContainer(containerId)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// go test -v -run Test_imageExist
func Test_imageExist(t *testing.T) {
	_ = InitTest()

	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containerId := "redis:3.2"
			_, err := cliForTest.imageExist(containerId)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

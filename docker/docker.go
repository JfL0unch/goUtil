package docker

type Intf interface {
	/*
		创建一个mysql database
		入参:
		name 容器名
		dbName	数据库名
		dbUser		数据库用户
		dbPwd		数据库密码
		portNo 	对外暴露端口号

		出参:
		containerId  容器ID
	*/
	StartMysqlDbContainer(name, dbName, dbUser, dbPwd, portNo string, volumeBinds []string) (containerId string, err error)

	StartRedisContainer(name, auth, portNo string, volumeBinds []string) (containerId string, err error)

	StopAndRemoveContainer(containerId string) error

	RunningContainers() []string

	// Clear 清理测试夹具
	Clear()
}

package config

import "flag"

type Config struct {
	ReplicaHost string
	ReplicaPort string
	ClusterType string
	Offset      int
	ID          string
	ListenAddr  string
}

var Configs Config

func SetConfig() {
	port := flag.String("port", "6379", "redis port")
	replicaHost := flag.String("replicaof", "", "redis master cluster")
	replicaPort := "0"

	flag.Parse()
	args := flag.Args()

	clusterType := "master"

	if len(*replicaHost) > 0 {
		clusterType = "slave"
	}
	if len(args) > 0 {
		replicaPort = args[0]
	}

	Configs = Config{
		ClusterType: clusterType,
		ListenAddr:  ":" + *port,
		ReplicaHost: *replicaHost,
		ReplicaPort: replicaPort,
		ID:          "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		Offset:      0,
	}

}

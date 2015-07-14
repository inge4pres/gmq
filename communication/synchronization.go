package gmqnet

import (
	m "gmq/configuration"
	_ "net"
)

var cluster map[string]*Server

func init() {
	cluster = make(map[string]*Server)
}

func ClusterDiscovery(params *m.Params) []*Server {
	return nil
}

func syncWithCluster(mex *Message) {

}

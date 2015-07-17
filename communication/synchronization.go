package gmqnet

import (
	"errors"
	m "gmq/configuration"
	"net"
	"strings"
	"time"
)

const (
	DEFAULT_CLUSTER_PORT = "4812"
	DEFAULT_TIMEOUT_MSEC = 4000
)

var cluster map[string]*Server

func init() {
	cluster = make(map[string]*Server)
}

func ClusterPeerDiscovery(params *m.Params) error {
	srvs := make(chan *Server, 254)

	if params.Cluster.Port == "" {
		params.Cluster.TimeoutMsec = DEFAULT_TIMEOUT_MSEC
	}
	if params.Cluster.Port == "" {
		params.Cluster.Port = DEFAULT_CLUSTER_PORT
	}

	ip := strings.Split(params.Cluster.Network, ".")
	if len(ip) != 4 {
		return errors.New("IP in NETWORK not in IPv4 format, IPv6 currently not supported")
	}
	for i := 1; i < 255; i++ {
		ip[3] = string(i)
		go func() {
			srvs <- dialServer(strings.Join(ip, "."), params.Cluster.Port, params.Cluster.TimeoutMsec)
		}()
		select {
		case server := <-srvs:
			cluster[server.LocalInet] = server
		default:
			continue
		}
	}
	checkLocalInet(params)
	return nil
}

func dialServer(ip, port string, timeout int64) *Server {
	if _, err := net.DialTimeout("tcp", ip+":"+port, time.Duration(timeout*1000000)); err != nil {
		return nil
	}
	return &Server{Port: port, LocalInet: ip, Proto: "tcp"}
}

func checkLocalInet(params *m.Params) {
	localaddr, err := net.Interfaces()
	if err != nil {
		cluster = nil
		return
	}
	for l := range localaddr {
		addr, _ := localaddr[l].Addrs()
		for a := range addr {
			if _, ok := cluster[addr[a].String()]; ok {
				delete(cluster, addr[a].String())
			}
		}
	}
	return
}

func syncMessage(mex *Message) error {
	errs := make(chan error)
	for c := range cluster {
		go func() {
			conn, err := net.Dial("tcp", cluster[c].LocalInet+":"+cluster[c].Port)
			if err != nil {
				errs <- err
			}
			written, err := conn.Write(WriteMessage(mex))
			if written < len(WriteMessage(mex)) {
				errs <- errors.New("Failed to write complete message in cluster")
			}
			errs <- err
		}()
	}
	return <-errs
}

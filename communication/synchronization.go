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

func ClusterPeerDiscovery(params *m.Params) ([]*Server, error) {
	srvs := make(chan *Server, 254)
	result := make([]*Server, 0)

	if params.Cluster.Port == "" {
		params.Cluster.TimeoutMsec = DEFAULT_TIMEOUT_MSEC
	}

	if params.Cluster.Port == "" {
		params.Cluster.Port = DEFAULT_CLUSTER_PORT
	}

	if strings.Index(params.Cluster.Cidr, "/") < 0 {
		return nil, errors.New("CIDR format in cluster configuration is invalid; use IP/PREFIX")
	}
	netinfo := strings.Split(params.Cluster.Cidr, "/")
	network := netinfo[0]
	prefix := netinfo[1]
	ip := strings.Split(network, ".")
	if len(ip) != 4 {
		return nil, errors.New("IP in CIDR not in IPv4 format, IPv6 currently not supported")
	}

	switch prefix {
	case "24":
		for i := 1; i < 254; i++ {
			ip[3] = string(i)
			go func() {
				srvs <- dialServer(strings.Join(ip, "."), params.Cluster.Port, params.Cluster.TimeoutMsec)
			}()
			result = append(result, <-srvs)
		}
	case "16":

	case "8":

	default:
		return nil, errors.New("PREFIX in CIDR must be 24, 16 or 8")
	}
	return result, nil
}

func dialServer(ip, port string, timeout int64) *Server {
	if _, err := net.DialTimeout("tcp", ip+":"+port, time.Duration(timeout * 10>>6)); err != nil {
		return nil
	}
	return &Server{Port: port, LocalInet: ip, Proto: "tcp"}

}

func syncMessage(mex *Message) {

}

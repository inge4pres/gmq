package gmqnet

import (
	"errors"
	"gmq/configuration"
	"net"
	_ "strings"
	"sync"
	"time"
)

const (
	DEFAULT_CLUSTER_PROTO = "tcp"
	DEFAULT_CLUSTER_PORT  = "4812"
	DEFAULT_TIMEOUT_MSEC  = 4000
)

var cluster map[string]*gmqconf.Server

func init() {
	cluster = make(map[string]*gmqconf.Server)
}

func ClusterPeerDiscovery(params *gmqconf.Params) error {
	srvs := make(chan *gmqconf.Server, 255)

	if params.Cluster.Port == "" {
		params.Cluster.TimeoutMsec = DEFAULT_TIMEOUT_MSEC
	}
	if params.Cluster.Port == "" {
		params.Cluster.Port = DEFAULT_CLUSTER_PORT
	}
	if params.Cluster.Proto == "" {
		params.Cluster.Proto = DEFAULT_CLUSTER_PROTO
	}
	ip, ipNet, err := net.ParseCIDR(params.Cluster.Cidr)
	if err != nil {
		return errors.New("Wrong CIDR in Cluster configuration")
	}
	//	for i := 1; i < 255; i++ {
	//		ip[3] = string(i)
	//		go func() {
	//			srvs <- dialServer(strings.Join(ip, "."), params.Cluster.Port, params.Cluster.Proto, params.Cluster.TimeoutMsec)
	//		}()
	//		select {
	//		case server := <-srvs:
	//			cluster[server.LocalInet] = server
	//		default:
	//			continue
	//		}
	//	}
	var wg sync.WaitGroup
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			srvs <- dialServer(ip, params.Cluster.Port, params.Cluster.Proto, params.Cluster.TimeoutMsec)
		}(ip.String())
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

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func dialServer(ip, port, proto string, timeout int64) *gmqconf.Server {
	if _, err := net.DialTimeout(proto, ip+":"+port, time.Duration(timeout*1000000)); err != nil {
		return nil
	}
	return &gmqconf.Server{Port: port, LocalInet: ip, Proto: proto}
}

func checkLocalInet(params *gmqconf.Params) {
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

func syncMessage(mex []byte) error {
	errs := make(chan error)
	for c := range cluster {
		go func() {
			conn, err := net.Dial(cluster[c].Proto, cluster[c].LocalInet+":"+cluster[c].Port)
			if err != nil {
				errs <- err
			}
			written, err := conn.Write(mex)
			if written < len(mex) {
				errs <- errors.New("Failed to write complete message synchronization in cluster")
			}
			errs <- err
		}()
	}
	return <-errs
}

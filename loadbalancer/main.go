package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, server []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         server,
	}

}

func (s *simpleServer) Address() string { return s.addr }
func (s *simpleServer) IsAlive() bool   { return true }

func (s *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	// Increment AFTER selecting a valid server
	lb.roundRobinCount++
	return server
}

func (lb *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to address %q\n ", targetServer.Address())
	targetServer.Serve(w, r)
}

func InitializeServersList(filename string) []Server {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	servers := []Server{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Skip empty lines
		}
		addr := strings.TrimSpace(line)
		fmt.Println("adding server ", addr)
		if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
			addr = "http://" + addr // Default to HTTP if no scheme is provided
		}
		servers = append(servers, newSimpleServer(addr))
	}

	return servers
}

func main() {
	servers := InitializeServersList("hosts.txt")

	lb := NewLoadBalancer("8000", servers)

	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving request at localhost:%s\n", lb.port)

	if err := http.ListenAndServe(":"+lb.port, nil); err != nil {
		fmt.Errorf("Couldn't listen on port :%v", lb.port)
	}
}

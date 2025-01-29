package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func NewSimpleServer(addr string) *SimpleServer {
	serverURL, err := url.Parse(addr)
	HandleErr(err)

	return &SimpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverURL),
	}
}

type LoadBalancer struct {
	Port            string
	RoundRobinCount int
	Servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		RoundRobinCount: 0, // Fixed: Start counter at 0 instead of -1
		Port:            port,
		Servers:         servers,
	}
}

func HandleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(0)
	}
}

func (s *SimpleServer) Address() string { return s.addr }
func (s *SimpleServer) IsAlive() bool   { return true }
func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func (lb *LoadBalancer) GetNextAvailableServer() Server {
	server := lb.Servers[lb.RoundRobinCount%len(lb.Servers)]
	for !server.IsAlive() {
		lb.RoundRobinCount++
		server = lb.Servers[lb.RoundRobinCount%len(lb.Servers)]
	}
	lb.RoundRobinCount++
	return server
}

func (lb *LoadBalancer) ServerProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.GetNextAvailableServer()
	fmt.Printf("FORWARDING REQUEST TO ADDRESS %q\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

func main() {
	servers := []Server{
		NewSimpleServer("https://www.duckduckgo.com"),
		NewSimpleServer("https://www.chatgpt.com"),
		NewSimpleServer("https://www.monkeytype.com"),
		NewSimpleServer("https://www.youtube.com"),
	}

	lb := NewLoadBalancer("8079", servers)

	http.HandleFunc("/", lb.ServerProxy)
	fmt.Println("Serving requests at http://localhost:" + lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}

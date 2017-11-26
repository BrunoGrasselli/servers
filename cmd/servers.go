package main

import (
  "log"
  "fmt"
  "os"
  "strings"
  "net/http"
  "net/http/httputil"
  "github.com/tkanos/gonfig"
)

type ServerConfiguration struct {
  Port int
  Name string
}

type Configuration struct {
  TLD string
  Servers []ServerConfiguration
}

var servers = make(map[string]int)
var configuration = Configuration{}

func proxy(w http.ResponseWriter, r *http.Request) {
  serverName := strings.Split(r.URL.Host, ".")[0]

  director := func(req *http.Request) {
    req = r
    req.URL.Scheme = "http"
    req.URL.Host = fmt.Sprintf("localhost:%d", servers[serverName])
  }
  backend := &httputil.ReverseProxy{Director: director}
  backend.ServeHTTP(w, r)
}

func pac(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "function FindProxyForURL (url, host) { if (dnsDomainIs(host, '.%s')) { return 'PROXY 127.0.0.1:8080'; } return 'DIRECT'; }", configuration.TLD)
}

func main() {
  err := gonfig.GetConf(os.Args[1], &configuration)

  if err != nil {
    log.Fatal(err)
  }

  for _, server := range configuration.Servers {
    fmt.Printf("Serving http://%s.%s\n", server.Name, configuration.TLD)
    servers[server.Name] = server.Port
  }

  http.HandleFunc("/", proxy)

  http.HandleFunc("/proxy.pac", pac)

  log.Fatal(http.ListenAndServe(":8080", nil))
}

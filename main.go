package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	interfaceName = flag.String("interface", "eth0", "interface name")
	addr          = flag.String("addr", ":30066", "addr")
)

type Data struct {
}

func (d *Data) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip, err := GetRa(*interfaceName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, addr := range addrs {
		if strings.Contains(addr.String(), ip[:len(ip)-1]) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(addr.String()))
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}

// sudo docker build . -t nas_net_tool:latest
// sudo docker run -itd --network host --name nas_net_tool nas_net_tool -interface enp7s0

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)
	http.ListenAndServe(*addr, &Data{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sig
	log.Println("exit")
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

type Config struct {
	Port          string
	IPHeader      string
	Zone          string
	Server        string
	TsigName      string
	TsigSecret    string
	TsigAlgorithm string
	Hosts         []HostConfig
}

type HostConfig struct {
	Hostname string
	User     string
	Password string
}

var config = &Config{}

func validHostname(hostname string) bool {
	for _, hc := range config.Hosts {
		if hc.Hostname+"."+config.Zone == hostname {
			return true
		}
	}
	return false
}

func validAuth(hostname, user, password string) bool {
	for _, hc := range config.Hosts {
		if hc.Hostname+"."+config.Zone == hostname && hc.User == user && hc.Password == password {
			return true
		}
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Fprint(w, "911\n")
		return
	}

	params := r.URL.Query()

	hostname := params.Get("hostname")
	if !validHostname(hostname) {
		fmt.Fprint(w, "nohost\n")
		return
	}

	myipstr := params.Get("myip")
	if myipstr == "" && config.IPHeader != "" {
		myipstr = r.Header.Get(config.IPHeader)
	}
	if myipstr == "" {
		fmt.Fprint(w, "911\n")
		return
	}
	myip := net.ParseIP(myipstr)
	if myip == nil {
		fmt.Fprint(w, "911\n")
		return
	}
	rrtype := "AAAA"
	if myip.To4() != nil {
		rrtype = "A"
	}

	user, password, ok := r.BasicAuth()
	if !ok || !validAuth(hostname, user, password) {
		fmt.Fprint(w, "badauth\n")
		return
	}

	if !sendDNSUpdate(hostname, rrtype, myipstr) {
		fmt.Fprint(w, "dnserr\n")
		return
	}

	fmt.Fprintf(w, "good %s\n", myipstr)
}

func sendDNSUpdate(hostname, rrtype, ip string) bool {
	rr, err := dns.NewRR(fmt.Sprintf("%s. 30 %s %s", hostname, rrtype, ip))
	if err != nil {
		return false
	}

	msg := new(dns.Msg)
	msg.SetUpdate(config.Zone + ".")
	msg.RemoveRRset([]dns.RR{rr})
	msg.Insert([]dns.RR{rr})

	msg.SetTsig(config.TsigName, config.TsigAlgorithm, 300, time.Now().Unix())

	client := new(dns.Client)
	client.TsigSecret = map[string]string{config.TsigName: config.TsigSecret}

	reply, _, err := client.Exchange(msg, config.Server+":53")
	if err != nil || reply.Rcode != dns.RcodeSuccess {
		return false
	}

	return true
}

func loadConfig(configfile string) {
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	configfile := flag.String("f", "/etc/dnsupd.json", "path of the config file to use")
	flag.Parse()

	loadConfig(*configfile)

	port := ":80"
	if config.Port != "" {
		port = ":" + config.Port
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(port, nil))
}

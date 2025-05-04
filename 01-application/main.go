package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
)

type PageData struct {
	Hostname string
	Address  string
	Author   string
}

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("couldn't find a suitable IP address")
}

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	ip, err := getLocalIP()
	if err != nil {
		log.Fatal(err)
	}

	author := os.Getenv("AUTHOR")
	if author == "" {
		author = "Unknown"
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Error parsing the template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("A request was received from %s", r.RemoteAddr)
		data := PageData{
			Hostname: host,
			Address:  ip,
			Author:   author,
		}
		tmpl.Execute(w, data)
	})

	fmt.Println("Listening on port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package main

import (
	"fmt"
	"os"
	"text/template"
)

type StaticHost struct {
	NebulaAddress string
	Endpoint      []string
}

type FirewallRule struct {
	Port   string   `json:"port"`
	Proto  string   `json:"proto"`
	Groups []string `json:"groups"`
	Any    bool     `json:"any"`
}

type NodeConfigSchema struct {
	Address          string         `json:"address"`
	Lighthouses      []string       `json:"lighthouses"`
	AmLighthouse     bool           `json:"am_lighthouse"`
	StaticHosts      []StaticHost   `json:"static_hosts"`
	ListenPort       int            `json:"listen_port"`
	FirewallInbound  []FirewallRule `json:"firewall_inbound"`
	FirewallOutbound []FirewallRule `json:"firewall_outbound"`
	Cipher           string         `json:"cipher"`
	Cert             string         `json:"cert"`
	Netname          string
}

func genConfig(config NodeConfigSchema) (err error) {
	templateData, err := Asset("cmd/neutron/template.yml")
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/etc/nebula/%s/nebula.yml", config.Netname)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	tmpl, err := template.New("NebulaConfig").Parse(string(templateData))
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, config)
	if err != nil {
		return err
	}
	return f.Close()
}

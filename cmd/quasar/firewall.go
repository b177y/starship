package main

type FirewallRule struct {
	Port   string   `json:"port"`
	Proto  string   `json:"proto"`
	Groups []string `json:"groups"`
	Any    bool     `json:"any"`
}

func defaultRules() (inbound, outbound []FirewallRule) {
	inbound = append(inbound, FirewallRule{
		Port:   "any",
		Proto:  "icmp",
		Any:    true,
		Groups: []string{},
	})
	outbound = append(outbound, FirewallRule{
		Port:   "any",
		Proto:  "any",
		Any:    true,
		Groups: []string{},
	})
	return inbound, outbound
}

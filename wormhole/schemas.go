package wormhole

type RequestJoinSchema struct {
	Netname  string
	Nodename string
	Hostname string
	PubKey   string
}

type NodeIdentitySchema struct {
	Netname  string
	Nodename string
}

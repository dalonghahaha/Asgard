package providers

import "Asgard/client"

var (
	MasterClient *client.Master
)

func RegisterMaster() {
	MasterClient = client.NewMaster()
	go MasterClient.Report()
}

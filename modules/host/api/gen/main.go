package main

import (
	"github.com/starius/api2"
	"github.com/EvilRedHorse/pubaccess-node/modules/host/api"
)

func main() {
	api2.GenerateClient(api.GetRoutes)
}

// go run ./cmd
// go run -tags prime ./cmd
// tinygo flash -target xxx ./cmd

package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/runner"
	solNew "github.com/SunWizePB-Git/solNew"
)

var (
	id           = dean.GetEnv("ID", "sol01")
	name         = dean.GetEnv("NAME", "Sol")
	deployParams = dean.GetEnv("DEPLOY_PARAMS", "target=demo")
	port         = dean.GetEnv("PORT", "8000")
	portPrime    = dean.GetEnv("PORT_PRIME", "8001")
	user         = dean.GetEnv("USER", "")
	passwd       = dean.GetEnv("PASSWD", "")
	dialURLs     = dean.GetEnv("DIAL_URLS", "")
	ssids        = dean.GetEnv("WIFI_SSIDS", "")
	passphrases  = dean.GetEnv("WIFI_PASSPHRASES", "")
)

func main() {
	sol := solNew.New(id, "solNew", name).(*solNew.Sol) //calls function to make sol variable; 
	sol.SetDeployParams(deployParams) //sets deploy parameters (see above)
	sol.SetWifiAuth(ssids, passphrases) //koyeb passes secrets into application
	sol.SetDialURLs(dialURLs) //THIS is what is entered to point to a particular hub
	runner.Run(sol, port, portPrime, user, passwd, dialURLs) //"just run it!"
}

package consulsd

import (
	"fmt"
	"game/base/discovery"
	"game/config"
	"strconv"
	"testing"
)

func regService(sd discovery.Discovery) {
	host := []string{
		"192.168.1.2",
		"192.168.1.3",
		"192.168.1.4",
		"192.168.1.5",
	}
	for i := 0; i < 3; i++ {
		e := sd.Register(&discovery.ServiceDesc{
			Name: "Game",
			ID:   "Game" + strconv.FormatInt(int64(i), 10),
			Host: host[i],
			Port: 3220 + i,
		})
		fmt.Println(e)
	}
}

func unregService(sd discovery.Discovery) {

	//sd.Deregister("Game_0")
	//sd.Deregister("Login_0")
	//sd.Deregister("Gateway_0")
	sd.Deregister("Battle_0")

}

func TestConsul(t *testing.T) {
	sd := NewDiscovery(config.ConsulAddr, "")
	unregService(sd)
	return

}
func TestHealth(t *testing.T) {
	sd := NewDiscovery(config.ConsulAddr, "")
	ddd := sd.QueryServices()
	for _, d := range ddd {
		data := sd.Query(d)
		for _, m := range data {
			fmt.Printf("sss %+v", m)
		}
	}
}

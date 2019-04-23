package consulsd

import (
	"encoding/json"
	"game/base/discovery"

	"github.com/cihub/seelog"
	"github.com/hashicorp/consul/api"
)

func (cd *consulDiscovery) Query(name string) (ret []*discovery.ServiceDesc) {
	var e error
	ret, e = cd.directQuery(name)
	if e != nil {
		seelog.Errorf("query Service Error %s %v", name, e)
		return
	}
	return
}

// from github.com/micro/go-micro/registry/consul_registry.go
func (cd *consulDiscovery) directQuery(name string) (ret []*discovery.ServiceDesc, err error) {
	result, _, err := cd.client.Health().Service(name, "", false, nil)
	if err != nil {
		return nil, err
	}

	for _, s := range result {
		if s.Service.Service != name {
			continue
		}
		sd := &discovery.ServiceDesc{
			Name:  s.Service.Service,
			ID:    s.Service.ID,
			Host:  s.Service.Address,
			Port:  s.Service.Port,
			Meta:  s.Service.Meta,
			Tag:   s.Service.Tags,
			State: &discovery.ServiceState{},
		}
		for _, c := range s.Checks {
			if c.Name == s.Service.ID {
				json.Unmarshal([]byte(c.Output), sd.State)
			}
		}
		if !isServiceHealth(s) {
			sd.State.Invalid = true
		}
		ret = append(ret, sd)
	}
	return
}

func (cd *consulDiscovery) QueryServices() (ret []string) {
	ss, e := cd.client.Agent().Services()
	if e != nil {
		return
	}
	for _, s := range ss {
		exists := false
		for _, r := range ret {
			if r == s.Service {
				exists = true
				break
			}
		}
		if !exists {
			ret = append(ret, s.Service)
		}
	}
	return
}

func (cd *consulDiscovery) Register(svc *discovery.ServiceDesc) error {
	err := cd.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:                svc.ID,
		Name:              svc.Name,
		Address:           svc.Host,
		Port:              svc.Port,
		Meta:              svc.Meta,
		EnableTagOverride: true,
		Check: &api.AgentServiceCheck{
			CheckID: svc.ID,
			Name:    svc.ID,
			TTL:     "10s",
		},
	})
	return err
}

func (cd *consulDiscovery) Deregister(svcid string) error {
	return cd.client.Agent().ServiceDeregister(svcid)
}
func (cd *consulDiscovery) UpdateService(serviceID string, status, info string) error {
	if e := cd.client.Agent().UpdateTTL(serviceID, info, status); e != nil {
		return e
	}
	return nil
}

package sd

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

type Consul struct {
	Client  *api.Client
	Service *connect.Service
}

func NewConsulRegistration(serviceName string, instanceName string, httpHealthCheckAddr string) (*Consul, error) {
	// Create a Consul API client
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	serviceDef := &api.AgentServiceRegistration{
		ID:   instanceName,
		Name: serviceName,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s/__health", httpHealthCheckAddr),
			Interval: "10s",
			Status:   api.HealthPassing,
		}}

	if err := consulClient.Agent().ServiceRegister(serviceDef); err != nil {
		return nil, err
	}

	consulService, err := connect.NewService(serviceName, consulClient)
	if err != nil {
		return nil, err
	}
	defer consulService.Close()

	return &Consul{
		Client:  consulClient,
		Service: consulService,
	}, nil
}

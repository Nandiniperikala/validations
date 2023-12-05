package eureka

type RegistrationManager interface {
	ManageDiscovery(configs RegistrationVariables)
	RegisterWithSerivceRegistry()
	SendHeartBeat(configs RegistrationVariables)
	DeRegisterFromServiceRegistry(configs RegistrationVariables)
}

type RegistrationVariables struct {
	ServiceRegistryURL string
	InstanceId         string
}



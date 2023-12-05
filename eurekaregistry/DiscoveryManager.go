package eureka

var manager *EurekaRegistrationManager

func ManageDiscovery(configs RegistrationVariables) {
	manager = new(EurekaRegistrationManager)
	manager.RegisterWithSerivceRegistry(configs)
	manager.SendHeartBeat(configs)
}

func Cleanup(configs RegistrationVariables) {
	manager.DeRegisterFromServiceRegistry(configs)
}

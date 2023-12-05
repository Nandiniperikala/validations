package eureka

import (
	app "com.medisure.validations/config"
	"com.medisure.validations/eurekaregistry/helper"
	"github.com/carlescere/scheduler"
	"github.com/micro/micro/v3/service/logger"
	"os"
	"runtime"
	"strconv"
	"time"
)

type AppRegistrationBody struct {
	Instance InstanceDetails `json:"instance"`
}

type InstanceDetails struct {
	InstanceId       string         `json:"instanceId"`
	HostName         string         `json:"hostName"`
	App              string         `json:"app"`
	VipAddress       string         `json:"vipAddress"`
	SecureVipAddress string         `json:"secureVipAddress"`
	IpAddr           string         `json:"ipAddr"`
	Status           string         `json:"status"`
	Port             Port           `json:"port"`
	SecurePort       Port           `json:"securePort"`
	HealthCheckUrl   string         `json:"healthCheckUrl"`
	StatusPageUrl    string         `json:"statusPageUrl"`
	HomePageUrl      string         `json:"homePageUrl"`
	DataCenterInfo   DataCenterInfo `json:"dataCenterInfo"`
	LeaseInfo        LeaseInfo      `json:"leaseInfo"`
}
type Port struct {
	Port    string `json:"$"`
	Enabled string `json:"@enabled"`
}

type DataCenterInfo struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

type LeaseInfo struct {
	RenewalIntervalInSecs int `json:"renewalIntervalInSecs"`
	DurationInSecs        int `json:"durationInSecs"`
}

// This struct shall be responsible for manager to manage registration with Eureka
type EurekaRegistrationManager struct {
}

func (erm EurekaRegistrationManager) RegisterWithSerivceRegistry(eurekaConfigs RegistrationVariables) {
	body := erm.getBodyForEureka("STARTING", eurekaConfigs)
	helper.MakePostCall(eurekaConfigs.ServiceRegistryURL+"validations", body, nil)
	time.Sleep(10 * time.Second)
	bodyUP := erm.getBodyForEureka("UP", eurekaConfigs)
	helper.MakePostCall(eurekaConfigs.ServiceRegistryURL+"validations", bodyUP, nil)
}

func (erm EurekaRegistrationManager) SendHeartBeat(eurekaConfigs RegistrationVariables) {
	logger.Infof("Sending HeartBeat")
	job := func() {
		helper.MakePutCall(eurekaConfigs.ServiceRegistryURL+"validations/"+eurekaConfigs.InstanceId, nil, nil)
	}
	// Run every 5 seconds but not now.
	scheduler.Every(5).Seconds().Run(job)
	runtime.Goexit()
}
func (erm EurekaRegistrationManager) DeRegisterFromServiceRegistry(configs RegistrationVariables) {
	bodyDOWN := erm.getBodyForEureka("DOWN", configs)
	helper.MakePostCall(configs.ServiceRegistryURL+"validations", bodyDOWN, nil)
}

func (erm EurekaRegistrationManager) getBodyForEureka(status string, configs RegistrationVariables) *AppRegistrationBody {
	httpport := app.GetVal("GO_MICRO_SERVICE_PORT")
	hostname, _ := os.Hostname()
	applicationName := "validations"
	env := os.Getenv("GO_MICRO_ACTIVE_PROFILE")
	if env == "prod" {
		hostname = "validations"
	}

	ipAddress, err := helper.ExternalIP()
	if err != nil {
		logger.Errorf("Enable to find IP address form OS")
	}

	renewalStr := app.GetVal("GO_MICRO_RENEWALINTERVALINSEC")
	var renewal int
	renewal, _ = strconv.Atoi(renewalStr)

	durationStr := app.GetVal("GO_MICRO_DURATIONINSECS")
	var duration int
	duration, _ = strconv.Atoi(durationStr)

	port := Port{httpport, "true"}
	securePort := Port{"8443", "false"}
	leaseInfo := LeaseInfo{renewal, duration}
	dataCenterInfo := DataCenterInfo{"com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo", "MyOwn"}

	homePageUrl := "http://" + hostname + ":" + httpport + "/"
	statusPageUrl := "http://" + hostname + ":" + httpport + "/status"
	healthCheckUrl := "http://" + hostname + ":" + httpport + "/healthcheck"

	instance := InstanceDetails{configs.InstanceId, hostname, applicationName, applicationName, applicationName,
		ipAddress, status, port, securePort, healthCheckUrl, statusPageUrl, homePageUrl, dataCenterInfo, leaseInfo}

	body := &AppRegistrationBody{instance}
	return body
}

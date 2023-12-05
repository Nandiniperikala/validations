package main

import (
	auth "com.medisure.validations/auth"
	"com.medisure.validations/controllers"
	"github.com/asim/go-micro/v3"
	rabbitmq "com.medisure.validations/rabbitmq"
	"github.com/micro/micro/v3/service/logger"
	eureka "com.medisure.validations/eurekaregistry"
	"github.com/google/uuid"
	"com.medisure.validations/handler"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/http"
	mhttp "github.com/go-micro/plugins/v3/server/http"
   "github.com/gorilla/mux"
	app "com.medisure.validations/config"
)

var configurations eureka.RegistrationVariables

func main() {
	defer cleanup()
	app.Setconfig()
	auth.SetClient()
	handler.InitializeMongoDb()
	service_registry_url :=app.GetVal("GO_MICRO_SERVICE_REGISTRY_URL")
	InstanceId := "validations:"+uuid.New().String()
	configurations = eureka.RegistrationVariables {ServiceRegistryURL:service_registry_url,InstanceId:InstanceId}
	port :=app.GetVal("GO_MICRO_SERVICE_PORT")
	srv := micro.NewService(
		micro.Server(mhttp.NewServer()),
    )
	opts1 := []micro.Option{
		micro.Name("validations"),
		micro.Version("latest"),
		micro.Address(":"+port),
	}
	srv.Init(opts1...)
	r := mux.NewRouter().StrictSlash(true)
	registerRoutes(r)		
	var handlers http.Handler = r
	
	go eureka.ManageDiscovery(configurations)
	go rabbitmq.ConsumerPoliciesToValidations() 

    if err := micro.RegisterHandler(srv.Server(), handlers); err != nil {
		logger.Fatal(err)
	}
	
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

func cleanup(){
	eureka.Cleanup(configurations)
}

func registerRoutes(router *mux.Router) {
	registerControllerRoutes(controllers.EventController{}, router)
}

func registerControllerRoutes(controller controllers.Controller, router *mux.Router) {
	controller.RegisterRoutes(router)
}

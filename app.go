package main

import (
	. "webapiserver/internal"
	. "webapiserver/resource"

	"gitlab.bcowtech.de/bcow-go/config"
	fasthttp "gitlab.bcowtech.de/bcow-go/host-fasthttp"
)

type ResourceManager struct {
	*RootResource `url:"/"`
	*DemoResource `url:"/Demo"`
}

func main() {
	app := App{}
	fasthttp.Startup(&app,
		fasthttp.UseResourceManager(&ResourceManager{})).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadYamlFile("config.${ENVIRONMENT}.yaml").
				LoadCommandArguments()
		}).Run()
}

package config

type RESTConfigData struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

type RESTConfig struct {
	Config *RESTConfigData
}

var restConfig = RESTConfig{}

func (config *RESTConfig) CreateStructure() interface{} {
	return &RESTConfigData{
		Host: "0.0.0.0",
		Port: 5020,
	}
}

func (config *RESTConfig) Set(data interface{}) {
	config.Config = data.(*RESTConfigData)
}

func Get() *RESTConfig {
	return &restConfig
}

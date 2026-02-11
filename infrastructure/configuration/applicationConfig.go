package configuration

import "github.com/spf13/viper"

var AppConfig *appConfig

type appConfig struct {
	ApplicationPort       int
	ApplicationName       string
	ServiceName           string
	Environment           string
	XenditAPIURL          string
	XenditAPIKey          string
	XenditTimeout         int // in milliseconds
	RedisHost             string
	RedisPort             int
	RedisPassword         string
	RedisDatabase         int
	YugabyteHost          string
	YugabytePort          int
	YugabyteUsername      string
	YugabytePassword      string
	YugabyteDatabase      string
	RabbitMQURI           string
	ElasticsearchAddress  string
	ElasticsearchUsername string
	ElasticsearchPassword string
	CallbackURLNicepay    string
	ReturnURLNicepay      string
	NicepayURL            string
}

func InitializeAppConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	AppConfig = &appConfig{}
	AppConfig.ApplicationPort = viper.GetInt("APP_PORT")
	AppConfig.ApplicationName = viper.GetString("APP_NAME")
	AppConfig.ServiceName = viper.GetString("SERVICE_NAME")
	AppConfig.Environment = viper.GetString("ENV")
	AppConfig.XenditAPIURL = viper.GetString("XENDIT_API_URL")
	AppConfig.XenditAPIKey = viper.GetString("XENDIT_API_KEY")
	AppConfig.XenditTimeout = viper.GetInt("XENDIT_TIMEOUT")
	AppConfig.RedisHost = viper.GetString("REDIS_HOST")
	AppConfig.RedisPort = viper.GetInt("REDIS_PORT")
	AppConfig.RedisPassword = viper.GetString("REDIS_PASSWORD")
	AppConfig.RedisDatabase = viper.GetInt("REDIS_DATABASE")
	AppConfig.YugabyteHost = viper.GetString("YUGABYTE_HOST")
	AppConfig.YugabytePort = viper.GetInt("YUGABYTE_PORT")
	AppConfig.YugabyteUsername = viper.GetString("YUGABYTE_USERNAME")
	AppConfig.YugabytePassword = viper.GetString("YUGABYTE_PASSWORD")
	AppConfig.YugabyteDatabase = viper.GetString("YUGABYTE_DATABASE")
	AppConfig.RabbitMQURI = viper.GetString("RABBITMQ_URI")
	AppConfig.ElasticsearchAddress = viper.GetString("ELASTICSEARCH_ADDRESS")
	AppConfig.ElasticsearchUsername = viper.GetString("ELASTICSEARCH_USERNAME")
	AppConfig.ElasticsearchPassword = viper.GetString("ELASTICSEARCH_PASSWORD")
	AppConfig.CallbackURLNicepay = viper.GetString("CALLBACK_URL_NICEPAY")
	AppConfig.ReturnURLNicepay = viper.GetString("RETURN_URL_NICEPAY")
	AppConfig.NicepayURL = viper.GetString("NICEPAY_URL")
}

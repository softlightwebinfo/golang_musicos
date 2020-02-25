package settings

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/prometheus/common/log"
	"github.com/tkanos/gonfig"
	"os"
	"time"
)

var JwtKey = []byte("passSecret")
var LoginExpirationDuration = time.Duration(24) * time.Hour
var JwtSigningMethod = jwt.SigningMethodHS256
var confDB = GetConfigEnvironment()
var InstanceDb *sql.DB
var InstanceDbTX *sql.Tx

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}
type ConfigModel struct {
	ImageDir  string
	StaticDir string
	UploadDir string
	Host      string
	Web       string
}
type GrafanaConfig struct {
	User     string
	Password string
	Url      string
}
type Environment struct {
	Database DatabaseConfig `json:"database"`
}
type SmtpEmail struct {
	Email    string
	Password string
	Host     string
	Port     int
	From     string
}
type MailjetConfig struct {
	Private string
	Public  string
	Email   string
	Name    string
}
type GraphqlConfig struct {
	Client string
}
type ScrapingConfig struct {
	Urls []string `json:"urls"`
}

func GetConfigEnvironment() Environment {
	configuration := Environment{}
	file := GetEnvironmentFile()
	err := gonfig.GetConf(fmt.Sprintf("./settings/%s", file), &configuration)
	if err != nil {
		log.Fatal("Error configuration environment", err)
	}
	return configuration
}
func GetEnvironment() string {
	args := os.Args
	if len(args) > 1 {
		return args[1]
	}
	return os.Getenv("ENVIRONMENT")
}
func GetEnvironmentFile() (env string) {
	environment := GetEnvironment()
	switch environment {
	case "local":
		env = "config.local.json"
	case "production":
		env = "config.production.json"
	case "development":
		env = "config.development.json"
	case "staging":
		env = "config.staging.json"
	default:
		env = "config.local.json"
	}
	return
}
func Database() DatabaseConfig {
	return DatabaseConfig{
		User:     confDB.Database.User,
		Password: confDB.Database.Password,
		Host:     confDB.Database.Host,
		Port:     confDB.Database.Port,
		Database: confDB.Database.Database,
	}
}
func Config() ConfigModel {
	return ConfigModel{
		StaticDir: "assets",
		ImageDir:  "images",
		UploadDir: "uploads",
		Host:      "http://www.musicosdelmundo.com:8000/",
		Web:       "https://www.musicosdelmundo.com",
	}
}
func GetImage(image string) string {
	config := Config()
	return "/" + config.ImageDir + "/" + image
}
func GetGrafanaConfig() GrafanaConfig {
	return GrafanaConfig{
		User:     "user",
		Password: "pass",
		Url:      "http://localhost:3003",
	}
}
func GetSmtEmailConfig() SmtpEmail {
	return SmtpEmail{
		Email:    "info@musicosdelmundo.com",
		Host:     "io",
		Password: "pass",
		Port:     25,
		From:     "info@musicosdelmundo.com",
	}
}
func GetMailjetConfig() MailjetConfig {
	return MailjetConfig{
		Public:  "asas",
		Private: "asas",
		Email:   "info@musicosdelmundo.com",
		Name:    "Músicos del mundo",
	}
}
func GetGraphqlConfig() GraphqlConfig {
	return GraphqlConfig{
		Client: "http://localhost:8080/v1/graphql",
	}
}
func GetScrapingConfig() ScrapingConfig {
	return ScrapingConfig{
		Urls: []string{
			"https://musicosdemallorca.com/search",
		},
	}
}

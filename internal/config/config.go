package config

import (
	"log"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

type Config struct {
	Server *Server
	Client *Client
	Router *Router

	DB          *AllDB
	EBoekHouden *EBoekHouden
	Stripe      *Stripe
}

type Client struct {
	ReadTimeout     time.Duration
	ReadBufferSize  int
	WriteTimeout    time.Duration
	WriteBufferSize int
}

type Router struct {
	Current                string
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	HandleOPTIONS          bool
}

type Server struct {
	Type               string
	Name               string
	URL                string
	Port               int
	MaxRequestBodySize int
	IsDevelopment      bool
	SendToExternal     bool
}

type AllDB struct {
	Postgres *StoragePostgres
	MySQL    *StorageMySQL
	Current  string
}

type StoragePostgres struct {
	Hostname string
	User     string
	Password string
	Database string
	Port     int
}

type StorageMySQL struct {
	Hostname        string
	User            string
	Password        string
	Database        string
	Port            int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type EBoekHouden struct {
	Username                       string
	SecurityCode1                  string
	SecurityCode2                  string
	UseLedgerAccountCodeForAll     bool
	UseLedgerAccountCodePerProduct bool
	UseLedgerAccountCodePerPlan    bool
	LedgerAccountCodeDebtors       string
	LedgerAccountCodeBank          string
	LedgerAccountCodeDefault       string
	LedgerAccountCodeProducts      map[string]string
	LedgerAccountCodePlans         map[string]string
}

type Stripe struct {
	Key    string
	Secret string
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.EBoekHouden),
		validation.Field(&c.Stripe),
	)
}

func (c EBoekHouden) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Username,
			validation.Required),
		validation.Field(&c.SecurityCode1,
			validation.Required),
		validation.Field(&c.SecurityCode2,
			validation.Required),
	)
}

func (c Stripe) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Key,
			validation.Required),
		validation.Field(&c.Secret,
			validation.Required),
	)
}

func New() (config *Config) {
	log.Println("Init config")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	config = &Config{
		Server: &Server{},
		Client: &Client{},
		Router: &Router{},
		DB: &AllDB{
			Current: "mysql",
			MySQL:   &StorageMySQL{},
		},
		EBoekHouden: &EBoekHouden{
			Username:      "",
			SecurityCode1: "",
			SecurityCode2: "",
		},
		Stripe: &Stripe{
			Key:    "",
			Secret: "",
		},
	}

	if err := viper.Unmarshal(config); err != nil {
		panic("unable to decode into config struct")
	}
	return config
}

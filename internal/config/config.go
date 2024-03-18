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
	Payment     *Payment
}

type Payment struct {
	Current string
	Stripe  *Stripe
}

type Stripe struct {
	Key    string
	Secret string
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
	LedgerAccountCode    EBoekHoudenLedgerAccountCode
	Username             string
	SecurityCode1        string
	SecurityCode2        string
	UseLedgerAccountCode EBoekHoudenUseLedgerAccountCode
}

type EBoekHoudenLedgerAccountCode struct {
	Products map[string]string
	Plans    map[string]string
	Debtors  string
	Bank     string
	Default  string
}

type EBoekHoudenUseLedgerAccountCode struct {
	ForAll     bool
	PerProduct bool
	PerPlan    bool
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.EBoekHouden),
		validation.Field(&c.Payment),
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
		validation.Field(&c.UseLedgerAccountCode,
			validation.Required),
		validation.Field(&c.LedgerAccountCode,
			validation.Required),
	)
}

func (c EBoekHoudenLedgerAccountCode) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Bank,
			validation.Required,
		),
		validation.Field(&c.Debtors,
			validation.Required,
		),
		validation.Field(&c.Default,
			validation.Required,
		),
		validation.Field(&c.Products),
		validation.Field(&c.Plans),
	)
}
func (c EBoekHoudenUseLedgerAccountCode) Validate() error {
	return validation.ValidateStruct(&c)
}

func (c Payment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Current,
			validation.Required),
		validation.Field(&c.Stripe,
			validation.When(c.Current == "stripe",
				validation.Required),
		),
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
		Payment: &Payment{
			Current: "stripe",
			Stripe: &Stripe{
				Key:    "",
				Secret: "",
			},
		},
	}

	if err := viper.Unmarshal(config); err != nil {
		panic("unable to decode into config struct")
	}
	return config
}

package config

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

var Conf = new(config)

//go:embed goldap-server-priv.pem
var priv []byte

//go:embed goldap-server-pub.pem
var pub []byte

type config struct {
	System       *SystemConfig       `mapstructure:"system" json:"system"`
	Logs         *LogsConfig         `mapstructure:"logs" json:"logs"`
	Database     *Database           `mapstructure:"database" json:"database"`
	Mysql        *MysqlConfig        `mapstructure:"mysql" json:"mysql"`
	Jwt          *JwtConfig          `mapstructure:"jwt" json:"jwt"`
	RateLimit    *RateLimitConfig    `mapstructure:"rate-limit" json:"rateLimit"`
	Ldap         *LdapConfig         `mapstructure:"ldap" json:"ldap"`
	Email        *EmailConfig        `mapstructure:"email" json:"email"`
	Registration *RegistrationConfig `mapstructure:"registration" json:"registration"`
}

func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %s", err))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/")

	if err = viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("failed to unmarshal config: %s", err))
		}
		Conf.System.RSAPublicBytes = pub
		Conf.System.RSAPrivateBytes = priv
	})

	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %s", err))
	}

	Conf.System.RSAPublicBytes = pub
	Conf.System.RSAPrivateBytes = priv

	loadEnvOverrides()
}

func loadEnvOverrides() {
	if v := os.Getenv("DB_DRIVER"); v != "" {
		Conf.Database.Driver = v
	}
	if v := os.Getenv("MYSQL_HOST"); v != "" {
		Conf.Mysql.Host = v
	}
	if v := os.Getenv("MYSQL_USERNAME"); v != "" {
		Conf.Mysql.Username = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		Conf.Mysql.Password = v
	}
	if v := os.Getenv("MYSQL_DATABASE"); v != "" {
		Conf.Mysql.Database = v
	}
	if v := os.Getenv("MYSQL_PORT"); v != "" {
		Conf.Mysql.Port, _ = strconv.Atoi(v)
	}
	if v := os.Getenv("LDAP_URL"); v != "" {
		Conf.Ldap.Url = v
	}
	if v := os.Getenv("LDAP_BASE_DN"); v != "" {
		Conf.Ldap.BaseDN = v
	}
	if v := os.Getenv("LDAP_ADMIN_DN"); v != "" {
		Conf.Ldap.AdminDN = v
	}
	if v := os.Getenv("LDAP_ADMIN_PASS"); v != "" {
		Conf.Ldap.AdminPass = v
	}
	if v := os.Getenv("LDAP_USER_DN"); v != "" {
		Conf.Ldap.UserDN = v
	}
	if v := os.Getenv("LDAP_USER_INIT_PASSWORD"); v != "" {
		Conf.Ldap.UserInitPassword = v
	}
	if v := os.Getenv("LDAP_DEFAULT_EMAIL_SUFFIX"); v != "" {
		Conf.Ldap.DefaultEmailSuffix = v
	}
	if v := os.Getenv("LDAP_USER_PASSWORD_ENCRYPTION_TYPE"); v != "" {
		Conf.Ldap.UserPasswordEncryptionType = v
	}
	if os.Getenv("LDAP_ALLOW_ANON_BINDING") == "true" {
		Conf.Ldap.AllowAnonBinding = true
	}
}

type SystemConfig struct {
	Mode            string `mapstructure:"mode" json:"mode"`
	UrlPathPrefix   string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Port            int    `mapstructure:"port" json:"port"`
	InitData        bool   `mapstructure:"init-data" json:"initData"`
	RSAPublicBytes  []byte `mapstructure:"-" json:"-"`
	RSAPrivateBytes []byte `mapstructure:"-" json:"-"`
}

type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type Database struct {
	Driver string `mapstructure:"driver" json:"driver"`
	Source string `mapstructure:"source" json:"source"`
}

type MysqlConfig struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}

type JwtConfig struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type RateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
}

type LdapConfig struct {
	Url                        string `mapstructure:"url" json:"url"`
	MaxConn                    int    `mapstructure:"max-conn" json:"maxConn"`
	BaseDN                     string `mapstructure:"base-dn" json:"baseDN"`
	AdminDN                    string `mapstructure:"admin-dn" json:"adminDN"`
	AdminPass                  string `mapstructure:"admin-pass" json:"adminPass"`
	UserDN                     string `mapstructure:"user-dn" json:"userDN"`
	UserInitPassword           string `mapstructure:"user-init-password" json:"userInitPassword"`
	GroupNameModify            bool   `mapstructure:"group-name-modify" json:"groupNameModify"`
	UserNameModify             bool   `mapstructure:"user-name-modify" json:"userNameModify"`
	DefaultEmailSuffix         string `mapstructure:"default-email-suffix" json:"defaultEmailSuffix"`
	UserPasswordEncryptionType string `mapstructure:"user-password-encryption-type" json:"userPasswordEncryptionType"`
	AllowAnonBinding           bool   `mapstructure:"allow-anon-binding" json:"allowAnonBinding"`
}

type EmailConfig struct {
	Host  string `mapstructure:"host" json:"host"`
	Port  string `mapstructure:"port" json:"port"`
	User  string `mapstructure:"user" json:"user"`
	Pass  string `mapstructure:"pass" json:"pass"`
	From  string `mapstructure:"from" json:"from"`
	IsSSL bool   `mapstructure:"is-ssl" json:"isSsl"`
}

type RegistrationConfig struct {
	Mode string `mapstructure:"mode" json:"mode"`
}

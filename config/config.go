package config

type Config struct {
    Server struct{
        Host string
        Port string
    }
    Mysql struct{
        DSN string
    }
    Redis struct{
        Addr string
        Password string
        DB int
    }
    SMS struct{
        HOST string `toml:"host"`
        Sid string `toml:"sid"`
        AuthToken string `toml:"auth_token"`
        AppID string `toml:"app_id"`
        AppToken string `toml:"app_token"`
    }
}

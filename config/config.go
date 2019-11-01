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
        HOST string
        Sid string
        AuthToken string
        AppID string
        AppToken string
    }
}

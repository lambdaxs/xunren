package main

import (
    "flag"
    "fmt"
    log "github.com/golang/glog"
    "github.com/gomodule/redigo/redis"
    "github.com/jinzhu/configor"
    "github.com/jinzhu/gorm"
    "github.com/labstack/echo"
    "github.com/lambdaxs/xunren/config"
    "github.com/lambdaxs/xunren/model"
    "github.com/lambdaxs/xunren/util"
    "strings"
    "time"
)

var (
    DB        *gorm.DB
    Cache     *redis.Pool
    SMSServer *util.SMSServer
    Cfg       *config.Config
)

const (
    slat = "huzhuxunzi"
)

func Init(cfg *config.Config) {
    db, err := util.NewDB(cfg.Mysql.DSN)
    if err != nil {
        panic(err)
    } else {
        DB = db
    }
    fmt.Println("init mysql success")
    Cache = util.NewRedisPool(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
    fmt.Println("init redis success")
    SMSServer = util.NewSMSServer(cfg.SMS.HOST, cfg.SMS.Sid, cfg.SMS.AuthToken, cfg.SMS.AppID, cfg.SMS.AppToken)
    fmt.Println(cfg.SMS.HOST,cfg.SMS.Sid, cfg.SMS.AuthToken, cfg.SMS.AppID, cfg.SMS.AppToken)
}

//短信登陆
func SendLoginCode(c echo.Context) error {
    reqModel := new(struct {
        Phone     string `json:"phone"`
        Sign      string `json:"sign"`
        Timestamp int64  `json:"timestamp"`
    })
    if err := c.Bind(reqModel); err != nil {
        return util.OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    //校验签名: phone:timestamp:slat
    srvSign := util.MD5Str(fmt.Sprintf("%s:%d:%s", reqModel.Phone, reqModel.Timestamp, slat))
    if reqModel.Sign != srvSign {
        return util.OutputError(c, 1, fmt.Errorf("请求签名错误"))
    }
    code := util.CaptchaCode()
    if err := SMSServer.SendLoginCode(reqModel.Phone, code, "483183"); err != nil {
        return util.OutputError(c, 2, fmt.Errorf("短信发送失败:%s", err.Error()))
    }
    //存入redis
    conn := Cache.Get()
    defer conn.Close()

    key := fmt.Sprintf("login_code_%s", reqModel.Phone)
    _, setErr := conn.Do("SET", key, code)
    if setErr != nil {
        return util.OutputError(c, 2, fmt.Errorf("server error"+setErr.Error()))
    }
    _, expErr := conn.Do("EXPIRE", key, 300)
    if expErr != nil {
        return util.OutputError(c, 2, fmt.Errorf("server error"+expErr.Error()))
    }
    return util.OutputData(c, 0, true, "success")
}

//用户登陆
func UserLogin(c echo.Context) error {
    reqModel := new(struct {
        Phone string `json:"phone"`
        Code  string `json:"code"`
    })
    if err := c.Bind(reqModel); err != nil {
        return util.OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }

    if strings.TrimSpace(reqModel.Phone) == "" {
        return util.OutputError(c, 1, fmt.Errorf("手机号不能为空"))
    }
    if strings.TrimSpace(reqModel.Code) == "" {
        return util.OutputError(c, 1, fmt.Errorf("验证码不能为空"))
    }

    //查询redis手机号验证码
    conn := Cache.Get()
    defer conn.Close()

    key := fmt.Sprintf("login_code_%s", reqModel.Phone)
    smsCode, err := redis.String(conn.Do("GET", key))
    if err != nil {
        if err == redis.ErrNil {
            err = fmt.Errorf("验证码已过期")
            return util.OutputError(c, 2, err)
        }
        return util.OutputError(c, 2, err)
    }
    if smsCode != reqModel.Code {
        err = fmt.Errorf("验证码错误")
        return util.OutputError(c, 2, err)
    }
    //验证成功，删除验证码
    if _, delErr := conn.Do("DEL", key); delErr != nil {
        log.Error("del login code error:" + delErr.Error())
    }
    //判断用户是否注册过
    now := time.Now().Unix()
    userData := model.User{}
    DB.Where("phone = ?", reqModel.Phone).First(&userData)
    if userData.ID == 0 { //新用户
        userData = model.User{
            Phone:    reqModel.Phone,
            Name:     fmt.Sprintf("尾号%s", reqModel.Phone[:7]),
            CreateAt: now,
            UpdateAt: now,
            IsNew:    true,
        }
        if saveErr := DB.Create(&userData); saveErr != nil {
            return util.OutputError(c, 2, fmt.Errorf("注册新用户失败"))
        }
        return util.OutputData(c, 0, userData, "")
    }
    //已经注册过,直接登陆
    return util.OutputData(c, 0, userData, "")
}

//发布
func InfoPublish(c echo.Context) error {
    reqModel := new(struct {
        Uid     int64    `json:"uid" form:"uid"`
        Title   string   `json:"title"`
        Content string   `json:"content"`
        Images  []string `json:"images"`
    })
    if err := c.Bind(reqModel); err != nil {
        return util.OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    info := model.Info{
        Uid:      reqModel.Uid,
        Title:    reqModel.Title,
        Content:  reqModel.Content,
        Images:   strings.Join(reqModel.Images, ","),
        CreateAt: time.Now().Unix(),
    }
    if err := DB.Create(info); err != nil {
        return util.OutputError(c, 2, fmt.Errorf("发布信息失败:%s", err.Error))
    }
    return util.OutputData(c, 0, info.ID, "")
}

//发布列表
func InfoList(c echo.Context) error {
    reqModel := new(struct {
        Uid    int64 `json:"uid" form:"uid"`
        InfoId int64 `json:"info_id"`
        Limit  int32 `json:"limit"`
    })
    if err := c.Bind(reqModel); err != nil {
        return util.OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    list := make([]model.Info, 0)
    if reqModel.InfoId == 0 {
        DB.Order("id desc").Limit(reqModel.Limit).Find(&list)
    } else {
        DB.Where("id < ?", reqModel.InfoId).Order("id desc").Limit(reqModel.Limit).Find(&list)
    }
    for i,item := range list {
        list[i].ImageList = strings.Split(item.Images,",")
    }
    return util.OutputData(c, 0, list, "")
}

//发布详情
func InfoDetail(c echo.Context) error {
    reqModel := new(struct {
        Uid    int64 `json:"uid" form:"uid"`
        InfoId int64 `json:"info_id"`
    })
    if err := c.Bind(reqModel); err != nil {
        return util.OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    infoData := model.Info{}
    DB.Where("id = ?").First(&infoData)
    infoData.ImageList = strings.Split(infoData.Images,",")
    return util.OutputData(c, 0, infoData, "")
}

//我发布的列表
func MyInfoList(c echo.Context) error {
    reqModel := new(struct {
        Uid    int64 `json:"uid" form:"uid"`
        InfoID int64 `json:"info_id"`
        Limit  int32 `json:"limit"`
    })
    if err := c.Bind(reqModel); err != nil {
        return util.OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    list := make([]model.Info, 0)
    DB.Where("uid = ?",reqModel.Uid).Order("id desc").Limit(reqModel.Limit).Find(&list)
    for i,item := range list {
        list[i].ImageList = strings.Split(item.Images,",")
    }
    return util.OutputData(c, 0, list, "")
}

func StartServer(addr string) {
    srv := echo.New()

    srv.POST("/api/v1/sms/send.json", SendLoginCode)
    srv.POST("/api/v1/user/login.json", UserLogin)
    srv.POST("/api/v1/info/publish.json", InfoPublish)
    srv.POST("/api/v1/info/detail.json", InfoDetail)
    srv.POST("/api/v1/info/list.json", InfoList)
    srv.POST("/api/v1/info/mylist.json", MyInfoList)

    if err := srv.Start(addr); err != nil {
        fmt.Println(err.Error())
        return
    }
}

func main() {
    var configPath string
    flag.StringVar(&configPath, "config", "", "config file path")
    flag.Parse()
    defer log.Flush()

    //加载配置
    Cfg = &config.Config{}
    if err := configor.Load(Cfg, configPath); err != nil {
        panic(err)
    }
    //初始化资源
    Init(Cfg)

    //启动服务器
    StartServer(fmt.Sprintf("%s:%s", Cfg.Server.Host, Cfg.Server.Port))
}

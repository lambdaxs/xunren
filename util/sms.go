package util

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

type SMSServer struct {
    Host      string
    Sid       string
    AuthToken string
    AppID     string
    AppToken  string
    Client    http.Client
}

type SMSContent struct {
    To         string   `json:"to"`
    AppID      string   `json:"appId"`
    TemplateID string   `json:"templateId"`
    Datas      []string `json:"datas"`
}

func (c *SMSContent) JSONBuffer() []byte {
    buf, _ := json.Marshal(c)
    return buf
}

type SMSResponse struct {
    StatusCode string `json:"statusCode"` //"000000"
}

func NewSMSServer(host, sid, authToken, appID, appToken string) *SMSServer {
    client := http.Client{
        Timeout: time.Second * 3,
    }
    return &SMSServer{
        Host:      host,      //"https://app.cloopen.com:8883",
        Sid:       sid,       //"8aaf0708670aef8401670cd0f7820567",
        AuthToken: authToken, //"227b7798396e42cebcb2c07f9f690002",
        AppID:     appID,     //"8aaf0708670aef8401670cd0f7df056d",
        AppToken:  appToken,  //"0276a86be9dbccbf98c521ed19f9466e",
        Client:    client,
    }
}

//登陆验证短信
func (s *SMSServer) SendLoginCode(phone string, code string, tempID string) error {
    return s.Send([]string{phone}, []string{code, "2"}, tempID)
}

func (s *SMSServer) Send(to []string, datas []string, templateID string) error {
    timestamp := time.Now().Format("20060102150405")
    sig := strings.ToUpper(MD5Str(fmt.Sprintf("%s%s%s", s.Sid, s.AuthToken, timestamp)))
    content := SMSContent{
        To:         strings.Join(to, ","),
        AppID:      s.AppID,
        TemplateID: templateID,
        Datas:      datas,
    }
    req, err := http.NewRequest("POST",
        fmt.Sprintf("%s/2013-12-26/Accounts/%s/SMS/TemplateSMS?sig=%s", s.Host, s.Sid, sig),
        bytes.NewBuffer(content.JSONBuffer()))
    if err != nil {
        return err
    }
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json;charset=utf-8")
    req.Header.Set("Authorization", Base64Str(fmt.Sprintf("%s:%s", s.Sid, timestamp)))
    resp, err := s.Client.Do(req)
    if err != nil {
        return err
    }
    buf, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    resp.Body.Close()
    fmt.Println(string(buf))

    result := SMSResponse{}
    if err := json.Unmarshal(buf, &result); err != nil {
        return err
    }
    if result.StatusCode == "000000" {
        return nil
    }
    return fmt.Errorf("短信发送失败:%s", result.StatusCode)
}

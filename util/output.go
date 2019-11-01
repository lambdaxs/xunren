package util

import (
    "github.com/labstack/echo"
    "net/http"
)

type Response struct {
    Code int `json:"code"`
    Data interface{} `json:"data"`
    Msg string `json:"msg"`
}

func OutputData(c echo.Context, status int, data interface{}, msg string) error {
    return c.JSON(http.StatusOK, Response{
        Code: status,
        Data: data,
        Msg:  msg,
    })
}

func OutputError(c echo.Context, status int , err error) error {
    return c.JSON(http.StatusOK, Response{
        Code: status,
        Data: nil,
        Msg:  err.Error(),
    })
}

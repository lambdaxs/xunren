package model

type User struct {
    ID       int64  `json:"id"`
    Phone    string `json:"phone"`
    Name     string `json:"name"`
    Avatar   string `json:"avatar"`
    Token    string `json:"token"`
    CreateAt int64  `json:"create_at"`
    UpdateAt int64  `json:"update_at"`
    IsNew    bool   `json:"is_new" gorm:"-"`
}

func (u *User) TableName() string {
    return "user"
}

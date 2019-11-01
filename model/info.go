package model

type Info struct {
    ID int64 `json:"id"`
    Uid int64 `json:"uid"`
    Title string `json:"title"`
    Content string `json:"content"`
    Images string `json:"image"`
    CreateAt int64 `json:"create_at"`
}

func (i *Info)TableName() string {
    return "info"
}

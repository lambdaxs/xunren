# 服务端API

## 发送验证码

#### url: /api/v1/sms/send.json
#### params: 

字段 | 类型 | 必填 | 含义
--- | --- | --- | ---
phone | string | 是 | 手机号
timestamp | int64 | 是 | 时间戳
sign | string | 是 | 签名

#### result:

```json
{
  "code":0,
  "data": true,
  "msg": "success" 
}
```

## 验证码登陆

#### url: /api/v1/user/login.json

#### params:

字段|类型|必填|含义
--- | --- | --- | ---
phone | string | 是 | 手机号
code | string | 是 | 验证码

#### result:

```json
{
  "code": 0,
  "data": {
    "uid": 100,
    "name": "尾号5086",
    "phone": "18686535086",
    "avatar": ""
  },
  "msg": "success"
}
```

## 发布寻人信息

#### url: /api/v1/info/publish.json

#### params:

字段|类型|必填|含义
--- | --- | --- | ---
uid | int64 | 是 | 用户id
title | string | 是 | 标题
content | string | 是 | 内容
images  | []string | 否 | 图片url数组

#### result

```json
{
  "code": 0,
  "data": 1232,
  "msg": "success"
}
```

## 寻人列表

#### url: /api/v1/info/list.json

#### params:

字段|类型|必填|含义
--- | --- | --- | ---
uid | int64 | 是 | 用户id
info_id | int64 | 是 | 上一条信息流id,获取最新的传0
limit | int | 是 | 一次获取的条数

#### result

```json
{
  "code": 0,
"data": [{
  "id": 131223,
  "title": "我儿乔峰",
  "content": "今日于加州走失...",
  "image_list": ["https://1.jpg","https://2.jpg"],
  "create_at": 154324454455
}]
}
```

## 寻人信息详情

#### url: /api/v1/info/detail.json

#### params:

字段|类型|必填|含义
--- | --- | --- | ---
uid | int64 | 是 | 用户id
info_id | int64 | 是 | 信息流id

#### result

```json
{
  "code": 0,
  "data": {
    "id": 131223,
    "title": "我儿乔峰",
    "content": "今日于加州走失...",
    "image_list": ["https://1.jpg","https://2.jpg"],
    "create_at": 154324454455
  },
  "msg": "success"
}
```

## 我的发布列表

#### url: /api/v1/info/mylist.json

#### params:

字段|类型|必填|含义
--- | --- | --- | ---
uid | int64 | 是 | 用户id
info_id | int64 | 是 | 上一条信息流id,获取最新的传0
limit | int | 是 | 一次获取的条数

#### result

```json
{
  "code": 0,
  "data": [{
    "id": 131223,
    "title": "我儿乔峰",
    "content": "今日于加州走失...",
    "image_list": ["https://1.jpg","https://2.jpg"],
    "create_at": 154324454455
  }],
  "msg": "success"
}
```
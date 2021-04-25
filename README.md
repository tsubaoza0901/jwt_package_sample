# jwt_package_sample

# 使い方

## 1. アクセス時トークン不要な URL

① トークンなしでアクセスできる URL

POST http://127.0.0.1:9010/public

```
{
    "name": "xxxxxxx"
}
```

② user 登録およびトークン取得のための URL

POST http://127.0.0.1:9010/signup

```
{
    "name": "admin",
    "password": "password"
}
```

※signup に POST した際の Response としてトークンを取得できる。トークンが必要な URL にアクセスする際には、そのトークンを header として設定してアクセスする。

## 2. アクセス時トークンが必要な URL

POST http://127.0.0.1:9010/api/private

```
{
    "name": "xxxxxxx"
}
```

※signup で取得したトークンを header に設定していないとアクセスできない。

# 参考 URL

Godoc | package jwt  
https://godoc.org/github.com/dgrijalva/jwt-go

Qiita | Go 言語で理解する JWT 認証 実装ハンズオン
https://qiita.com/po3rin/items/740445d21487dfcb5d9f

Qiita | Go 言語で Echo を用いて認証付き Web アプリの作成
https://qiita.com/x-color/items/24ff2491751f55e866cf

Qiita | Golang の Echo で JWT
https://qiita.com/kiyc/items/959283ff84be99c42ad6

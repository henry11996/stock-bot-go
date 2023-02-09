此為 Telegram, Discord 查詢股價 bot

# 部署

必要Key
- fugle api token
- telegram api token or discord token

# 可使用指令

以 台積電 2330 為範例

### 顯示即時五檔
```
Telegram
/2330
/台積電

Discord
/tw 2330
/tw 台積電
```

### 顯示股票資訊
```
Telegram
/2330 i
/台積電 i

Discord
/tw 2330 i
/tw 台積電 i
```

### 顯示三大法人買賣超日/月報表
```
Telegram
#當日
/2330 d
/台積電 d

#當月
/2330 m
/台積電 m

#2021/05/01
/2330 d 20210501
/台積電 d 20210501

#2021/04
/2330 m 20210401
/台積電 m 20210401

Discord
#當日
/tw 2330 d
/tw 台積電 d

#當月
/tw 2330 m
/tw 台積電 m

#2021/05/01
/tw 2330 d 20210501
/tw 台積電 d 20210501

#2021/04
/tw 2330 m 20210401
/tw 台積電 m 20210401
```

# 開發

開啟 ngrok
``` sh
ngrok http 80
```

執行
```
go run .
```
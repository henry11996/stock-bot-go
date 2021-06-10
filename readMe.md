此為 Telegram 查詢股價 bot

# 部署

必要Key
- fugle api token
- telegram api token

部署至herko
https://devcenter.heroku.com/articles/getting-started-with-go#deploy-the-app

記得在 herkuo 的 Config Vars 加入
- FUGLE_API_TOKEN
- TELEGRAM_APITOKEN
- HEROKU_APP_NAME

# 可使用指令

以 台積電 2330 為範例

### 顯示即時五檔
```
/2330
/台積電
```

### 顯示股票資訊
```
/2330 i
/台積電 i
```

### 顯示三大法人買賣超日/月報表
```
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
```

# 開發中項目
- 即時圖表 使用 plot 套件
- 排成通知三大法人總買賣超金額
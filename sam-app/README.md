# AWS SAM サンプル

* API Gateway + Lambda/go + DynamoDBで、いろいろgoのライブラリを試しているサンプル（作成中）
## 構成イメージ
![構成イメージ](image/demo.drawio.png)

## SAM ビルド
```sh
cd sam-app
sam build
# viper用のconfigを.aws-samフォルダにコピー（Windowsでのxcopyの例）
xcopy /I config .aws-sam\build\GetUsersFunction\config
xcopy /I config .aws-sam\build\PostUsersFunction\config
```
## SAM デプロイ
```sh
sam deploy --guided
# 2回目以降は
sam deploy
```
## 実行
* curlコマンドを実行
```sh
# POSTコマンドの例
curl -X POST -H "Content-Type: application/json" -d '{ "name" : "Taro"}' https://42b4c7bk9g.execute-api.ap-northeast-1.amazonaws.com/Prod/users

# Getコマンドの例（users/の後にPOSTコマンドで取得したユーザIDを指定）
curl https://civuzxdd14.execute-api.ap-northeast-1.amazonaws.com/Prod/users/d4d6cb7f-7691-11ec-9520-1ee887dd490e
```
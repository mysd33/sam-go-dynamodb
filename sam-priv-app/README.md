# Private API APIGatewayサンプルの下準備

## 1. VPCおよびサブネット、InternetGateway等の作成
```sh
cd cfn
aws cloudformation validate-template --template-body file://cfn-vpc.yaml
aws cloudformation create-stack --stack-name Demo-VPC-Stack --template-body file://cfn-vpc.yaml
```

## 3. EC2(Bastion)用Security Groupの作成
```sh
aws cloudformation validate-template --template-body file://cfn-sg.yaml
aws cloudformation create-stack --stack-name Demo-SG-Stack --template-body file://cfn-sg.yaml
```
## 2. EC2(Basion)の作成
```sh
aws cloudformation validate-template --template-body file://cfn-bastion-ec2.yaml
aws cloudformation create-stack --stack-name Demo-Bastion-Stack --template-body file://cfn-bastion-ec2.yaml
```
* CloudFormationの出力「BastionDNSName」に表示されるドメイン名の値でアクセスできます。

## 3. NAT Gatewayの作成とプライベートサブネットのルートテーブル更新
```sh
aws cloudformation validate-template --template-body file://cfn-ngw.yaml
aws cloudformation create-stack --stack-name Demo-NATGW-Stack --template-body file://cfn-ngw.yaml
```
## 4. AWS SAMでLambda/API Gatewayの実行
* SAMビルド
```sh
cd ../sam-priv-app
sam build
```
* 必要に応じてローカル実行可能
```sh
sam local invoke
sam local start-api
curl http://127.0.0.1:3000/hello
```
* SAMデプロイ
```sh
sam deploy --guided
#2回目以降は
sam deploy
```

* EC2(Bation)へTeraTermでログインして、動作確認
```sh
curl https://5h5zxybd3c.execute-api.ap-northeast-1.amazonaws.com/Prod/hello/
```
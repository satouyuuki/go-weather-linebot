# 社畜専用お天気通知bot
毎朝7時に傘を持っていくべきかどうかを自動で通知してくれるbot


### インフラ構成図
![インフラ構成図](assets/weather-bot.drawio.png)

### aws_cliのプロファイルをセットする方法
```
$ export TF_VAR_AWS_PROFILE=xxxxx
```

### dockerの基本的な使い方(debug)
```
# imageからcontainerを作成して起動する
$ docker run -p 9000:8080 ${container_image_name}
# エンドポイント-> http://localhost:9000/2015-03-31/functions/function/invocations

# コンテナを停止
$ docker stop ${contaienrid}

# コンテナをスタート
$ docker start ${contaienrid} -a

# コンテナを更新する流れ(全てhost側でcommand実施)
### build
$ GOOS=linux GOARCH=amd64 go build -o main
### container 起動
$ docker start ${contaienrid}
### 実行ファイルをhostからcontainerにコピー
$ docker cp main ${containerid}:/ 
### container 再起動
$ docker restart ${contaienrid}
```


### 参考になったサイト
https://hands-on.cloud/terraform-deploy-python-lambda-container-image/
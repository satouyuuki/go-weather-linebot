```
# まずECRを作成する
$ terraform apply -target="aws_ecr_repository.weather_bot_repo"

# 次にlambda container imageをecrにpushする
$ docker build -t weather-lambda . --platform=linux/amd64

$ aws ecr get-login-password --region ap-northeast-1 --profile cli_only_user | docker login --username AWS --password-stdin ${aws_account_id}.dkr.ecr.ap-northeast-1.amazonaws.com

$ docker tag ${image_id} ${aws_account_id}.dkr.ecr.ap-northeast-1.amazonaws.com/weather_bot_repo:0.0.1

$ docker push ${aws_account_id}.dkr.ecr.ap-northeast-1.amazonaws.com/weather_bot_repo:latest

# 最後に残りのリソースをデプロイする
$ terraform apply 
```


# 毎日渋谷の天気を通知するLINEbot

# ローカルでlambdaをテストする
```
$ docker run -p 9000:8080 weather-lambda
```

エンドポイント: 
`http://localhost:9000/2015-03-31/functions/function/invocations`

### aws_cliのプロファイルをセットする
```
$ export TF_VAR_AWS_PROFILE=xxxxx
```

### 参考になったサイト
https://hands-on.cloud/terraform-deploy-python-lambda-container-image/
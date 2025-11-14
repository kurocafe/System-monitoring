# System-monitoring
GO!

# 環境起動方法
1. 初回のみ
 ```
 docker compose up --build
 ```
2. 初回以降
```
docker compose start
```

3. 環境コンテナ削除
```
docker compose down
```

4. 環境コンテナ一時停止
```
docker compose stop
```

# Go の走らせ方

 goファイルが存在するディレクトリに移動
 ```
 cd <directory name>
 ```
 Go ファイルを走らせる
 ```
 go run <go file name>
 ```  
 # 複数のファイルがある場合
 ```
 go run <file 1> <file 2> <file 3>
 ```

 ### 簡略化したいならモジュールを作成してから走らせる
 ```
 go mod init <random name ok> && go run .
 ```

 # 良いファイル構造(according to chat gpt)
 ```
 project/
    cmd/
        myapp/
            main.go
    internal/
        utils/
        game/
        ai/
    pkg/
        math/

 ```


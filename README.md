# JarvisGo

A Golang Version of Jarvis bot.

## Run Jarvis

```bash
# Under JarvisGo folder
$ go mod tidy
$ go run cmd/server/main.go
```

Also, you can just run `./start.bat` if you like:)

## Config Jarvis

You might notice that you will not be able to start Jarvis successfully until you correctly configure it.

We have two configuration files to edit.

* `config.yml`: We are using **go-cqhttp** as our login tools, please refer to [go-cqhttp configuration](https://docs.go-cqhttp.org/guide/config.html#%E9%85%8D%E7%BD%AE%E4%BF%A1%E6%81%AF) for more. Treat `config.yml.example` as an good example.

* `JarvisConfig.yml`: This is the configuration file for Jarvis bot, the `JarvisConfig.yml.example` is an good example.

  ```yaml
  # Set this to true if you would like Jarvis to respond to group message
  enable-group: true
  
  # Administrators of Jarvis. Please be careful when you set this value. A good and simple value for this field is your **own** qquid (Replace 12345678 with your own qquid).
  masters:
    - 12345678
  
  # Working directory of Jarvis. This field determines where `jdata/` and `jlog/` is placed. **Do not** use relative path, please make sure you are using absolute path. If you don't know what is "path", leave this field empty, like `working-directory: `, and follow the steps in "Run Jarvis".
  working-directory: /home/abdtyx/code/JarvisGo/
  ```

  
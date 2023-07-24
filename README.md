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

  

## Structure of JarvisGo (updated before 0d73e60)

```yaml
.
├── cmd
│   └── server
│       └── main.go           # main function of Jarvis
├── config
│   └── config.go             # load configuration of Jarvis (JarvisConfig.yml)
├── config.yml.example        # go-cqhttp configuration example
├── errors
│   └── errors.go             # placeholder, for future uses
├── go.mod                    # go modules
├── go.sum                    # go summary
├── handler
│   └── handler.go            # message handler. Receive post from gin, then resolve which service will be called
├── JarvisConfig.yml.example  # Configuration of Jarvis
├── message
│   └── message.go            # helper functions to send private and group message
├── README.md                 # PLEASE README
├── service
│   ├── service.go            # all services provided by Jarvis, including helper functions used by service functions
│   └── timers.go             # timed message, which means message under this classification will be sent at a specific time
└── start.bat                 # batch file for windows, also works for linux

7 directories, 13 files
```

## TODO List (In no particular order)

* ~~Deploy docker containers, adopting agile software development.~~ Done.

* Abstract Jarvis kernel, using Trie tree to optimize keyword detection. Jarvis kernel is supposed to be a skeleton bot server. Users can register their response functions using kernel methods.
* Use MySQL to enrich Jarvis Corpus.
* Build detailed tests to make sure JarvisGo behaves as expected. Adopting Ginkgo.
* Activation manager for XJTUANA.

## Author's words

**`JarvisGo`** is an optimized and refactored version of my previous project `Jarvis (Python)`. For some reasons, `Jarvis` is not open to the public. Several months ago, after careful consideration, I decided to refactor `Jarvis` to make it faster, better, and more comprehensive. And this is what you see now, **`JarvisGo`**! Enjoy **`Jarvis`**!

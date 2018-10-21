# audigo
3D Led CubeのPCレス(Raspberry pi)音響サービス  

<!-- toc -->  
* [💊  Requirements](#-requirements)
* [📌 Installing](#-installing)
* [🎧  Usage](#-usage)
* [🌏  REST Api](#-rest-api)
* [🎃  Notes](#-notes)
<!-- tocstop -->  

# Getting Started
## 💊 Requirements

**ALL**  
* git
* dep
    ```sh
    $ go get -u github.com/golang/dep/cmd/dep
    ```
* Go 1.11 or later

**Linux**  
```sh
$ sudo apt install libasound2-dev
```
  
  
## 📌 Installing

1. Goto GOPATH  
    **WIndows**
    ```sh
    $ cd %GOPATH%
    ```

    **Others**
    ```sh
    $ cd $HOME/go
     or
    $ cd $GOPATH
    ```

2. Get src
    ```sh
    $ git clone https://github.com/YGFYHD2018/audigo.git ./src/github.com/YGFYHD2018/audigo
    $ cd ./src/github.com/YGFYHD2018/audigo
    $ dep ensure
    ```

3. Build
    ```sh
    $ go build
    ```

# 🎧 Usage
Start audio service.  
```sh
$ audigo
```

## 🔨 Options

### port
add port number. default port 80

    ```sh
    Listening port 5701
    $ audigo 5701

    Listening port 8080
    $ audigo 8080

    Listening port 80
    $ audigo
    ```

## 📖 help
```sh
NAME:
   audigo - Audio service by LED CUBU

USAGE:
   audigo.exe [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     server, s  Instant server mode. (default)
     client, c  Instant client REPL mode.
     repl, r    Instant local REPL mode.
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cd value     change current directory by repl
   --help, -h     show help
   --version, -v  print the version

client OPTIONS:
   --cd value                change current directory by repl
   --domain value, -d value  set request domain url by client (default: "http://audigo.local")
```

## 📂 Directory layout

Add sound file in audigo/asset/audio

```sh
audigo
 |-- audigo
 |-- asset
      |-- audio
           |-- bgm_wave.wav
           |-- hogehoge.mp3
           |-- foobar.wav
           |-- ...

```


# 🌏️ REST Api
| REST | URI                             | note                          | arguments (json)     |
|------|---------------------------------|-------------------------------|----------------------|
| GET  | /audio/v1/ping                  | I Can Fly !                   | none                 |
| POST | /audio/v1/init/\<content id>    | init players in memory        | none                 |
| POST | /audio/v1/play/\<content id>    | play sound                    | src: "bgm_wave.wav"<br> (file name in ./asset/audio/) <br><br>loop: true or false<br> (loop play or single play) <br><br>stop: true or false<br> (start and pause or normal play)        |
| POST | /audio/v1/stop/\<content id>    | stop content player sound     | none                 |
| POST | /audio/v1/pause/\<content id>   | pause content player sound    | none                 |
| POST | /audio/v1/resume/\<content id>  | resume content player sound   | none                 |
| POST | /audio/v1/volume/\<content id>  | change volume                 | vol: 2 (0 - n, 0 is silent)          |


# 🌏 Api
| REST | URI                             | note                          | arguments     |
|------|---------------------------------|-------------------------------|---------------|
| GET  | /audio/v1/ping                  | I Can Fly !                   | none          |
| POST | /audio/v1/init/\<content id>    | init players in memory        | none          |
| POST | /audio/v1/play/\<content id>    | play sound                    | src: "bgm_wave.wav" (file name in ./asset/audio/) <br>loop: true or false (loop play or single play) <br>stop: true or false (start and pause or normal play)        |
| POST | /audio/v1/stop/\<content id>    | stop content player sound     | none          |
| POST | /audio/v1/pause/\<content id>   | pause content player sound    | none          |
| POST | /audio/v1/resume/\<content id>  | resume content player sound   | none          |
| POST | /audio/v1/volume/\<content id>  | change volume                 | vol: 2 (0 - n, 0 is silent)          |



# 🎃 Notes

| Platform / Architecture        | x86 | x86_64 |
|--------------------------------|-----|--------|
| Windows (7, 10 or Later)       | -   | ✓     |
| Rasbian (STRETCH or Later)     | ✓  | -      |
| OSX (10.14 or Later)           | -   | ✓     |


以上  
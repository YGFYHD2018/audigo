# audigo
3D Led CubeのPCレス(Raspberry pi)音響サービス  

<!-- toc -->  
* [💊  Requirements](#-requirements)
* [📌 Installing](#-installing)
* [🎧  Usage](#-usage)
* [🎃  Notes](#-notes)
<!-- tocstop -->  

# Getting Started
## 💊 Requirements

* git
* Go 1.11 or later

## 📌 Installing

1. Goto GOPATH  
    **WIndows**
    ```sh
    $ cd %GOPATH%/go
    ```

    **Others**
    ```sh
    $ cd $HOME/go
    ```

2. Get src
    ```sh
    $ git clone https://github.com/code560/audigo.git ./src/github.com/code560/audigo
    $ cd ./src/github.com/code560/audigo
    $ dep ensure
    ```

# 🎧 Usage
Start audio service.
```sh
$ go run audigo.go 80
```

## 🔨 Commands

### port
add port number. default port 8080


# 🎃 Notes

| Platform / Architecture        | x86 | x86_64 |
|--------------------------------|-----|--------|
| Windows (7, 10 or Later)       |     | ✓     |
| Rasbian (?)                    |     |        |
| OSX (?)                        |     |        |


以上  
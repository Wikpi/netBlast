# netBlast - terminal Chat App
netblast is a simple terminal Chat App written in Go

![image](https://user-images.githubusercontent.com/66695611/235474854-641b0229-f833-4c35-a3be-c7593da87641.png)

**Client**

The whole Client side is based on [Bubble Tea](https://github.com/charmbracelet/bubbletea), a go framework based on [The Elm Architecture](https://guide.elm-lang.org/architecture/).
While that helps to create a nice and cozy frontend, the [nhooyr/websocket](https://github.com/nhooyr/websocket) establishes a connection and helps to communicate with the server.

**Server**

The Server side is a simple http server host, which is further upgraded by [nhooyr/websocket](https://github.com/nhooyr/websocket), that helps to enhance the whole expereince.

* Huge inspirtaion derived from [tiny](https://github.com/osa1/tiny), a simple terminal IRC client.

# Features
**Clean UI**
* Chat room is only shown to registered users
* Messages are accompanied by a timestamp and seperated by full date, if they were written on a different day
* Every name is unique with randomly picked color to stand out from the others.

**Error handling**
* All errors are checked and recorded in respective logs files

# Running / Installation
Tested on Windows and Linux.

To run, simply launch the `server/main.exe` file or run in command line:  
>go run server/main.go

Afterwards clients can join either by running `client/main.exe` or in command line:
>go run clinet/main.go

![2023-05-01 18-56-32](https://user-images.githubusercontent.com/66695611/235483409-93815da2-ae86-4116-bdf8-f9f40704745d.gif)


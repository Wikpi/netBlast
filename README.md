# netBlast - terminal Chat App
![logo](https://user-images.githubusercontent.com/66695611/236682936-765b5685-9b16-4b08-91a9-29c0e8bc7cd3.png)

netblast is a simple terminal Chat App written in Go

**Client**

The whole Client side is based on [Bubble Tea](https://github.com/charmbracelet/bubbletea), a go framework based on [The Elm Architecture](https://guide.elm-lang.org/architecture/).
While that helps to create a nice and cozy frontend, the [nhooyr/websocket](https://github.com/nhooyr/websocket) establishes a connection and helps to communicate with the server.
Also uses the [Autolycus](https://github.com/Wikpi/Autolycus) module to scrape unique colors for the users.

**Server**

The Server side is a simple http server host, which is further upgraded by [nhooyr/websocket](https://github.com/nhooyr/websocket), that helps to enhance the whole expereince.

* Huge inspirtaion derived from [tiny](https://github.com/osa1/tiny), a simple terminal IRC client.

#
![chat](https://user-images.githubusercontent.com/66695611/236644775-e5403f6f-0983-4ef3-a36a-2613732195d5.png)

# Features
**Clean UI**
* Chat room is only shown to registered users.
* Messages are accompanied by a timestamp and seperated by full date, if they were written on a different day.
* Every name is unique and has a randomly picked color to stand out from the others ([Autolycus](https://github.com/Wikpi/Autolycus)).

**Error handling**
* All errors are checked and recorded in respective logs files.

**Clean layout**

# Running / Installation
To view the project locally, clone the repository:
```
git clone https://github.com/Wikpi/netBlast
```
(Similar documentation can be found in the `./docs` sub folder.)

Afterwards to make the executables, which you can find and launch seperatly in the `./build` (if it is not present, the directory will be automatically created), run:
```
make all
```

Or run the server in your cmd:
```
make server
```

After which, launch your client:
```
make client
```

# Project layout
**cmd** - holds all the main applications of the project, `./server` for server side scripts and `./client` for client.

**pkg** - stores additional packages used throughout the project.

**build** - path, in which the build executables are made and stored, `./server` for server side app and `./client` for client.

**assets** - random assortment of media, used for decorating the project page.

**docs** - short documentation of the project.

#

![2023-05-01 18-56-32](https://user-images.githubusercontent.com/66695611/235483409-93815da2-ae86-4116-bdf8-f9f40704745d.gif)

# ***Still being updated***


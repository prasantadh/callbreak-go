This software is currently in alpha version. Program might crash
for players and there might be api breaking changes for developers.
Issues and pull requests are welcome.

## Installation And Usage

This repo depends on golang so you will need golang installed.
Refer to the [official golang repository](https://go.dev/doc/install) 
for installation based on your OS. 

### Linux/MacOS
```bash
## install the package
go install github.com/prasantadh/callbreak-go/cmd/callbreak-go@latest
## run the server
~/go/bin/callbreak-go server
## on a different terminal window, use the client to connect and play
~/go/bin/callbreak-go client
```

### Windows CMD
```powershell
## install the package
go install github.com/prasantadh/callbreak-go/cmd/callbreak-go@latest
## run the server
%HOMEPATH%\go\bin\callbreak-go server
## on a different terminal window, use the client to connect and play
%HOMEPATH%\go\bin\callbreak-go client
```


## Rules
There are primarily two moves: call and break. The detailed rules are 
documented in [docs/rules.md](docs/rules.md).

## Credits
This project is a go implementation of initial work with python at
[callbreak](https://github.com/prasantadh/callbreak)

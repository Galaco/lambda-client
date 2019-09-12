
# Lambda Client
> Lambda Client is a game engine written in golang designed that loads Valve's Source Engine projects. 

[![GoDoc](https://godoc.org/github.com/Galaco/lambda-client?status.svg)](https://godoc.org/github.com/Galaco/lambda-client)
[![Go report card](https://goreportcard.com/badge/github.com/galaco/lambda-client)](https://goreportcard.com/badge/github.com/galaco/lambda-client)
[![GolangCI](https://golangci.com/badges/github.com/galaco/lambda-client.svg)](https://golangci.com)
[![codecov](https://codecov.io/gh/Galaco/lambda-client/branch/master/graph/badge.svg)](https://codecov.io/gh/Galaco/lambda-client)
[![CircleCI](https://circleci.com/gh/Galaco/lambda-client.svg?style=svg)](https://circleci.com/gh/Galaco/lambda-client)

The end goal is to be able to point this application at a source engine game and be able to
load and play that games levels. Where this progresses beyond that, needs to be decided. Most likely this will be come either a thin client for multiple
source games with game specific code layered on top (target multiplayer as priority), or the full server simulation for single player games
would be written (targeting single player as priority).

![de_dust2](https://cdn.galaco.me/github/lambda-client/readme/de_dust2.gif)

## Current features
You can build this right now, and, assuming you set the configuration to point to an existing Source game installation (this is tested primarily against CS:S):
* Loads game data files from projects gameinfo.txt
* Load BSP maps
* Load high-resolution texture data for bsp faces, including pakfile entries
* Full visibility data support
* Staticprop loading (working, but is incomplete)
* Basic entdata loading (dynamic and physics props)

## Installation
Windows, Mac & Linux are all supported.

There is a small amount of configuration required to get this project running, beyond `go get`.
* For best results, you need a source engine game installed already.
* Copy `config.example.json` to `config.json`, and update the `gameDirectory` property to point to whatever game installation
you are targeting (e.g. HL2 would be `<steam_dir>/steamapps/common/hl2`).

## Contributing
1. Fork it (<https://github.com/galaco/lambda-client/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request
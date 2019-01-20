# Ludus
--------
Playful love

ludus is a Love2d build and release command

ludus is aimed at making it easy to build love into a distributable package and run that package constantly so that your are aware of all your bugs

```
Usage:
  ludus [command]
Available Commands:
  build       build all versions of love
  clean       remove all build artifacts
  help        Help about any command
  run         run your game
  version     Print the version number of wrp
Flags:
  -a, --author string             author name
  -b, --build_directory string    directory where builds are outputted relative to your config file (default "build")
  -d, --description string        short description of the game
  -e, --email string              email of the author
  -h, --help                      help for ludus
  -p, --homepage string           homepage for the game
  -i, --identifier string         short description of the game (default "com.love.workspace")
  -l, --love_version string       version of love to build
  -s, --source_directory string   directory that contains your source relative to your config file (default ".")
  -t, --title string              title of the game (default "workspace")
  -v, --version string            version of your game
Use "ludus [command] --help" for more information about a command.
```

## Todo
- [ ] lua bytecode compile
- [ ] App Icons
    - [ ] osx plist change
- [ ] Better logging/output

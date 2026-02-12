# MC-Gomoku
Minecraft mod to Connect to Gomoku's Server

To build the mod execute `./gradlew build` (will remap the sources the 1st time)

Then in Minecraft `1.21.11` with fabric add the mod from `MC-Client/build/libs/` to the mod folder with `fabric-api`, `architectury` and `cloth-config`

`fabric-api` is for the connection between your code and Minecraft\
`architectury` and `cloth-config` are only used for the config (Based URL), so you don't have to rebuild to change the Gomoku server

### In Game:
`/go` create an vs AI\
`/go create local` to create an player vs player local\
`/go create remote` to create an player vs player remotly\
`/go join [code]` to join an player vs player\
`/go giveUp` to giveUp the game

Closing the client will also make you gave up


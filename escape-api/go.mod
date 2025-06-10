module escape-api

go 1.23.4

replace escape-engine => ../escape-engine

require (
	escape-engine v0.0.0-00010101000000-000000000000
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require github.com/gorilla/websocket v1.5.3 // indirect

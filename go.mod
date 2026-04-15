module steam-launch-custom/wrapper

go 1.23.0

require steam-launch-custom/wrapper/devtool v0.0.0 // indirect

tool steam-launch-custom/wrapper/devtool

replace steam-launch-custom/wrapper/devtool => ./devtool

require github.com/hajimehoshi/go-steamworks v0.0.0-20251207152439-f178e387e2a4

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

require (
	github.com/ebitengine/purego v0.9.1 // indirect
	github.com/spf13/cobra v1.10.2
)

package main

import "devtool/commands"

func main() {

	commands.RootCommand.Execute()

	// var cmd = os.Args[1]

	// switch cmd {
	// case "build":
	// 	build()
	// 	return
	// case "patch":
	// 	version(Patch)
	// 	return
	// case "minor":
	// 	version(Minor)
	// 	return
	// case "major":
	// 	version(Major)
	// 	return
	// default:
	// 	fmt.Printf("Unknwon command: %s", cmd)
	// }
}

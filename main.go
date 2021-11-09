package main

import (
	"github.com/yishixincheng/go-tools/cli"
	"log"
	"os"
)

func main()  {
	if len(os.Args) < 2 {
		log.Fatalln("\n参数缺失，示例如下\ngo run main.go sql2code --table=cdk")
	}
	cmdName := os.Args[1]
	if err := cli.RegisterCmd(cmdName); err != nil {
		log.Fatalf("%v", err)
	}
	if err := cli.Run(); err != nil {
		log.Fatalf("execute run fail: %v", err)
	}
}

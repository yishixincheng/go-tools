package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	_ "github.com/yishixincheng/go-tools/config"
)

var rootCmd = &cobra.Command{}
var registerMap = make(map[string]*cobra.Command)

// RegisterCmd 注册命令
func RegisterCmd(cmdName string) error {
	if cmd, ok := registerMap[cmdName]; ok {
		rootCmd.AddCommand(cmd)
		return nil
	}
	return fmt.Errorf("命令%s无效", cmdName)
}

// Run 执行命令
func Run() error {
	return rootCmd.Execute()
}

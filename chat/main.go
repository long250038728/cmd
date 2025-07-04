package main

import (
	"fmt"
	"github.com/long250038728/cmd/chat/client"
	"github.com/spf13/cobra"
)

func main() {
	chatCron := client.NewChatCorn()

	rootCmd := &cobra.Command{
		Use:   "llmchat",
		Short: "快速llm工具",
	}
	rootCmd.AddCommand(chatCron.Chat())
	rootCmd.AddCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

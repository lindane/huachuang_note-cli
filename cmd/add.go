package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add \"笔记内容\"",
	Short: "添加一条笔记",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content := strings.Join(args, " ")
		id, err := store.AddNote(content)
		if err != nil {
			return fmt.Errorf("添加笔记失败: %w", err)
		}
		fmt.Printf("笔记已添加，ID: %d\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search \"关键词\"",
	Short: "搜索包含关键词的笔记",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := strings.Join(args, " ")
		notes, err := store.SearchNotes(keyword)
		if err != nil {
			return fmt.Errorf("搜索笔记失败: %w", err)
		}
		if len(notes) == 0 {
			fmt.Printf("未找到包含'%s'的笔记\n", keyword)
			return nil
		}
		fmt.Printf("找到 %d 条包含'%s'的笔记：\n\n", len(notes), keyword)
		for _, note := range notes {
			fmt.Printf("ID: %d\n", note.ID)
			fmt.Printf("创建时间: %s\n", note.CreatedAt.Format(time.RFC3339))
			fmt.Printf("内容:\n%s\n", note.Content)
			fmt.Println("----------------------------------------")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

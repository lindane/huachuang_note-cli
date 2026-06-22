package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有笔记",
	RunE: func(cmd *cobra.Command, args []string) error {
		notes, err := store.ListNotes()
		if err != nil {
			return fmt.Errorf("获取笔记列表失败: %w", err)
		}
		if len(notes) == 0 {
			fmt.Println("暂无笔记")
			return nil
		}
		fmt.Printf("%-6s  %-55s  %s\n", "ID", "内容摘要", "创建时间")
		fmt.Println("------  -------------------------------------------------------  -------------------")
		for _, note := range notes {
			summary := truncate(note.Content, 50)
			created := note.CreatedAt.Format(time.RFC3339)
			fmt.Printf("%-6d  %-55s  %s\n", note.ID, summary, created)
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "删除指定 ID 的笔记",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id uint64
		_, err := fmt.Sscanf(args[0], "%d", &id)
		if err != nil {
			return fmt.Errorf("无效的 ID: %s", args[0])
		}
		err = store.DeleteNote(id)
		if err != nil {
			return fmt.Errorf("删除笔记失败: %w", err)
		}
		fmt.Printf("笔记 %d 已删除\n", id)
		return nil
	},
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
}

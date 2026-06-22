package cmd

import (
	"errors"
	"fmt"
	"time"

	"huachuang-note/storage"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [ID]",
	Short: "查看指定 ID 的笔记详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id uint64
		_, err := fmt.Sscanf(args[0], "%d", &id)
		if err != nil {
			return fmt.Errorf("无效的 ID: %s", args[0])
		}
		note, err := store.GetNote(id)
		if err != nil {
			if errors.Is(err, storage.ErrNoteNotFound) {
				return fmt.Errorf("笔记 ID %d 不存在", id)
			}
			return fmt.Errorf("获取笔记失败: %w", err)
		}
		fmt.Printf("ID: %d\n", note.ID)
		fmt.Printf("创建时间: %s\n", note.CreatedAt.Format(time.RFC3339))
		if note.UpdatedAt != nil {
			fmt.Printf("修改时间: %s\n", note.UpdatedAt.Format(time.RFC3339))
		}
		fmt.Println("内容:")
		fmt.Println(note.Content)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

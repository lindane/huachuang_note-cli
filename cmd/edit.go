package cmd

import (
	"errors"
	"fmt"
	"strings"

	"huachuang-note/storage"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [ID] \"新内容\"",
	Short: "编辑指定 ID 的笔记内容",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id uint64
		_, err := fmt.Sscanf(args[0], "%d", &id)
		if err != nil {
			return fmt.Errorf("无效的 ID: %s", args[0])
		}
		newContent := strings.Join(args[1:], " ")
		err = store.EditNote(id, newContent)
		if err != nil {
			if errors.Is(err, storage.ErrNoteNotFound) {
				return fmt.Errorf("笔记 ID %d 不存在", id)
			}
			return fmt.Errorf("编辑笔记失败: %w", err)
		}
		fmt.Printf("笔记 %d 已更新\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

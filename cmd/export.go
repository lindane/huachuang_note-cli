package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出所有笔记为 JSON 文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		notes, err := store.ListNotes()
		if err != nil {
			return fmt.Errorf("获取笔记列表失败: %w", err)
		}
		filename := fmt.Sprintf("notes_export_%s.json", time.Now().Format("20060102"))
		data, err := json.MarshalIndent(notes, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化笔记失败: %w", err)
		}
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}
		fmt.Printf("已导出 %d 条笔记到 %s\n", len(notes), filename)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

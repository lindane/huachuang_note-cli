package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"huachuang-note/storage"

	"github.com/spf13/cobra"
)

var store *storage.Store

var rootCmd = &cobra.Command{
	Use:   "note",
	Short: "简单的命令行笔记管理工具",
	Long:  "一个基于 BoltDB 的轻量级命令行笔记管理工具，支持添加、列出和删除笔记。",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		dbPath := os.Getenv("NOTE_DB_PATH")
		if dbPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("获取用户主目录失败: %w", err)
			}
			dbPath = filepath.Join(home, ".note.db")
		}
		store, err = storage.Open(dbPath)
		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if store != nil {
			return store.Close()
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

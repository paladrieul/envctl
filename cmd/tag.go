package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envctl/internal/store"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags on environment variable keys",
	}

	var target string

	addCmd := &cobra.Command{
		Use:   "add <key> <tag>",
		Short: "Add a tag to a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Root().PersistentFlags().GetString("dir")
			s := store.New(dir)
			if err := s.AddTag(target, args[0], args[1]); err != nil {
				return fmt.Errorf("add tag: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tagged %s with %q in [%s]\n", args[0], args[1], target)
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <key> <tag>",
		Short: "Remove a tag from a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Root().PersistentFlags().GetString("dir")
			s := store.New(dir)
			if err := s.RemoveTag(target, args[0], args[1]); err != nil {
				return fmt.Errorf("remove tag: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed tag %q from %s in [%s]\n", args[1], args[0], target)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tags for a target",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Root().PersistentFlags().GetString("dir")
			s := store.New(dir)
			tags, err := s.LoadTags(target)
			if err != nil {
				return fmt.Errorf("load tags: %w", err)
			}
			if len(tags) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no tags found")
				return nil
			}
			for key, ts := range tags {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %v\n", key, ts)
			}
			return nil
		},
	}

	for _, sub := range []*cobra.Command{addCmd, removeCmd, listCmd} {
		sub.Flags().StringVarP(&target, "target", "t", "default", "deployment target")
		tagCmd.AddCommand(sub)
	}

	GetRootCmd().AddCommand(tagCmd)
}

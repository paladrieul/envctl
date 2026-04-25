package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/envctl/internal/store"
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show the audit log of environment variable mutations",
		RunE:  runAudit,
	}
	auditCmd.Flags().StringP("target", "t", "", "filter entries by target name")
	auditCmd.Flags().StringP("action", "a", "", "filter by action (set, delete, import, copy, rename)")

	NewRootCmd().AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	dir, err := store.DefaultDir()
	if err != nil {
		return fmt.Errorf("resolve store dir: %w", err)
	}

	al := store.NewAuditLog(dir)
	entries, err := al.ReadAll()
	if err != nil {
		return fmt.Errorf("read audit log: %w", err)
	}

	targetFilter, _ := cmd.Flags().GetString("target")
	actionFilter, _ := cmd.Flags().GetString("action")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tTARGET\tACTION\tKEY\tENCRYPTED")

	count := 0
	for _, e := range entries {
		if targetFilter != "" && e.Target != targetFilter {
			continue
		}
		if actionFilter != "" && string(e.Action) != actionFilter {
			continue
		}
		encStr := "no"
		if e.Encrypted {
			encStr = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Target,
			e.Action,
			e.Key,
			encStr,
		)
		count++
	}
	w.Flush()

	if count == 0 {
		fmt.Println("No audit entries found.")
	}
	return nil
}

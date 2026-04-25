package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envctl/internal/store"
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage environment variable snapshots",
	}

	saveCmd := &cobra.Command{
		Use:   "save <target>",
		Short: "Save a snapshot of a target's current variables",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshotSave,
	}
	saveCmd.Flags().StringP("label", "l", "", "Optional label for the snapshot")

	listCmd := &cobra.Command{
		Use:   "list <target>",
		Short: "List all snapshots for a target",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshotList,
	}

	restoreCmd := &cobra.Command{
		Use:   "restore <target> <snapshot-index>",
		Short: "Restore a target from a snapshot by index (1-based)",
		Args:  cobra.ExactArgs(2),
		RunE:  runSnapshotRestore,
	}

	snapshotCmd.AddCommand(saveCmd, listCmd, restoreCmd)
	NewRootCmd().AddCommand(snapshotCmd)
}

func runSnapshotSave(cmd *cobra.Command, args []string) error {
	target := args[0]
	label, _ := cmd.Flags().GetString("label")

	s, err := store.New("")
	if err != nil {
		return err
	}

	snap, err := s.SaveSnapshot(target, label)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved: %s [%s]\n", snap.CreatedAt.Format("2006-01-02 15:04:05"), snap.Label)
	return nil
}

func runSnapshotList(cmd *cobra.Command, args []string) error {
	target := args[0]

	s, err := store.New("")
	if err != nil {
		return err
	}

	snapshots, err := s.ListSnapshots(target)
	if err != nil {
		return err
	}
	if len(snapshots) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No snapshots found for target %q\n", target)
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tTIMESTAMP\tLABEL\tVARS")
	for i, snap := range snapshots {
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\n", i+1, snap.CreatedAt.Format("2006-01-02 15:04:05"), snap.Label, len(snap.Vars))
	}
	return w.Flush()
}

func runSnapshotRestore(cmd *cobra.Command, args []string) error {
	target := args[0]
	var idx int
	if _, err := fmt.Sscanf(args[1], "%d", &idx); err != nil || idx < 1 {
		return fmt.Errorf("invalid snapshot index %q: must be a positive integer", args[1])
	}

	s, err := store.New("")
	if err != nil {
		return err
	}

	snapshots, err := s.ListSnapshots(target)
	if err != nil {
		return err
	}
	if idx > len(snapshots) {
		return fmt.Errorf("snapshot index %d out of range (1-%d)", idx, len(snapshots))
	}

	snap := snapshots[idx-1]
	if err := s.RestoreSnapshot(snap); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Restored %q from snapshot %s\n", target, snap.CreatedAt.Format("2006-01-02 15:04:05"))
	_ = os.Stdout.Sync()
	return nil
}

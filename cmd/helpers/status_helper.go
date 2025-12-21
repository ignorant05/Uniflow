package helpers

import (
	"fmt"
	"time"
)

func FormatStatus(status string) string {
	switch status {
	case "queued":
		return "Queued"
	case "in_progress":
		return "In Progress"
	case "completed":
		return "Completed"
	case "waiting":
		return " Waiting"
	default:
		return status
	}
}

func FormatConclusion(conc string) string {
	switch conc {
	case "success":
		return "Success"
	case "failure":
		return "Failure"
	case "cancelled":
		return "Cancelled"
	case "skipped":
		return " Skipped"
	case "timed_out":
		return "Timed Out"
	case "action_required":
		return "Action Required"
	case "neutral":
		return "Neutral"
	case "":
		return "Pending"
	default:
		return conc
	}
}

func FormatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())

		if mins == 1 {
			return "One minute ago"
		}

		return fmt.Sprintf("%d minutes ago", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())

		if hours == 1 {
			return "One hour ago"
		}

		return fmt.Sprintf("%d minutes ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)

		if days == 1 {
			return "One day ago"
		}

		return fmt.Sprintf("%d days ago", days)
	} else {
		return t.Format("Feb 6, 2004")
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

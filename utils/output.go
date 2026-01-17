package utils

import "regexp"

func SanitizeOutput(input string) string {
	// Matches /sandbox/job_.../ and replaces with ./
	re := regexp.MustCompile(`/sandbox/[^/ \n]+/`)
	return re.ReplaceAllString(input, "./")
}

func TruncateOutput(input string, limit int) string {
	if len(input) > limit {
		return input[:limit] + "\n... [Output Truncated]"
	}
	return input
}

package server

import "strings"

func CountTokens(text string) int {
    return len(strings.Fields(text))
}

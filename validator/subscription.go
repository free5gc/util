package validator

import "unicode"

// Check empty strings, control characters (including NUL),
// Check characters that would affect URL parsing such as '/', '?', '#'.
func IsValidSubscriptionID(id string) bool {
    if len(id) == 0 {
        return false
    }

    for _, r := range id {
        if r <= 31 || r == 127 { // control characters
            return false
        }
        switch r {
        case '/', '?', '#':
            return false
        }
        if unicode.IsSpace(r) {
            return false
        }
    }

    return true
}

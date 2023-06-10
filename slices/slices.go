package slices

func ContainsString(slice []string, s string) bool {
    for _, val := range slice {
        if val == s {
            return true
        }
    }
    return false
}

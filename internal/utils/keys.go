package utils

import "sort"

// SortedKeys returns a sorted list of keys from the map.
func SortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// MergedSortedKeys returns a sorted list of unique keys from both maps.
func MergedSortedKeys(m1, m2 map[string]interface{}) []string {
	keysMap := make(map[string]struct{}, len(m1)+len(m2))
	for k := range m1 {
		keysMap[k] = struct{}{}
	}
	for k := range m2 {
		keysMap[k] = struct{}{}
	}

	keys := make([]string, 0, len(keysMap))
	for k := range keysMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

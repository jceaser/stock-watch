//***********************************************************************************************100
/* Common functions in this project */

package common

import (
	//"sort"
)

/*func KeysFrom[K comparable, V any] (m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}*/

/*
Returns a sorted list of keys from a map.
Generic type is the map type:
	names := common.SortStringKeys[Command](Cmd.List)
*/
/*func SortStringKeys[V any] (m map[string]V) []string {
	keys := KeysFrom[string, V](m)
	sort.Strings(keys)
	return keys
}*/

func Contains[T comparable](list []T, item T) bool {
    for _, value := range list {
        if value == item {
            return true
        }
    }
    return false
}
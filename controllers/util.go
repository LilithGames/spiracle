package controllers

import "reflect"

func contains(items []string, x string) bool {
	for i := range items {
		if items[i] == x {
			return true
		}
	}
	return false
}

func keys(dict interface{}) []string {
	kvs := reflect.ValueOf(dict).MapKeys()
	r := make([]string, len(kvs))
	for i := range kvs {
		r[i] = kvs[i].String()
	}
	return r
}

func union(a []string, b []string) []string {
	t := make(map[string]struct{}, len(a)+len(b))
	for i := range a {
		t[a[i]] = struct{}{}
	}
	for i := range b {
		t[b[i]] = struct{}{}
	}
	return keys(t)
}

package controllers

import (
	"reflect"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func min(a []time.Duration) *time.Duration {
	if len(a) == 0 {
		return nil
	}
	m := a[0]
	for i := range a {
		if a[i] < m {
			m = a[i]
		}
	}
	return &m
}

func ptr(t metav1.Time) *metav1.Time {
	return &t
}

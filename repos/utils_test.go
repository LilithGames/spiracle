package repos

import (
    "testing"
    "time"
)

func TestConcurrentMap_SetGet(t *testing.T) {
    cm := NewConcurrentMap[int]()

    // Test set
    cm.Set("one", 1, time.Second*2)
    cm.Set("two", 2, 0)

    // Test get
    value, ok := cm.Get("one")
    if !ok {
        t.Error("Expected to find key 'one'")
    }
    if value != 1 {
        t.Errorf("Expected key 'one' to have value '1', but got '%v'", value)
    }

    value, ok = cm.Get("two")
    if !ok {
        t.Error("Expected to find key 'two'")
    }
    if value != 2 {
        t.Errorf("Expected key 'two' to have value '2', but got '%v'", value)
    }

    // Test expired item
    time.Sleep(time.Second * 3)
    _, ok = cm.Get("one")
    if ok {
        t.Error("Expected key 'one' to be expired and not found")
    }
}

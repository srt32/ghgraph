package cache

import (
	"testing"
	"time"
)

func TestCache_GetSet(t *testing.T) {
	c := NewCache(5 * time.Minute)
	token := "test-token"
	key := "test-key"
	value := "test-value"

	// Test Set and Get
	c.Set(token, key, value)
	got := c.Get(token, key)

	if got == nil {
		t.Fatal("Expected value, got nil")
	}

	if got.(string) != value {
		t.Errorf("Expected %q, got %q", value, got)
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	c := NewCache(5 * time.Minute)
	token := "test-token"
	key := "non-existent"

	got := c.Get(token, key)
	if got != nil {
		t.Errorf("Expected nil for non-existent key, got %v", got)
	}
}

func TestCache_TokenScoping(t *testing.T) {
	c := NewCache(5 * time.Minute)
	key := "shared-key"
	token1 := "token1"
	token2 := "token2"
	value1 := "value1"
	value2 := "value2"

	// Set values for different tokens with same key
	c.Set(token1, key, value1)
	c.Set(token2, key, value2)

	// Verify each token gets its own value
	got1 := c.Get(token1, key)
	got2 := c.Get(token2, key)

	if got1 == nil || got1.(string) != value1 {
		t.Errorf("Token1: expected %q, got %v", value1, got1)
	}

	if got2 == nil || got2.(string) != value2 {
		t.Errorf("Token2: expected %q, got %v", value2, got2)
	}
}

func TestCache_Expiration(t *testing.T) {
	ttl := 100 * time.Millisecond
	c := NewCache(ttl)
	token := "test-token"
	key := "test-key"
	value := "test-value"

	// Set value
	c.Set(token, key, value)

	// Verify it exists
	got := c.Get(token, key)
	if got == nil {
		t.Fatal("Expected value immediately after Set")
	}

	// Wait for expiration
	time.Sleep(ttl + 50*time.Millisecond)

	// Verify it's expired
	got = c.Get(token, key)
	if got != nil {
		t.Errorf("Expected nil after expiration, got %v", got)
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache(5 * time.Minute)
	token := "test-token"
	key := "test-key"
	value := "test-value"

	// Set and verify
	c.Set(token, key, value)
	if c.Get(token, key) == nil {
		t.Fatal("Expected value after Set")
	}

	// Delete and verify
	c.Delete(token, key)
	got := c.Get(token, key)
	if got != nil {
		t.Errorf("Expected nil after Delete, got %v", got)
	}
}

func TestCache_Clear(t *testing.T) {
	c := NewCache(5 * time.Minute)
	token := "test-token"
	otherToken := "other-token"

	// Set multiple values for the test token
	c.Set(token, "key1", "value1")
	c.Set(token, "key2", "value2")
	c.Set(token, "key3", "value3")

	// Set a value for a different token
	c.Set(otherToken, "key1", "other-value")

	// Clear all entries for test token
	c.Clear(token)

	// Verify test token entries are cleared
	if c.Get(token, "key1") != nil {
		t.Error("Expected key1 to be cleared")
	}
	if c.Get(token, "key2") != nil {
		t.Error("Expected key2 to be cleared")
	}
	if c.Get(token, "key3") != nil {
		t.Error("Expected key3 to be cleared")
	}

	// Verify other token's entry still exists
	got := c.Get(otherToken, "key1")
	if got == nil || got.(string) != "other-value" {
		t.Errorf("Expected other token's value to remain, got %v", got)
	}
}

func TestCache_MakeKey(t *testing.T) {
	c := NewCache(5 * time.Minute)
	token1 := "token1"
	token2 := "token2"
	key := "test-key"

	// Generate keys for different tokens
	fullKey1 := c.makeKey(token1, key)
	fullKey2 := c.makeKey(token2, key)

	// Verify keys are different (token-scoped)
	if fullKey1 == fullKey2 {
		t.Error("Expected different keys for different tokens")
	}

	// Verify keys have the expected format (hash:key)
	if len(fullKey1) < len(key)+17 { // 16 hex chars + ":" + key
		t.Error("Key format unexpected")
	}

	// Verify same token+key produces same cache key
	fullKey1Again := c.makeKey(token1, key)
	if fullKey1 != fullKey1Again {
		t.Error("Expected same cache key for same token+key")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache(5 * time.Minute)
	token := "test-token"
	iterations := 100

	// Concurrent writes
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < iterations; j++ {
				key := "key"
				value := id*iterations + j
				c.Set(token, key, value)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify cache didn't crash and can still work
	c.Set(token, "final-key", "final-value")
	got := c.Get(token, "final-key")
	if got == nil || got.(string) != "final-value" {
		t.Error("Cache state corrupted after concurrent access")
	}
}

func TestCache_CleanupGoroutine(t *testing.T) {
	ttl := 200 * time.Millisecond
	c := NewCache(ttl)
	token := "test-token"

	// Set multiple values
	for i := 0; i < 5; i++ {
		c.Set(token, "key"+string(rune(i)), "value")
	}

	// Wait for cleanup to run (TTL + buffer)
	time.Sleep(ttl + 150*time.Millisecond)

	// Set a new value to verify cache still works
	c.Set(token, "new-key", "new-value")

	// Verify new value exists
	got := c.Get(token, "new-key")
	if got == nil || got.(string) != "new-value" {
		t.Error("Cache should still work after cleanup")
	}

	// Old values should be cleaned up (may or may not be gone yet depending on timing)
	// This is non-deterministic, so we just verify the cache is still functional
}

func TestCache_DifferentValueTypes(t *testing.T) {
	c := NewCache(5 * time.Minute)
	token := "test-token"

	// Test string
	c.Set(token, "string-key", "string-value")
	if got := c.Get(token, "string-key"); got.(string) != "string-value" {
		t.Error("String value mismatch")
	}

	// Test int
	c.Set(token, "int-key", 42)
	if got := c.Get(token, "int-key"); got.(int) != 42 {
		t.Error("Int value mismatch")
	}

	// Test struct
	type testStruct struct {
		Name  string
		Count int
	}
	testVal := testStruct{Name: "test", Count: 10}
	c.Set(token, "struct-key", testVal)
	got := c.Get(token, "struct-key")
	if gotStruct, ok := got.(testStruct); !ok || gotStruct.Name != "test" || gotStruct.Count != 10 {
		t.Error("Struct value mismatch")
	}

	// Test pointer
	ptrVal := &testStruct{Name: "ptr", Count: 20}
	c.Set(token, "ptr-key", ptrVal)
	gotPtr := c.Get(token, "ptr-key").(*testStruct)
	if gotPtr.Name != "ptr" || gotPtr.Count != 20 {
		t.Error("Pointer value mismatch")
	}
}

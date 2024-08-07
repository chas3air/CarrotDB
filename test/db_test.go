package test

import (
	"testing"
	"time"

	"github.com/chas3air/CarrotDB/db"
)

func TestInit(t *testing.T) {
	db := db.Init()
	if len(db.Values) != 0 {
		t.Errorf("Expected DB.values map to have capacity 10, got %d", len(db.Values))
	}
	if len(db.TempKey) != 0 {
		t.Errorf("Expected DB.tempKey map to have capacity 10, got %d", len(db.TempKey))
	}
	if db.RecoverTime != 5 {
		t.Errorf("Expected DB.recoverTime to be 5, got %d", db.RecoverTime)
	}
}

func TestSetLifetime(t *testing.T) {
	db := db.Init()
	db.SetLifetime(10)
	if db.RecoverTime != 10 {
		t.Errorf("Expected DB.recoverTime to be 10, got %d", db.RecoverTime)
	}
}

func TestGetItem(t *testing.T) {
	db := db.Init()
	db.SetItem("key1", "value1")

	// Test case: Existing key
	item, err := db.GetItem("key1")
	if err != nil {
		t.Errorf("GetItem returned an error for existing key: %v", err)
	}
	if item != "value1" {
		t.Errorf("GetItem returned unexpected value: got %v, want %v", item, "value1")
	}

	// Test case: Non-existing key
	item, err = db.GetItem("non-existing-key")
	if err == nil {
		t.Errorf("GetItem did not return an error for non-existing key")
	}
	if item != nil {
		t.Errorf("GetItem returned unexpected value: got %v, want %v", item, nil)
	}
}

func TestSetItem(t *testing.T) {
	db := db.Init()
	db.SetItem("key1", "value1")
	item, err := db.GetItem("key1")
	if err != nil {
		t.Errorf("GetItem returned an error: %v", err)
	}
	if item != "value1" {
		t.Errorf("SetItem did not set the value correctly. Got %v, want %v", item, "value1")
	}
}

func TestSetTempItem(t *testing.T) {
	db := db.Init()
	db.SetTempItem("key1", "value1", 5)
	item, err := db.GetItem("key1")
	if err != nil {
		t.Errorf("GetItem returned an error: %v", err)
	}
	if item != "value1" {
		t.Errorf("SetTempItem did not set the value correctly. Got %v, want %v", item, "value1")
	}

	timedItem, ok := db.TempKey["key1"]
	if !ok {
		t.Errorf("SetTempItem did not set the timed item in DB.tempKey")
	}
	if time.Now().Sub(timedItem.Borntime) > 5*time.Second {
		t.Errorf("SetTempItem did not set the correct expire time. Got %v, want %v", time.Now().Sub(timedItem.Borntime), 5*time.Second)
	}
}

func TestDeleteItem(t *testing.T) {
	db := db.Init()
	db.SetItem("key1", "value1")
	db.DeleteItem("key1")
	_, err := db.GetItem("key1")
	if err == nil {
		t.Errorf("DeleteItem did not delete the item correctly")
	}
}

func TestClear(t *testing.T) {
	db := db.Init()
	db.SetItem("key1", "value1")
	db.SetItem("key2", "value2")
	db.Clear()
	if len(db.GetItems()) != 0 {
		t.Errorf("Clear did not clear the DB.values map correctly")
	}
}

func TestGetItems(t *testing.T) {
	db := db.Init()
	db.SetItem("key1", "value1")
	db.SetItem("key2", "value2")
	items := db.GetItems()
	if len(items) != 2 {
		t.Errorf("GetItems did not return the correct number of items. Got %d, want %d", len(items), 2)
	}
	if items["key1"] != "value1" {
		t.Errorf("GetItems did not return the correct value for key1. Got %v, want %v", items["key1"], "value1")
	}
	if items["key2"] != "value2" {
		t.Errorf("GetItems did not return the correct value for key2. Got %v, want %v", items["key2"], "value2")
	}
}

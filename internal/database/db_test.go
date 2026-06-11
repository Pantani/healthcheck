package database

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestDatabaseSaveAndGetData(t *testing.T) {
	ctx := context.Background()
	server := miniredis.RunT(t)

	db, err := Init(ctx, "redis://"+server.Addr()+"/0")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer db.Close()

	if err := db.SaveData(ctx, "api", "status", 42); err != nil {
		t.Fatalf("SaveData() error = %v", err)
	}

	var got interface{}
	if err := db.GetData(ctx, "api", "status", &got); err != nil {
		t.Fatalf("GetData() error = %v", err)
	}
	if got != float64(42) {
		t.Fatalf("GetData() = %v (%T), want 42", got, got)
	}
}

func TestDatabaseGetMissingData(t *testing.T) {
	ctx := context.Background()
	server := miniredis.RunT(t)

	db, err := Init(ctx, "redis://"+server.Addr()+"/0")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer db.Close()

	var got interface{}
	if err := db.GetData(ctx, "api", "missing", &got); err == nil {
		t.Fatal("GetData() error = nil, want missing key error")
	}
}

func TestRedisOptionsRejectsEmptyURL(t *testing.T) {
	if _, err := redisOptions(""); err == nil {
		t.Fatal("redisOptions() error = nil, want error")
	}
}

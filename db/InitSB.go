package db

import (
	"fmt"
	"github.com/supabase-community/supabase-go"
	"os"
)

func InitSupabase() (*supabase.Client, error) {
	anonKey := os.Getenv("ANON_KEY")
	url := os.Getenv("SUPABASE_URL")

	if anonKey == "" || url == "" {
		return nil, fmt.Errorf("wrong env config")
	}

	return supabase.NewClient(url, anonKey, nil)
}

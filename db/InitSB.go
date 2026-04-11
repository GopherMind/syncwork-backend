package db

import (
	"github.com/supabase-community/supabase-go"
	"os"
)

var SB *supabase.Client
func InitSupabase() {
	anonKey := os.Getenv("ANON_KEY")
	url := os.Getenv("SUPABASE_URL")

	client, err := supabase.NewClient(url, anonKey, nil)
    if err != nil {
        panic("Ошибка подключения к Supabase")
	
    }
	SB = client
}



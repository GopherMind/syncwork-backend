package db

import (
	"github.com/supabase-community/supabase-go"
	"os"
)

var SB *supabase.Client
func InitSupabase() {
	secretKey := os.Getenv("SECRET_KEY")
	url := os.Getenv("SUPABASE_URL")

	client, err := supabase.NewClient(url, secretKey, nil)
    if err != nil {
        panic("Ошибка подключения к Supabase")

    }
	SB = client
}



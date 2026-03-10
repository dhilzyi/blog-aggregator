package main

import (
	"blog-aggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, c command) error {
		userData, err := s.db.GetUser(context.Background(), s.Cfg.Username)
		if err == sql.ErrNoRows {
			fmt.Println("user is not exist within database")
			os.Exit(1)
		} else if err != nil {
			log.Fatal(err)
		}

		return handler(s, c, userData)
	}
}

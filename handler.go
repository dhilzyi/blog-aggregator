package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func handlerLogins(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("not enough arguments were provided")
	} else if len(cmd.arguments) < 1 {
		return fmt.Errorf("username argument is required")
	}

	ctx := context.Background()
	queryData, err := s.db.GetUser(ctx, cmd.arguments[0])
	if err == sql.ErrNoRows {
		fmt.Println("user is not exist within database")
		os.Exit(1)
	} else if err != nil {
		return err
	}

	s.Cfg.Username = queryData.Name
	if err := config.Write(s.Cfg); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("User has been set to", cmd.arguments[0])

	return nil
}

func handleRegister(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("not enough arguments")
	}

	ctx := context.Background()
	queryResult, err := s.db.GetUser(ctx, cmd.arguments[0])
	if err != sql.ErrNoRows && err != nil {
		return err
	}

	if queryResult.Name == cmd.arguments[0] {
		fmt.Println("user already exist")
		os.Exit(1)
	}

	arg := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}

	dataUser, err := s.db.CreateUser(ctx, arg)
	if err != nil {
		return err
	}

	fmt.Println("user is created successfuly")

	s.Cfg.Username = dataUser.Name
	if err := config.Write(s.Cfg); err != nil {
		return err
	}

	fmt.Printf("\nid: %v\ncreatedAt:%v\nname:%s\n\n", dataUser.ID, dataUser.CreatedAt, dataUser.Name)

	return nil
}

func handlerReset(s *state, _ command) error {
	ctx := context.Background()
	if err := s.db.ResetUsers(ctx); err != nil {
		return err
	}

	fmt.Println("Reset all rows in 'users' table successfully")
	return nil
}

func handlerUsers(s *state, _ command) error {

	ctx := context.Background()
	queryData, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for i := range queryData {
		prefix := ""
		if queryData[i].Name == s.Cfg.Username {
			prefix = "(current)"
		}
		fmt.Printf("* %s %s\n", queryData[i].Name, prefix)
	}
	return nil
}

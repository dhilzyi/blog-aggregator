package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dhilzyi/blog-aggregator/internal/config"
	"github.com/dhilzyi/blog-aggregator/internal/database"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handlerLogins(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
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
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("not enough arguments. it needs 'name'")
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

	fmt.Printf("\nid: %v\ncreatedAt:%v\nname:%s\n", dataUser.ID, dataUser.CreatedAt, dataUser.Name)

	return nil
}

func handlerReset(s *state, _ command) error {
	ctx := context.Background()
	if err := s.db.ResetFeedFollows(ctx); err != nil {
		return err
	}
	if err := s.db.ResetFeeds(ctx); err != nil {
		return err
	}
	if err := s.db.ResetUsers(ctx); err != nil {
		return err
	}

	fmt.Println("Reset all rows for 'users' and 'feeds' table successfully")
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

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("arguments not enough")
	}

	timeIntervalReq, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every", timeIntervalReq)
	ticker := time.NewTicker(timeIntervalReq)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("arguments not enough. it needs 'name' and 'url'")
	}

	ctx := context.Background()

	input := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}

	queryData, err := s.db.AddFeed(ctx, input)
	if err != nil {
		return err
	}

	newCmd := command{
		arguments: []string{cmd.arguments[1]},
	}

	if err := handlerFollow(s, newCmd, user); err != nil {
		return err
	}

	fmt.Printf("Adding feed to database successfully\n")
	fmt.Fprintf(os.Stdout, "\nID: %v\nName: %s\nUrl:%v\nUserID: %v\n", queryData.ID, queryData.Name, queryData.Url, queryData.UserID)

	return nil
}

func handlerFeeds(s *state, _ command) error {
	ctx := context.Background()
	queryData, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for i := range queryData {
		inst := queryData[i]
		userData, err := s.db.GetUserWithID(ctx, inst.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("\nEntry %d\n", i+1)
		fmt.Printf("Name: %s\n", inst.Name)
		fmt.Printf("Url: %v\n", inst.Url)
		fmt.Printf("User Name: %s\n", userData.Name)
	}
	fmt.Println("")
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("url arguments is required")
	}

	ctx := context.Background()
	feedData, err := s.db.GetFeedFromURL(ctx, cmd.arguments[0])
	if err != nil {
		return err
	}

	input := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedData.ID,
	}

	queryData, err := s.db.CreateFeedFollow(ctx, input)
	if err != nil {
		return err
	}

	fmt.Printf("\nFeed Name: %s\n", queryData.FeedName)
	fmt.Printf("Username: %s\n ", queryData.UserName)

	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	ctx := context.Background()
	followData, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for i := range followData {
		inst := followData[i]
		fmt.Printf("- %s\n", inst.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("arguments provided is not enough")
	}
	ctx := context.Background()

	input := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.arguments[0],
	}
	if err := s.db.DeleteFeedFollow(ctx, input); err != nil {
		return err
	}

	fmt.Printf("Delete follow feed successfully\n")
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.arguments) != 0 {
		limitInt, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			fmt.Println(err)
		}
		limit = limitInt
	}
	ctx := context.Background()

	input := database.GetPostsUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	data, err := s.db.GetPostsUser(ctx, input)
	if err != nil {
		return err
	}

	for i := range data {
		inst := data[i]
		fmt.Println(" -", inst.Title)
		fmt.Println("Published Date: ", inst.PublishedAt)
		fmt.Println("Source RSS: ", inst.FeedSourceName)
		fmt.Println("")
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feedData, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	if err := s.db.MarkFeedFetched(ctx, feedData.ID); err != nil {
		return err
	}
	fmt.Printf("\nMark fetched success for feed id: %v\n", feedData.ID)

	newFeed, err := fetchFeed(ctx, feedData.Url)
	if err != nil {
		return err
	}

	if len(newFeed.Channel.Item) < 1 {
		return fmt.Errorf("feed returning 0 items")
	}
	fmt.Printf("Printing feed to the console from '%s'\n", feedData.Name)
	fmt.Printf("%d items is received.\n\n", len(newFeed.Channel.Item))

	count := 0
	for i := range newFeed.Channel.Item {
		if count >= 30 {
			fmt.Printf("\nPrinting is limited to 30 items.\n")
			break
		}

		inst := newFeed.Channel.Item[i]
		if inst.Title == "" {
			continue
		}

		desc := sql.NullString{}
		if inst.Description != "" {
			desc.String = inst.Description
			desc.Valid = true
		}

		layout := "Mon, 02 Jan 2006 15:04:05 Z0700"
		convertedTime, err := time.Parse(layout, inst.PubDate)
		if err != nil {
			fmt.Println(err)
		}

		input := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       inst.Title,
			Url:         inst.Link,
			Description: desc,
			PublishedAt: convertedTime,
			FeedID:      feedData.ID,
		}

		_, err = s.db.CreatePost(ctx, input)
		if err != nil {
			fmt.Println(err)
		}

		count++
	}
	fmt.Printf("\nSucceed record %d items to database.\n", count)

	return nil
}

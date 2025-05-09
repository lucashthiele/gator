package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

const defaultLimitSize = 2

func printPosts(posts []database.GetPostsByUserRow) {
	fmt.Println("Posts:")
	fmt.Printf("\n----------------------------------\n")
	for _, post := range posts {
		fmt.Printf("Post: %s\n\n", post.Title)
		fmt.Printf("Published at: %s\n", post.PublishedAt.Time.Format("01/02/2006"))
		fmt.Printf("%s\n\n", post.Description.String)
		fmt.Printf("Access this post through this link: %s\n", post.Url)
		fmt.Printf("\n----------------------------------\n")
	}
}

func HandlerBrowse(s *model.State, cmd model.Command, user *database.User) error {
	expectedArguments := 1

	if len(cmd.Arguments) > expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou can only inform the limit",
			expectedArguments,
			len(cmd.Arguments))
	}

	limit := defaultLimitSize
	var err error
	if len(cmd.Arguments) == 1 {
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("wrong type for parameter limit")
		}
	}

	getPostsParams := database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := s.Db.GetPostsByUser(context.Background(), getPostsParams)
	if err != nil {
		return err
	}

	printPosts(posts)

	return nil
}

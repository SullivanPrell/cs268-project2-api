package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"cs268-project2-api/graph/generated"
	"cs268-project2-api/graph/model"
	"cs268-project2-api/mongo"
	"cs268-project2-api/user"
	"fmt"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.Login, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUser) (*model.User, error) {
	validatedUser, errors := user.ValidateInfo(input)
	if errors.Errors == true {
		fmt.Print(errors.Message)
		return &model.User{
			ID:            "",
			Email:         "",
			FirstName:     "",
			LastName:      "",
			DateOfBirth:   "",
			Major:         "",
			Minor:         "",
			WillingToHelp: false,
			Posts:         nil,
			Comments:      nil,
			PostIds:       nil,
			CommentIds:    nil,
			ClassesTaken:  nil,
			EmailVerified: nil,
			Token:         nil,
			Error:         &errors,
		}, nil
	}
	user, errors := mongo.CreateUser(validatedUser)
	if errors.Errors == true {
		fmt.Print(errors.Message)
		return &model.User{
			ID:            "",
			Email:         "",
			FirstName:     "",
			LastName:      "",
			DateOfBirth:   "",
			Major:         "",
			Minor:         "",
			WillingToHelp: false,
			Posts:         nil,
			Comments:      nil,
			PostIds:       nil,
			CommentIds:    nil,
			ClassesTaken:  nil,
			EmailVerified: nil,
			Token:         nil,
			Error:         &errors,
		}, nil
	}
	user.Error = &errors
	return &user, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifyEmail(ctx context.Context, input model.VerifyEmailInput) (*model.VerifyEmail, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, input model.UserInput) (*model.User, error) {
	returnUser := mongo.FindOneUser(input.Email)
	return &returnUser, nil
}

func (r *queryResolver) Posts(ctx context.Context, input model.PostsInput) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Post(ctx context.Context, input model.PostInput) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Comments(ctx context.Context, input model.CommentsInput) ([]*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Comment(ctx context.Context, input model.CommentInput) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

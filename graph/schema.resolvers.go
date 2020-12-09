package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"cs268-project2-api/graph/generated"
	"cs268-project2-api/graph/model"
	"cs268-project2-api/post"
	"cs268-project2-api/user"
	"cs268-project2-api/userAuth"
	"fmt"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.Login, error) {
	loginReturn := userAuth.Login(input)
	return &loginReturn, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUser) (*model.User, error) {
	returnUser := user.CreateNewUser(input)
	return &returnUser, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	returnPost := post.CreateNewPost(input)
	return &returnPost, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, input model.UserInput) (*model.User, error) {
	returnUser := user.FindUser(input)
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.


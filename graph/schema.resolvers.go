package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"cs268-project2-api/comment"
	"cs268-project2-api/graph/generated"
	"cs268-project2-api/graph/model"
	"cs268-project2-api/post"
	thread2 "cs268-project2-api/thread"
	"cs268-project2-api/user"
	"cs268-project2-api/userAuth"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.Login, error) {
	loginReturn := userAuth.Login(input)
	return &loginReturn, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUser) (*model.User, error) {
	returnUser := user.CreateNewUser(input)
	return &returnUser, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.PostSingle, error) {
	returnPost, errors := post.CreateNewPost(input)
	returnPost.Error = &errors
	return &returnPost, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateComment) (*model.PostSingle, error) {
	returnPost, errors := comment.CreateNewComment(input)
	returnPost.Error = &errors
	return &returnPost, nil
}

func (r *queryResolver) User(ctx context.Context, input model.UserInput) (*model.User, error) {
	returnUser := user.FindUser(input)
	return &returnUser, nil
}

func (r *queryResolver) Posts(ctx context.Context, input model.PostsInput) (*model.Posts, error) {
	posts, errors := post.GetPosts(input)
	returnVal := model.Posts{
		Posts: posts,
		Error: &errors,
	}
	return &returnVal, nil
}

func (r *queryResolver) Post(ctx context.Context, input model.PostInput) (*model.PostSingle, error) {
	returnPost, errors := post.GetPost(input)
	returnPost.Error = &errors
	return &returnPost, nil
}

func (r *queryResolver) PostsByUser(ctx context.Context, input model.PostsByUserIDInput) (*model.Posts, error) {
	posts, errors := post.GetPostsByUserID(input)
	returnVal := model.Posts{
		Posts: posts,
		Error: &errors,
	}
	return &returnVal, nil
}

func (r *queryResolver) PostByThread(ctx context.Context, input model.PostsByThreadInput) (*model.Posts, error) {
	posts, errors := post.GetPostsByThread(input)
	returnVal := model.Posts{
		Posts: posts,
		Error: &errors,
	}
	return &returnVal, nil
}

func (r *queryResolver) Comments(ctx context.Context, input model.CommentsInput) (*model.Comments, error) {
	comments, errors := comment.GetComments(input)
	returnVal := model.Comments{
		Comments: comments,
		Error:    &errors,
	}
	return &returnVal, nil
}

func (r *queryResolver) CommentsByPostID(ctx context.Context, input model.CommentsByPostIDInput) (*model.Comments, error) {
	comments, errors := comment.GetCommentsByPostID(input)
	returnVal := model.Comments{
		Comments: comments,
		Error:    &errors,
	}
	return &returnVal, nil
}

func (r *queryResolver) Thread(ctx context.Context, input model.ThreadInput) (*model.ThreadSingle, error) {
	thread, errors := thread2.GetThread(input)
	thread.Errors = &errors
	return &thread, nil
}

func (r *queryResolver) Threads(ctx context.Context, input model.ThreadsInput) (*model.Threads, error) {
	threads, errors := thread2.GetThreads(input)
	returnVal := model.Threads{
		Threads: threads,
		Errors:  &errors,
	}
	return &returnVal, nil
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

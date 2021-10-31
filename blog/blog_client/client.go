package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jisuhan3201/go-grpc/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect : %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Jisu",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", res)
	blogID := res.GetBlog().GetId()

	// Read blog
	fmt.Println("Reading the blog")

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "invalid id"})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v\n", readBlogErr)
	}
	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// Update blog
	fmt.Println("Updating the blog")

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Jisu",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v\n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	// Delete blog
	fmt.Println("Deleting the blog")
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v\n", err)
	}
	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// List blog
	fmt.Println("List blogs")

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unexpected error : %v", err)
		}

		fmt.Println(res.GetBlog())
	}
}

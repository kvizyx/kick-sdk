package kicksdk

import (
	"context"
	"fmt"
	"github.com/glichtv/kick-sdk/internal/urloptional"
	"net/http"
)

type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type Categories struct {
	client *Client
}

func (c *Client) Categories() Categories {
	return Categories{client: c}
}

type SearchCategoriesInput struct {
	Query string
}

// Search searches for Categories based on the search input.
//
// Reference: https://docs.kick.com/apis/categories#categories
func (c Categories) Search(ctx context.Context, input SearchCategoriesInput) (Response[[]Category], error) {
	resource := Resource{
		Type: ResourceTypeAPI,
		Path: "public/v1/categories",
	}

	request := NewRequest[[]Category](
		ctx,
		c.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodGet,
			AuthType: AuthTypeUserToken,
			URLValues: urloptional.Values{
				"q": urloptional.Single(input.Query),
			},
		},
	)

	return request.Execute()
}

type GetCategoryByIDInput struct {
	CategoryID int
}

// GetByID retrieves Category based on it's ID.
//
// Reference: https://docs.kick.com/apis/categories#categories-category_id
func (c Categories) GetByID(ctx context.Context, input GetCategoryByIDInput) (Response[Category], error) {
	resource := Resource{
		Type: ResourceTypeAPI,
		Path: fmt.Sprintf("%s/%d", "public/v1/categories", input.CategoryID),
	}

	request := NewRequest[Category](
		ctx,
		c.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodGet,
			AuthType: AuthTypeUserToken,
		},
	)

	return request.Execute()
}

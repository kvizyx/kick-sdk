package kicksdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/glichtv/kick-sdk/internal/urloptional"
)

type Category struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

type CategoriesResource struct {
	client *Client
}

func (c *Client) Categories() CategoriesResource {
	return CategoriesResource{client: c}
}

type SearchCategoriesInput struct {
	Query string
}

// Search searches for CategoriesResource based on the search input.
//
// Reference: https://docs.kick.com/apis/categories#categories
func (c CategoriesResource) Search(ctx context.Context, input SearchCategoriesInput) (Response[[]Category], error) {
	resource := c.client.NewResource(ResourceTypeAPI, "public/v1/categories")

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
func (c CategoriesResource) GetByID(ctx context.Context, input GetCategoryByIDInput) (Response[Category], error) {
	resource := c.client.NewResource(
		ResourceTypeAPI,
		fmt.Sprintf("%s/%d", "public/v1/categories", input.CategoryID),
	)

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

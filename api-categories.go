package kickkit

import (
	"context"
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
	return Categories{
		client: c,
	}
}

type (
	SearchCategoriesInput struct {
		Query string
	}

	SearchCategoriesOutput struct {
		Categories       []Category
		ResponseMetadata ResponseMetadata
	}
)

func (c Categories) Search(ctx context.Context, input SearchCategoriesInput) (SearchCategoriesOutput, error) {
	const resource = "categories"
	return SearchCategoriesOutput{}, nil
}

type (
	GetCategoryByIDInput struct {
		CategoryID int
	}

	GetCategoryByIDOutput struct {
		Category         Category
		ResponseMetadata ResponseMetadata
	}
)

func (c Categories) ByID(ctx context.Context, input GetCategoryByIDInput) (GetCategoryByIDOutput, error) {
	const resource = "categories"
	return GetCategoryByIDOutput{}, nil
}

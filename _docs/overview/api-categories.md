# Categories

Documentation for the categories API. Official documentation is [here](https://docs.kick.com/apis/categories). 

## Payloads
Various categories API payloads (entities). 
```go
type Category struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}
```

## Search Categories

Searching for categories that matching ```some-game-name``` query string.

```go
package main

import (
	"context"
	"fmt"
	"log"

	kicksdk "github.com/glichtv/kick-sdk"
)

func main() {
	client := kicksdk.NewClient(
		kicksdk.WithAccessTokens(kicksdk.AccessTokens{
			UserAccessToken: "user-access-token",
		}),
	)

	categories, err := client.Categories().Search(
		context.Background(),
		kicksdk.SearchCategoriesInput{
			Query: "some-game-name",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%v\n", categories)
}
```

## Get Category

Getting category with the ID ```1```.

```go
package main

import (
	"context"
	"fmt"
	"log"

	kicksdk "github.com/glichtv/kick-sdk"
)

func main() {
	client := kicksdk.NewClient(
		kicksdk.WithAccessTokens(kicksdk.AccessTokens{
			UserAccessToken: "user-access-token",
		}),
	)

	category, err := client.Categories().GetByID(
		context.Background(),
		kicksdk.GetCategoryByIDInput{
			CategoryID: 1,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%v\n", category)
}
```

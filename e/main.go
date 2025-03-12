package main

import (
	"context"
	"fmt"
	kickkit "github.com/glichtv/kick-kit"
)

func main() {
	client := kickkit.NewClient(
		kickkit.WithCredentials(kickkit.Credentials{
			ClientID:     "01JP0AB647B21FPMQAWAN915AE",
			ClientSecret: "52b66d5d4d09379d5d20748875114b567d08cc1ab2d42470c083a27b2b9c5c1d",
			RedirectURI:  "http://localhost:8080/api/auth/callback",
		}),
		kickkit.WithAccessTokens(kickkit.AccessTokens{
			UserAccessToken: "MZNIZJU4OTUTZMNHNY0ZOTC3LWIXOWITZGY0MDU5MJKXZGRL", // NZVHOWFJNGMTYTGWOS01MJCXLWI2YJUTMTK2YWQ3ZTG3ZJHI
		}),
	)

	//res, err := client.Chat().PostMessage(
	//	context.Background(),
	//	kickkit.PostChatMessageInput{
	//		BroadcasterUserId: 0,
	//		Content:           "aboba!",
	//		PosterType:        kickkit.MessagePosterUser,
	//	},
	//)
	//if err != nil {
	//	panic(err)
	//}

	res, err := client.Users().InspectToken(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Data)
	fmt.Printf("%d %s\n", res.ResponseMetadata.StatusCode, res.ResponseMetadata.KickMessage)
}

package kicksdk

import "time"

type EventType = string

const (
	EventTypeChatMessage             EventType = "chat.message.sent"
	EventTypeChannelFollow           EventType = "channel.followed"
	EventTypeChannelSubRenewal       EventType = "channel.subscription.renewal"
	EventTypeChannelSubGifts         EventType = "channel.subscription.gifts"
	EventTypeChannelSubCreated       EventType = "channel.subscription.new"
	EventTypeLivestreamStatusUpdated EventType = "livestream.status.updated"
)

type (
	Emote struct {
		ID        string          `json:"emote_id"`
		Positions []EmotePosition `json:"positions"`
	}

	EmotePosition struct {
		Start int `json:"s"`
		End   int `json:"e"`
	}

	Broadcaster struct {
		IsAnonymous    bool   `json:"is_anonymous"`
		UserID         int    `json:"user_id"`
		Username       string `json:"username"`
		IsVerified     bool   `json:"is_verified"`
		ProfilePicture string `json:"profile_picture"`
		ChannelSlug    string `json:"channel_slug"`
	}
)

type (
	EventChatMessage struct {
		MessageID   string      `json:"message_id"`
		Broadcaster Broadcaster `json:"broadcaster"`
		Sender      Broadcaster `json:"sender"`
		Content     string      `json:"content"`
		Emotes      []Emote     `json:"emotes"`
	}

	EventChannelFollow struct {
		Broadcaster Broadcaster `json:"broadcaster"`
		Follower    Broadcaster `json:"follower"`
	}

	EventChannelSubscriptionRenewal struct {
		Broadcaster Broadcaster `json:"broadcaster"`
		Subscriber  Broadcaster `json:"subscriber"`
		Duration    int         `json:"duration"`
		CreatedAt   time.Time   `json:"created_at"`
	}

	EventChannelSubscriptionGifts struct {
		Broadcaster Broadcaster   `json:"broadcaster"`
		Gifter      Broadcaster   `json:"gifter"`
		Giftees     []Broadcaster `json:"giftees"`
		CreatedAt   time.Time     `json:"created_at"`
	}

	EventChannelSubscriptionCreated struct {
		Broadcaster Broadcaster `json:"broadcaster"`
		Subscriber  Broadcaster `json:"subscriber"`
		Duration    int         `json:"duration"`
		CreatedAt   time.Time   `json:"created_at"`
	}

	EventLivestreamStatusUpdated struct {
		Broadcaster Broadcaster `json:"broadcaster"`
		IsLive      bool        `json:"is_live"`
		Title       string      `json:"title"`
		StartedAt   time.Time   `json:"started_at"`
		EndedAt     time.Time   `json:"ended_at"`
	}
)

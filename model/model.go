package model

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
)

type MentionType int

const (
	MENTION_TYPE_ROLES MentionType = iota
	MENTION_TYPE_USERS
	MENTION_TYPE_EVERYONE
)

func (m MentionType) String() string {
	return [...]string{"roles", "users", "everyone"}[m]
}

type BotConfig struct {
	PublicKey                  string
	ApplicationCommandHandlers map[string]func(ctx context.Context, applicationCommand ApplicationCommand) (InteractionResponse, error)
}

type InteractionResponse struct {
	Type int                     `json:"type"`
	Data InteractionResponseData `json:"data"`
}

type InteractionResponseData struct {
	TTS             bool                                    `json:"tts"`
	Content         string                                  `json:"content,omitempty"`
	Embeds          []any                                   `json:"embeds,omitempty,default=[]"`
	AllowedMentions *InteractionResponseDataAllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           int                                     `json:"flags,omitempty"`
	Components      any                                     `json:"components,omitempty"`
	Attachments     *InteractionResponseDataAttachment      `json:"attachments,omitempty"`
}

type InteractionResponseDataAllowedMentions struct {
	Parse       []MentionType `json:"parse,omitempty"`
	Roles       []string      `json:"roles,omitempty"`
	Users       []string      `json:"users,omitempty"`
	RepliedUser bool          `json:"replied_user,omitempty"`
}

type InteractionResponseDataAttachment struct {
	ID           string `json:"id"`
	Filename     string `json:"filename"`
	Description  string `json:"description"`
	ContentType  string `json:"content_type"`
	Size         int    `json:"size"`
	URL          string `json:"url"`
	ProxyURL     string `json:"proxy_url"`
	Height       int    `json:"height,omitempty"`
	Width        int    `json:"width,omitempty"`
	Ephemeral    bool   `json:"ephemeral,omitempty"`
	DurationSecs int    `json:"duration_secs,omitempty"`
	Waveform     any    `json:"waveform,omitempty"`
}

type Interaction struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	Type          int    `json:"type"`
	Data          string `json:"data,omitempty"`
	GuildID       string `json:"guild_id"`
	Channel       any    `json:"channel,omitempty"` //todo: map in unmarshal
	ChannelID     string `json:"channel_id"`
	Member        any    `json:"member,omitempty"` //todo: map in unmarshal
	User          User   `json:"user"`             //todo: map in unmarshal
	Token         string `json:"token"`
	Version       int    `json:"version"`
	Message       any    `json:"message,omitempty"` //todo: map in unmarshal
	AppPermission int    `json:"application_permission,omitempty"`
	Locale        string `json:"locale,omitempty"`
	GuildLocale   string `json:"guild_locale,omitempty"`
}

func (i *Interaction) UnmarshalJSON(data []byte) error { //todo: map rest of keys
	v := map[string]any{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v["id"] != nil {
		i.ID = v["id"].(string)
	}
	if v["application_id"] != nil {
		i.ApplicationID = v["application_id"].(string)
	}
	if v["type"] != nil {
		i.Type = int(v["type"].(float64))
	}
	if v["data"] != nil {
		applicationCommand := &ApplicationCommand{}
		bytes, err := json.Marshal(v["data"].(map[string]interface{}))
		if err != nil {
			return err
		}

		err = json.Unmarshal(bytes, applicationCommand)
		if err != nil {
			return err
		}

		b, err := json.Marshal(applicationCommand)

		i.Data = string(b)
	}
	if v["guild_id"] != nil {
		i.GuildID = v["guild_id"].(string)
	}
	if v["channel_id"] != nil {
		i.ChannelID = v["channel_id"].(string)
	}
	if v["token"] != nil {
		i.Token = v["token"].(string)
	}
	if v["application_permission"] != nil {
		i.AppPermission = int(v["application_permission"].(float64))
	}
	if v["locale"] != nil {
		i.Locale = v["locale"].(string)
	}
	if v["guild_locale"] != nil {
		i.GuildLocale = v["guild_locale"].(string)
	}

	return nil
}

type ApplicationCommand struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     int    `json:"type"`
	Resolved any    `json:"resolved,omitempty"`
	Options  any    `json:"options,omitempty"`
	GuildID  string `json:"guild_id,omitempty"`
	TargetID string `json:"target_id,omitempty"`
}

func (m ApplicationCommand) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", m.ID)
	e.Str("name", m.Name)
	e.Int("type", m.Type)
	e.Interface("resolved", m.Resolved)
	e.Interface("options", m.Options)
	e.Str("guild_id", m.GuildID)
	e.Str("target_id", m.TargetID)
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot,omitempty"`
	System        bool   `json:"system,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
	Banner        string `json:"banner,omitempty"`
	AccentColor   int    `json:"accent_color,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Verified      bool   `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         int    `json:"flags,omitempty"`
	PremiumType   int    `json:"premium_type,omitempty"`
	PublicFlags   int    `json:"public_flags,omitempty"`
}

/*
 Copyright (c) 2023 Ian Hulsbus

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Constructor
type DiscordClient interface{}

type Messages struct {
	discordClient *discordgo.Session
}

func MessagesConstructor(client *discordgo.Session) *Messages {
	return &Messages{
		discordClient: client,
	}
}

func (m Messages) UserJoinedChannelMessage(guildID string, channelID string, user discordgo.User) {
	message := fmt.Sprintf("▶️ User %s joined the channel!", user.Mention())
	m.discordClient.ChannelMessageSend(channelID, message)
}

func (m Messages) UserLeftChannelMessage(guildID string, channelID string, user discordgo.User) {
	message := fmt.Sprintf("🚮 User %s left the channel!", user.Mention())
	m.discordClient.ChannelMessageSend(channelID, message)
}

func (m Messages) JoinableChannelEmbed(guildID string, messageChannel string, channel *discordgo.Channel) error {

	embed := discordgo.MessageEmbed{
		Title:       fmt.Sprintf(`Joinable channel "%s"`, channel.Name),
		Type:        discordgo.EmbedTypeRich,
		Description: channel.Topic,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Channel link",
				Value:  channel.Mention(),
				Inline: false,
			},
			{
				Name:   "Use the reactions below to join or leave the channel",
				Inline: false,
			},
		},
	}

	message, err := m.discordClient.ChannelMessageSendEmbed(messageChannel, &embed)
	if err != nil {
		return err
	}

	err = m.discordClient.MessageReactionAdd(message.ChannelID, message.ID, "▶️")
	if err != nil {
		return err
	}

	err = m.discordClient.MessageReactionAdd(message.ChannelID, message.ID, "🚮")
	if err != nil {
		return err
	}

	return nil
}

func (m Messages) GetMessagesInChannel(channelID string) ([]*discordgo.Message, error) {
	var channelMessages []*discordgo.Message
	var beforeID string

	for {
		messages, err := m.discordClient.ChannelMessages(channelID, 100, beforeID, "", "")
		if err != nil {
			return nil, err
		}

		channelMessages = append(channelMessages, messages...)

		if len(messages) == 100 {
			beforeID = messages[99].ID
		} else {
			break
		}
	}

	return channelMessages, nil
}

func (m Messages) DeleteMessage(channelID, messageID string) error {
	err := m.discordClient.ChannelMessageDelete(channelID, messageID)
	if err != nil {
		return err
	}

	return nil
}

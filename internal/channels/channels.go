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

package channels

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Constructor
type DiscordClient interface{}

type Channels struct {
	discordClient *discordgo.Session
}

func ChannelsConstructor(client *discordgo.Session) *Channels {
	return &Channels{
		discordClient: client,
	}
}

func (c Channels) CreateTextChannel(guildID string, channelData discordgo.GuildChannelCreateData) (*discordgo.Channel, error) {
	len := len(channelData.Name)
	if len < 2 || len > 100 {
		return nil, fmt.Errorf("channel name length must be between 2 and 100 characters. current length: %d", len)
	}

	channel, err := c.discordClient.GuildChannelCreateComplex(guildID, channelData)
	if err != nil {
		errMsg := fmt.Sprintf("Channel creation failed. Error was: %s", err)

		if channel.ID != "" {
			_, err := c.discordClient.ChannelDelete(channel.ID)
			if err != nil {
				errMsg = fmt.Sprintf("%s. Additionally, the following error occured when reverting changes: %s", errMsg, err)
			}
		}

		return nil, errors.New(errMsg)
	}

	return channel, nil
}

func (c Channels) DeleteTextChannel(channelID string) error {
	_, err := c.discordClient.ChannelDelete(channelID)
	if err != nil {
		return err
	}

	return nil
}

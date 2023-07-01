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

package roles

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Constructor
type DiscordClient interface{}

type Roles struct {
	discordClient *discordgo.Session
}

func RolesConstructor(client *discordgo.Session) *Roles {
	return &Roles{
		discordClient: client,
	}
}

func (r Roles) RetrieveRoles(guildID string) ([]*discordgo.Role, error) {
	roles, err := r.discordClient.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r Roles) CreateRole(guildID string, data *discordgo.RoleParams) (*discordgo.Role, error) {
	len := len(data.Name)
	if len < 2 || len > 100 {
		return nil, fmt.Errorf("role name length must be between 2 and 100 characters. current length: %d", len)
	}

	role, err := r.discordClient.GuildRoleCreate(guildID, data)
	if err != nil {
		errMsg := fmt.Sprintf("Role creation failed. Error was: %s", err)

		if role.ID != "" {
			err := r.discordClient.GuildRoleDelete(guildID, role.ID)
			if err != nil {
				errMsg = fmt.Sprintf("%s. Additionally, the following error occured when reverting changes: %s", errMsg, err)
			}
		}

		return nil, errors.New(errMsg)
	}

	return role, nil
}

func (r Roles) DeleteRole(guildID, roleID string) error {
	err := r.discordClient.GuildRoleDelete(guildID, roleID)
	if err != nil {
		return err
	}

	return nil
}

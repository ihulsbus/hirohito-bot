# hirohito-bot
---

This is the main bot for the DAGC Discord server. 

This bot fulfills the following purposes:

* Allow users to join/leave optional channels
* Allow admins to automatically archive inactive channels
* Allow admins to create joinable channels

## Building and packaging
TBD

## Running the bot

The bot uses env variables for it's configuration. The table blow shows the supported environment variables.

|varname                |type     |required |
|---                    |---      |---      |
|HIROHITO_LOGLEVEL      |`string` |no       |
|HIROHITO_PREFIX        |`string` |no       |
|HIROHITO_DISCORD_TOKEN |`string` |yes      |

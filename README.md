Sheave-Bot
==========

Simple implementation of a irc bot which responds with upcoming meetups for a User Group.

Also contains functionality to respond to previous messages in the channel with all the words replaced with their anagrams when possible.

##Setup
Depends on having a Go already being installed.
> go get github.com/Ropes/Sheave-Bot

Create config file for authenticating with the IRC server; "sheave.conf" in main working directory. eg Freenode setup:
> {
    "Server": "chat.freenode.net:6667",
    "UserName": "botname",
    "RealName": "FullBotName",
    "Passwd": "password",
    "Channels": ["#pdxgo", "#pdxbots", "#channelstojoin"]
}

Build executable
> go build

###Running

Run executable
> ./sheave


# d3tour's Nomi Room Manager

[Nomi](https://nomi.ai) is a platform that offers AI companions for human users to chat with. They have opened v1 of their API which enables Nomi chatting that occurs outside of the Nomi app or website in a feature called "Rooms" which function like a group chat. Some integrations, like my [NomiKin-Discord](https://github.com/d3tourrr/NomiKin-Discord) project allow users to chat with their Nomis on other platforms, and some of those integrations create Nomi Rooms on behalf of you, the user.

Although Rooms function like a Group Chat, they are separate. In-app Group Chats are *not* Rooms, and vice versa. Nomi has not provided a way for users to see which Rooms they have, which Nomis are in which rooms, create, delete or modify rooms, outside of directly accessing their API. There is no Room management in the Nomi app. Since integrations like [NomiKin-Discord](https://github.com/d3tourrr/NomiKin-Discord) abstract the API calls away from a user, a user could be left without any way of knowing about or being able to manage their Rooms.

This small commandline application offers a way for Nomi users to manage their Nomi Rooms. You can: 

* List Nomis
* List Rooms
* Create new Rooms
* Delete existing Rooms
* Add a Nomi to a Room
* Remove a Nomi from a Room
* Update a Room Name, Note or Backchanneling setting

The Nomi Room Manager is not intended as a way for you to interact with your Nomi. It is provided to streamline the process of managing your Rooms.

# Using Nomi Room Manager

Go to the [Releases](https://github.com/d3tourrr/NomiRoomManager/releases/latest) and download the correct app for your system. Nomi Room Manager runs on Windows, macOS and Linux. After downloading the app, run it like you would run any other.

## Windows

1. Open a Command Prompt or PowerShell terminal
1. Navigate to the location the app was downloaded
1. Type `.\NomiRoomManager-Windows.exe`

## macOS

1. Open a Terminal
1. Navigate to the location where the app was downloaded
1. Type `./NomiRoomManager-macOS`

## Linux

1. Open a Terminal
1. Navigate to the location where the app was downloaded
1. Type `./NomiRoomManager-Linux`

# Earthwalker

![Create new game dialog](readme/image_create_new_game.png)
![Ingame](readme/image_ingame.png)
![Summary](readme/image_summary.png)

Earthwalker is a game of a similar concept to [GeoGuessr](geoguessr.com).
You get dropped somewhere in the world in Google Street View, and the goal is that you find out where you are,
and guess better than all of your friends. You can play against time, restrict the game to an area, and more.

It's free and open source, and the idea is that people host it themselves to play with their friends. No Google
API keys are needed, as Earthwalker "fools" the public Google Street View a bit. This is technically against Google TOS,
so I am not hosting a public version of this myself.

## How do I play it?

You need to host earthwalker or find a friend who hosts it.

## Okay, how do I host it?

### Disclaimer

First, I need to tell you that this program is technically against Google's Terms of Service, as it hides some UI elements on Street View,
and filters out information from some Google packets that are sent to Street View. Hosting this game is your own responsibility.

### Ok, I want to play the game now.

Are you set? Okay, first of all you need a server (which could also be your own computer), and some ports that are forwarded on it.
Next, you need to install [Git](https://git-scm.com/) and [Go](https://golang.org/).
Clone this repo and build the program:

    git clone https://gitlab.com/glatteis/earthwalker.git
    cd earthwalker
    go build

The executable should be called `earthwalker` or `earthwalker.exe`. Run `./earthwalker(.exe) -h` to see how to use it.

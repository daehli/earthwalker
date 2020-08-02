# Earthwalker

Earthwalker is a game of a similar concept to [GeoGuessr](https://geoguessr.com).
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

Are you set? Okay, first of all you need a server (which could also be your own computer), and some ports that are forwarded on it (but you can
come back to these later).
Just write us, for instance write an issue, if you need any help.

#### Setup directly on Windows

I would strongly recommend installing a Windows Subsystem for Linux, for instance [the Debian one](https://www.microsoft.com/en-us/p/debian/),
and then following the "Setup directly on Linux" walkthrough. This does not mean this doesn't run directly on Windows,
but it's going to save you a lot of pain, probably. If you are hardcore and do want it directly on Windows, good luck to you,
just follow the "Setup directly on Linux" steps and install stuff differently.
If you're setting up directly on Windows, you might also need [make](http://www.gnu.org/software/make/),
if you're somewhere else this should be preinstalled.

#### Setup directly on Linux (or in the Windows Subsystem for Linux)

Next, you need to install [Git](https://git-scm.com/), [Go](https://golang.org/) and [node](https://nodejs.org/en/download/).
This can be done through `apt` if you're on Debian:

    apt-get install git
    apt-get install golang-go
    curl -sL https://deb.nodesource.com/setup_14.x | bash -
    apt-get install -y nodejs

If you're on another distribution, the above installations steps might be different.

Now, clone this repo and build the program:

    git clone https://gitlab.com/glatteis/earthwalker.git
    cd earthwalker
    make

The executable should be called `earthwalker` or `earthwalker.exe`.
Just running it and going to `localhost:8080` in your browser should work for you.
If you are running a server, you probably configure a custom port to work with your nginx or apache config.
How to do this is described in the help: `./earthwalker(.exe) -h`.

##### Updating

You can update earthwalker by running `git pull` in its directory and then running `make` again.

#### Docker

To use the docker container you have to run the following commands (given you already have docker installed and configured).
    
    git clone https://gitlab.com/glatteis/earthwalker.git
    cd earthwalker
    docker build -t earthwalker:local .
    docker run -p 8080:8080 earthwalker:local

It might not be necessary to use the `-t earthwalker:local` param, but it makes it a little prettier.
The website should be hosted at `localhost:8080`. The port can be remapped via docker.

#### Configuration

Some configuration options can be read from environment variables, a `.toml` file, or command line arguments; these are summarized below.  In all cases, command line arguments override environment variables, which override `.toml` values.  All configuration options are strings.  Using absolute paths is recommended.
A sample configuration file, `config.toml.sample`, is provided.  Rename/copy it to `config.toml` to start configuring.

<details>
<summary>Table of configuration options.</summary>

| Command Line Flag | Environment Variable                              | `.toml` Key          | Default                                                  | Comments |
|-------------------|---------------------------------------------------|----------------------|----------------------------------------------------------|----------|
|                   | EARTHWALKER_CONFIG_PATH                           |                      | ./config.toml                                            | Location of the `.toml` configuration file |
| port              | EARTHWALKER_PORT                                  | Port                 | 8080                                                     |          |
|                   | EARTHWALKER_DB_PATH                               | DBPath               | ./badger/                                                 | Location of the database directory |
|                   | EARTHWALKER_STATIC_PATH                           | StaticPath           | location of executable (usually `earthwalker`)           | Absolute path to the directory containing `static` and `templates` |
|                   |                                                   | TileServerURL        | https://tiles.wmflabs.org/osm/{z}/{x}/{y}.png            | URL of a raster tile server.  This determines what you see on the map. |
|                   |                                                   | NoLabelTileServerURL | https://tiles.wmflabs.org/osm-no-labels/{z}/{x}/{y}.png  | As above, but this value is used when a map creator has turned labels off. |

</details>

## Images

![Create new game dialog](readme/image_create_new.png)
![Ingame](readme/image_ingame.png)
![Summary](readme/image_summary.png)

# Earthwalker

Earthwalker is a game of a similar concept to [GeoGuessr](https://geoguessr.com).
You get dropped somewhere in the world in Google StreetView, and the goal is to you find out where you are and guess your location more precisely than all of your friends. You can play against the clock, restrict the game to an area, and more.

It's free and open source, and the idea is that people host it themselves to play with their friends. No Google API keys are needed, as Earthwalker "fools" the public Google Street View a bit. This is technically against Google TOS, so I am not hosting a public version of this myself.

## How do I play?

You need to host Earthwalker yourself or find a friend who hosts it. Don't worry; hosting isn't too difficult.
The following guide will focus on the simplest case: hosting and playing on the same computer. If you want to play
on other computers too, see the "hosting outside of your network" paragraph (but first focus on getting earthwalker installed :P).
If you need any help, you can create an Issue on this project's GitLab page.

If hosting on your own computer isn't an option or you want to just quickly try it out,
you can also for instance use [PlayWithDocker](https://labs.play-with-docker.com/) (this is by Docker, we're not affiliated).
After signing in, click on "Add new instance", copy the commands in the "Docker" instructions below inside the console and then click on
"8080" next to "Open Port". It will take you to a link where earthwalker should be running, which you can also share with your friends.
This isn't a good permanent solution though, as you can only play for four hours and all your maps will be lost after that time.

#### Disclaimer

First, I need to tell you that this program is technically against Google's Terms of Service, as it hides some UI elements on Street View, and filters out information from some Google packets that are sent to Street View. Hosting this game is your own responsibility.

#### Hosting on Linux or the Windows Subsystem for Linux (WSL)

Start by installing [Git](https://git-scm.com/), [Go](https://golang.org/) and [node](https://nodejs.org/en/download/).
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

You should now be able to run the `earthwalker` executable to start the server, and then go to `localhost:8080` in your browser to start playing!

#### Hosting on Windows

I would recommend installing a Windows Subsystem for Linux, for instance [the Debian one](https://www.microsoft.com/en-us/p/debian/),
and then following the steps for "Hosting on Linux" above. 

You can also host the game directly on Windows, but it may be a bit more complicated.    
As on Linux, you'll need to install [Git](https://git-scm.com/), [Go](https://golang.org/) and [node](https://nodejs.org/en/download/).  If you want to use [make](http://www.gnu.org/software/make/), you'll need to install it as well.  
Once you have those components installed and working, clone the repo and cd (change directory) into it:

    git clone https://gitlab.com/glatteis/earthwalker.git
    cd earthwalker

Then, run `make` if you've installed it, or if not:  
Compile the server:  

    go build 

And compile the front end:  

    cd frontend-svelte
    npm install
    npm run build

You should now be able to run `earthwalker.exe` to start the server, and then go to `localhost:8080` in your browser to start playing!

#### Docker

To use the docker container you have to run the following commands (given you already have docker installed and configured).
    
    git clone https://gitlab.com/glatteis/earthwalker.git
    cd earthwalker
    docker build -t earthwalker:local .
    docker run -p 8080:8080 earthwalker:local

It might not be necessary to use the `-t earthwalker:local` param, but it makes it a little prettier.
The website should be hosted at `localhost:8080`. The port can be remapped via docker.

#### Hosting outside of your network

There are multiple possibilities with which you can host outside of your network.

1. You already have a server somewhere with a cool nginx installation that has all ports forwarded anyway - well, you know you don't need to read this paragraph.
2. You have full control over your router and you don't want to use any third party services - find out your IP Address,
and choose any port of your liking - like 8080 or 7777. Open it to TCP in the router, and then start earthwalker on that port on your PC.
Then, <your-ip>:<opened-port> should be open to others.
3. Configuring your router is either impossible or too clumsy and you're fine with using third party services - take a look at
[tunnelto.dev](https://tunnelto.dev/) or [ngrok](https://ngrok.com/). They run on your PC with earthwalker and people will be able to
access earthwalker through these services. Note that all traffic will flow through their servers, so don't spill your secrets.
Remember to disable the service after stopping earthwalker.

#### Configuration

We've provided a handful of configuration options, which are read from your environment variables, a `.toml` file, or command line arguments (these are all summarized below).  In all cases, command line arguments override environment variables, which override `.toml` values.  All configuration options are strings.  Using absolute paths is recommended.  
You can rename or copy the provided sample configuration file, `config.toml.sample`, to `config.toml` to get started.

<details>
<summary>Table of configuration options.</summary>

| Command Line Flag | Environment Variable                              | `.toml` Key          | Default                                                  | Comments |
|-------------------|---------------------------------------------------|----------------------|----------------------------------------------------------|----------|
|                   | EARTHWALKER_CONFIG_PATH                           |                      | ./config.toml                                            | Location of the `.toml` configuration file |
| port              | EARTHWALKER_PORT                                  | Port                 | 8080                                                     |          |
|                   | EARTHWALKER_DB_PATH                               | DBPath               | ./badger                                                 | Location of the database directory |
|                   | EARTHWALKER_STATIC_PATH                           | StaticPath           | location of executable (usually `earthwalker`)           | Absolute path to the directory containing `static` and `templates` |
|                   |                                                   | TileServerURL        | https://tiles.wmflabs.org/osm/{z}/{x}/{y}.png            | URL of a raster tile server.  This determines what you see on the map. |
|                   |                                                   | NoLabelTileServerURL | https://tiles.wmflabs.org/osm-no-labels/{z}/{x}/{y}.png  | As above, but this value is used when a map creator has turned labels off. |

</details>

#### Updating

You can update earthwalker by running `git pull` in its directory, and then running `make` or following the compilation instructions again.

## Contributing

Contributions are welcome!  Check out [our TODO list on Trello](https://trello.com/b/cGc4oTqf/earthwalker) and the Issues page for this GitLab repo.  The application is written mostly in Go (back end) and Svelte/JavaScript (front end).

Even if you're not a developer, please submit an Issue on GitLab if you find any bugs or would like to request a feature.

## Images

![Create new game dialog](readme/image_create_new.png)
![Ingame](readme/image_ingame.png)
![Summary](readme/image_summary.png)

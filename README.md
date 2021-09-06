# LIVE

Youtube has a bug in their playlist feature when adding a live channel video to it.

The channel live will stop at some point of time and the live video become a normal video. so you have to regularily remove the video from your playlist and go to the channel and add the new live video to your list again.

This program will list all live videos for a list of youtube channels in one page. and will keep itself up to date with the latest live videos when it expires.

# What is it?

- A Go program
- Doesn't depend on any package
- An HTTP server that listens on port 3000
- It uses Youtube API because the list of live videos are not included in the public RSS feed for the channel

# Adding more channels

- open the youtube channel you want to follow
- Show page source
- Search for "channel_id=" to find the RSS URL and copy the channel ID
- open `channels.go`
- Add the channel ID to the list

# Setting up

- You need a google application with youtube api v3 enabled
- copy `.env.sample` to `.env`
- set the application key in the `.env` file
- make sure the `.env` file is loaded in your shell when running the program
- run the program with `go run .`
- it will listen on port 3000 so open `localhost:3000`
- it will get the information in the background every 1 hour and update the list of channels

# Contributing

- Keep it simple. try as hard as possible not to add any dependencies
- Add more channels
- Add useful features like searching in the list by name

# Notable files

- main.go : includes the program main entry point
- bin/css : a bash script that downloads bulma library and change some variables and compile it with SASS
- channels.go : include the list of channels IDs

# Deployment

the deployment script in `bin/deploy` can be used to deploy the program to remote server given you have ssh access to it.

- You need docker and docker-compose installed.
- the script will clone, build, run the container to your remote server
- The command takes for format `bin/deploy master user@server-ip-address web` will deploy `master` branch to your `server-ip-address` using user `user` credentials and deploy the service named `web`.
- Service listen on port 4001 inside the machine, you need to forward that with nginx or haproxy from port 80 based on your rules.

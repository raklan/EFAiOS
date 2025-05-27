An online version of [Escape From the Aliens in Outer Space](https://www.eftaios.com/), a game of strategy and bluff published by Osprey Games. 

My friends and I love playing this game, but wished we had a way to more easily adjust the game to our liking, as well as adding our own content, such as new roles, maps, and items. With this in mind, I looked around for an online version, but didn't see any that A) were standalone and B) had ways to create my own content. Given my experience in similar projects, I decided to make an online version myself.

My goal with this is twofold: First, a lightweight version that's easy to set up and play, similar to Jackbox Games, Kahoot, or Candlelight (my college senior project), and Second, to support as much customization as possible. I would love to eventually have some sort of system in place for user-generated roles/items and such, but currently, only maps can be created by users, and anything else must be implemented by myself.

I'll probably get it deployed in some fashion (and post the URL in this README when I do), but if for whatever reason that doesn't happen or you just want to run it yourself, you'll have three options:
  - I'll post a standalone executable in the releases that you can just download and run (since the entire project is done in Go, which compiles into user-friendly executables quite nicely)
  - I'll also probably push a docker image of the project that you can set up your own container for
  - You can clone the repository, install Go on your machine, and simply run the command `go run ./escape-api` from the root directory.

Either way, the site will run on port 80 of whatever machine is running it, so the default web connection should get you up and running.

![efaiosHomepage](https://github.com/user-attachments/assets/05326570-d2e1-464b-8da7-ed50fa05bbfe)
An online version of [Escape From the Aliens in Outer Space](https://www.eftaios.com/), a game of strategy and bluff published by Osprey Games. 

My friends and I love playing this game, but wished we had a way to more easily adjust the game to our liking, as well as adding our own content, such as new roles, maps, and items. With this in mind, I looked around for an online version, but didn't see any that A) were standalone and B) had ways to create my own content. Given my experience in similar projects, I decided to make an online version myself.

My goal with this is twofold: First, a lightweight version that's easy to set up and play, similar to Jackbox Games, Kahoot, or Candlelight (my college senior project), and Second, to support as much customization as possible. I would love to eventually have some sort of system in place for user-generated roles/items and such, but currently, only maps can be created by users, and anything else must be implemented by myself.

I'll probably get it deployed in some fashion (and post the URL in this README when I do), but if for whatever reason that doesn't happen or you just want to run it yourself, you'll have three options:
  - (Recommended) I've posted a ZIP folder in the releases, which is the easiest way to get set up - instructions below  
  - You can clone the repository, install Go on your machine, and simply run the command `go run ./escape-api` from the root directory. Be aware this option does not come with any premade maps.
  - I'll also probably push a docker image of the project that you can set up your own container for. Be aware that this option will not come with any premade maps.

Either way, the site will run on port 80 of whatever machine is running it, so the default web connection should get you up and running.

## Features
- Built-in, easy-to-use map editor for creating any map you can think of!
  ![image](https://github.com/user-attachments/assets/1b56d025-6eaa-4fd8-9f37-4bc932454943)

- Every Card & Role found in the Ultimate edition of the Tabletop version, as well as host controls over how many of each can be in the game
  ![image](https://github.com/user-attachments/assets/9f6f6916-0e20-4e57-abee-5eb96c9afd30)

- The Release Version comes prebundled with a recreation of every map found in the Ultimate edition of the Tabletop version
  ![image](https://github.com/user-attachments/assets/eb2e31d8-adba-47cc-a62d-972e28ce2d8d)

- Lost connection? Accidentally closed the page? No worries, we have rejoin support!
  ![image](https://github.com/user-attachments/assets/44e996fb-cc3a-40f8-8156-bd80421a1a8e)
  
- Easy to set up and play. Simply visit the Maps page, click "Play" on the map you want, and share the Room Code with the other players. Those players can join your lobby by entering their Name and the Room Code on the home page.
- Built-in Compendium with Game Rules and explanation of every card, role, etc.

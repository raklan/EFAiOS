![efaiosHomepage](https://github.com/user-attachments/assets/05326570-d2e1-464b-8da7-ed50fa05bbfe)
An online version of [Escape From the Aliens in Outer Space](https://www.eftaios.com/), a game of strategy and bluff published by Osprey Games. 

My friends and I love playing this game, but wished we had a way to more easily adjust the game to our liking, as well as adding our own content, such as new roles, maps, and items. With this in mind, I looked around for an online version, but didn't see any that A) were standalone, not requiring ANY external dependencies to run such as Java, C#, etc. and B) had ways to create my own content. Given my experience in similar projects, I decided to make an online version myself.

My goal with this is twofold: First, a lightweight version that's easy to set up and play, similar to Jackbox Games, Kahoot, or Candlelight (my college senior project), and Second, to support as much customization as possible. I would love to eventually have some sort of system in place for user-generated roles/items and such, but currently, only maps can be created by users, and anything else must be implemented by myself.

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
  
  ![image](https://github.com/user-attachments/assets/db17d4fc-966e-42a9-b8f9-7e2d01ca32f5)

- Built-in Compendium with Game Rules and explanation of every card, role, etc.
  ![image](https://github.com/user-attachments/assets/012a973d-d1b8-4d38-8717-4b984e9bef1a)

- Easy Map Sharing (Instructions below)


## How to Run
I'll probably get it deployed in some fashion (and post the URL in this README when I do), but if for whatever reason that doesn't happen or you just want to run it yourself, you'll have three options:
  ### Recommended/Easiest
  Check out the [Releases Page](https://github.com/raklan/EFAiOS/releases), which is the easiest way to get set up. Every release will have a ZIP folder included called EFAiOS-vX.Y.Z (replace X.Y.Z with the version number) that you can download and extract. There will be a .exe file inside called efaios.exe that you can run.
  ### Other Options
  - You can clone the repository, install Go on your machine, and simply run the command `go run ./escape-api` from the root directory. Be aware this option does not come with any premade maps.
  - I'll also probably push a docker image of the project that you can set up your own container for. Be aware that this option will not come with any premade maps.

However you run the project, the site will run on port 80 of whatever machine is running it, so connecting from any web browser should get you up and running.

## Map Sharing
With the map editor, I wanted to make sharing maps between computers easy; The idea is that if my friends want to make maps, they can run the server themselves, make a map, and easily send a file to me that I can import into the "official" server so we can all play their map. With that in mind, here's how to share maps:
### Exporting a map to GIVE TO SOMEONE
 1. Whatever folder you run the server from will have a folder called `maps` that includes every map saved to your local version of the server. Inside you'll see a bunch of files named `map_{gibberish here}.json` (If you have the prebundled official maps, I've given those files special names so they'll look different)
 2. Find the file of the map you want to share. There's not a great way to do this part, but you have two options.
    - Option 1 is you can open up each `.json` file in notepad or some other text editor. At the beginning of the file you'll see is a spot that says `"id":"{gibberish here}", "name":"{Map Name Here}"`. Using this, you can check the map name for the map you want.
    - Option 2 is you can start the server, open up your browser, and go to the Maps page. Find the map you want to share and hit "Edit" - this opens up the Map Editor with that map. Next, look at your browser's address bar - at the end of the address will be a spot that says `?id={text here}`. Take note of the text at the end, and find the file where the part of the file name after `map_` matches the text you see in the address bar, and that's your map.
 3. Once you've found the file of the map you want to share, you can just send that file to whoever you're sharing it with!
### Importing a map someone GAVE YOU
 1. Whatever folder you run the server from will have a folder called `maps` that includes every map saved to your local version of the server. Inside you'll see a bunch of files named `map_{gibberish here}.json` (If you have the prebundled official maps, I've given those files special names so they'll look different)
 2. Whoever shared the map with you should have given you a `.json` file with a title something like `map_{text here}.json`.
 3. Simply copy the file they gave you into this `maps` folder and it should show up the next time you go to the Maps page in your browser, no server restart required!

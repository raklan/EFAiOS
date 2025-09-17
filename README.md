![title](https://github.com/user-attachments/assets/5c8ef89a-18b4-4391-ae7e-4b7b030fc77d)
An online version of [Escape From the Aliens in Outer Space](https://www.eftaios.com/), a game of strategy and bluff published by Osprey Games. 

My friends and I love playing this game, but wished we had a way to more easily adjust the game to our liking, as well as adding our own content, such as new roles, maps, and items. With this in mind, I looked around for an online version, but didn't see any that A) were standalone, not requiring ANY external dependencies to run such as Java, C#, etc. and B) had ways to create my own content. Given my experience in similar projects, I decided to make an online version myself.

My goal with this is twofold: First, a lightweight version that's easy to set up and play, similar to Jackbox Games, Kahoot, or Candlelight (my college senior project), and Second, to support as much customization as possible. I would love to eventually have some sort of system in place for user-generated roles/items and such, but currently, only maps can be created by users, and anything else must be implemented by myself.

## Features
- Every Card & Role found in the Ultimate edition of the Tabletop version, as well as host controls over how many of each can be in the game
  ![image](https://github.com/user-attachments/assets/2fd42acb-a49e-4b12-98d2-e746841e090a)

- The Release Version comes prebundled with an editable recreation of every map found in the Ultimate edition of the Tabletop version
  ![image](https://github.com/user-attachments/assets/070f76a2-abd9-4e62-80bb-70fcfdc99755)

- Screen Drawing built-in so you can keep notes right on your device, as well as an Event Log to ensure you don't miss important turns.
  ![image](https://github.com/user-attachments/assets/275a5c02-5bcb-4780-b849-b42eef0b57f9)

- Post-game recap - See a play-by-play of every player's turn after the game ends.
  ![image](https://github.com/user-attachments/assets/21ec01e7-15bc-489c-9eb3-d20ef9c7e01e)

- Built-in, easy-to-use map editor for creating any map you can think of!
  ![image](https://github.com/user-attachments/assets/ac85ee68-b161-4e11-9148-a2508bfdb1f0)

- Lost connection? Accidentally closed the page? No worries, we have rejoin support!
  ![image](https://github.com/user-attachments/assets/44e996fb-cc3a-40f8-8156-bd80421a1a8e)
  
- Easy to set up and play. To Host, click the "Host" button, select your map, and share the Room Code with your friends. Then, they can enter that Room Code on the home page to join your game!
  ![image](https://github.com/user-attachments/assets/fa97247a-81c9-4d38-bffe-053e95ff19c4)

- Built-in Compendium with Game Rules and explanation of every card, role, etc.
  ![image](https://github.com/user-attachments/assets/6abaa3c7-3e3b-45d0-b7a9-d6f0b798f244)

- [Easy Map Sharing](#map-sharing)

## How to Run
I'll probably get it deployed in some fashion (and post the URL in this README when I do), but if for whatever reason that doesn't happen or you just want to run it yourself, you'll have three options:
  ### Recommended/Easiest
  Check out the [Releases Page](https://github.com/raklan/EFAiOS/releases), which is the easiest way to get set up. Every release will have a ZIP folder included called EFAiOS-vX.Y.Z (replace X.Y.Z with the version number) that you can download and extract. There will be a .exe file inside called `efaios.exe` that you can run.
  ### Other Options
  - You can clone the repository, install Go on your machine, and simply run the command `go run ./escape-api` from the root directory. 
  - I'll also probably push a docker image of the project that you can set up your own container for. Be aware that this option will not come with any premade maps.

However you run the project, the site will run on port 80 of whatever machine is running it, so any device on the same internet network can connect to your machine's IP via any web browser to play!

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

## Roadmap
If you'd like to see what new features I have planned, known bugs, or to report your own bug you've found, check out the [Issues Page](https://github.com/raklan/EFAiOS/issues)! All work I do is tracked there, so if you post there I'm very likely to respond.

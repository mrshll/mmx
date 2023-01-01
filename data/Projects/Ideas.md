## Radio hub (working title "yourelisteningto.com")
The pandemic brought to the fore methods of communication over distance. {reading list, Reading} All the Light We Cannot See, setting up a microradio station (dormant) and sharing a lot of playlists with family and friends has kept the technology of "radio" simmering in the back of my mind as something I'd like to explore more. In setting up my microradio station, it became clear that while technologies like [Icecast](https://icecast.org/) are mature and powerful, they are by no means easy. My pulse > [mpd](https://www.musicpd.org/) > Icecast > website flow was a Rube Goldberg machine. I'd love for a small program that allowed anyone to easily set up an online station, with a page that served as hub of stations around the world.

That hub could be the best place to start - indexing stations that stream on Twitch and via their own streaming stacks. If I wanted to take it further and provide some mechanism to do the streaming itself, I could:
+ build a native application that allows a user to mix audio-out and an audio-in (microphone) and combine them into a stream. ([{https://github.com/moriyoshi/pulsego/, pulsego](http://code.google.com/p/portaudio-go/,portaudio-go}))
+ build a web application that uses the browser's "share tab with audio" functionality to stream audio

There may already be some protocols or open source building blocks to do this sort of streaming (I recall some sort of library being compatible with things like Twitch too). Both would expose the stream in the hub at some reserved address, i.e. *yourelisteningto.com/mrshll*

## Online Battleline (card game, working title "Formations")
From what I can tell, there does not exist an online version of [Battleline](https://www.gmtgames.com/p-939-battle-line-11th-printing.aspx). It would be fun to make one, possibly using [boardgame.io](https://boardgame.io/) as the backing engine.

## {mmx} internal
Aside from paper notebooks, I haven't found digital notetaking or brainstorming venues that have _stuck_. I sometimes want a version of mmx that is either a native or web-based client with editor capabilities to parse {mmxup} at runtime and create structure - sort of like [Obsidian](https://obsidian.md/).

Web stuff is easy, but it would be a learning experience to build and deploy a native app across platforms. Since mmx was originally written in Golang, something like [Fyne](https://developer.fyne.io/) could be interesting to play with.

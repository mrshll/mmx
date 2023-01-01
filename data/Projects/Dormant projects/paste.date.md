2021-01-23

In what is certainly my most mundane project to date, I've created an application that loads calendar events from Google Calendar via a {https://twitter.com/marsh/status/1353404000469659648, prayer to the OAuth gods}, and uses lots of date and timezone math to find availablities based on the specified meeting parameters. I haven't had the opportunity to solve an algorithmically interesting problem recently -- so much of my technical work lately is technical piping and static sites. So it was fun to sit down and design an approach to solve a backpack-type problem.

![Screenshot of paste.date](img/pastedate-1.png)

I was prompted to build this by the existing solutions' reliance on flashy formatted pastes into an email and the ever-presumptive "calendar link." I wanted a simple program that give me plaintext times that I could quickly copy and paste. I might add a link option at a later date, but I'd keep it simple. It currently lives on localhost, but I hope to wrap up signup and iron out the last few bugs this month.

I've paused this project. I've I revisit it, I'd like to make it primarily cli-based.

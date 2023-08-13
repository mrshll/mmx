2023-03-05

Meadow is growing, work is A Lot, time is short, and the weather is weird. I think that sums up the time elapsed since [[Recently, January 2023]].
![Meadow and Dax, living in some sort of improbable unity](img/PXL_20230215_153011317.jpg)

My days end and begin anew with Meadow's snores and smiles. Her glad recognition, in the valley of fatigue, is stronger than a cup of coffee. We've found a rhythm this past month and have started exploring more confidently; last week we went on our first hike together in the Middlesex Fells with our pals Sam, Aliza, and Jude. 

Coincidentally, I'm now on the board of the [Friends of the Fells](https://www.friendsofthefells.org/) which aims to protect and enhance the nature of the Fells.

![Wrapping up Meadow's first hike, with our friends Sam, Aliza, and Jude](img/IMG_20230212_191046.jpg)

Getting back to work after my too-brief parental leave has been intense. Between the company ramping up a number of start-of-year efforts, taking meetings while feeding Meadow, and spinning up a fundraising process for the company, the days leave me worn out. I love my work, the [team I've built over the years](https://upstream.tech/about), the impact we have; adding Meadow into the mix has made me realize that I have to, at least in this phase, let many of my creative projects slumber and be especially disciplined to make time for recharging exercise in husk-state. Tough to swallow, but necessary for now.

## Reading
Books finished recently:

- _[[You're Paid What You're Worth]]_ by Jake Rosenfeld
- _Neuromancer_ by William Gibson

Up next:

- _The House of Spirits_ by Isabel Allende
- _Make Shift: Dispatches from the Post-Pandemic Future_ by Various Authors, part of the Twelve Tomorrows series

Good stuff on the net:

- [ChatGPT is a Blurry JPEG of the Web](https://www.newyorker.com/tech/annals-of-technology/chatgpt-is-a-blurry-jpeg-of-the-web) by Ted Chiang
- [Reverse engineering the Facebook Messenger API](https://intuitiveexplanations.com/tech/messenger); I enjoyed this and shared it with our engineering team - the playful and curious approach, as well as some of the specific methods and tools used encapsulate what I think is core to a software engineer who is self-sufficient and good at debugging.
- [Visual design rules you can safely follow every time](https://anthonyhobday.com/sideprojects/saferules/)

## Playing
Over the past few years, I've noticed that my preference of game has shifted from immersive, story-driven non-violent-ish single player games, to a seemingly disjoint group of multiplayer games. First it was Sea of Thieves at the start of the pandemic, then it was Hell Let Loose, and most recently DayZ. 

Immersion has always been something I look for in a game, to lose myself in the game-world. When I started playing multiplayer games in high school, I gravitated towards games with immersion through emergent social interactions. My friend Henry, who I play often with, hit the nail on the head recently when he said that the sub-genre we found so intriguing could be called "multiplayer story engines."

![In Sea of Thieves, my pals and I became a pirate crew](img/sot.png)

In contrast to multiplayer games that require twitchy reactions and move at a lightning pace, story engines tend to unfold more slowly, and the game mechanic leans towards exploration and interaction, with a sparse or large game world. Proximity-based voice or text chat is also a common theme that leads to emergent and interesting interaction with other players.

All of this is to say, I've been playing DayZ recently, which has all of these qualities in addition to a complicated (brutal) and intriguing game loop. I found it via a YouTuber named [SourSweet](https://www.youtube.com/watch?v=jpZejlBxbXc)(link cw: violence, zombies), who edits together his adventures into episodes that are massively popular. A story engine at its finest.

## Elsewhere
Internet pal [Eli](https://eli.li/) shared [Mynah](https://git.sr.ht/~eli_oat/mynah), which helps you keep a digital garden/wiki from a directory of markdown files. I was smitten with the approach Eli took, which was to use a simple bash script that leveraged [pandoc](https://pandoc.org/) to convert from markdown to html. It sparked a whole host of ideas, including one I toyed with (but had to pause) to create an Instapaper-esque cli that takes a web URL and uses html > markdown > html as a questionable compression method to create a local saved copy of web articles.

I also came across the concept of ssh as an application protocol, a la [charm.sh](https://charm.sh/). That turned me on to [their Go libraries](https://charm.sh/libs/), to create rich and dynamic command line tools. I had some fun toying with their examples and making a CLI for the aforementioned web archiving tool.

Check out Oppen's [video devlog of Playdate audio applications](https://www.youtube.com/@oppenlabyorkshire1240) - very fast and polished progress that has been fun to follow.

It's been a long time since coding was a core part of my job, but I've found time to sneak it in the in-betweens as it helps me decompress for other, more tiring responsibilities. I recently built a set of spatial visualizations that show hydrological, meteorological, and earth observation data as a spatial time series (be proud, Tufte!). I enjoy creating visualizations like these by rendering plain ol' svg elements with support from a melange of interpolation and scale libraries. Build some Javascript tooltip library? Nah, let's use `<path><title>content</title></path>`!

![A screenshot of the web dashboard of company's service HydroForecast showing a new spatial view that shoes the distribution of model inputs across hydrological sub-basins.](img/hf-spatial.png)

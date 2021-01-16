---
name: Building mmx
date: 2020-01-15
host: mmx
bref: An exercise in self-built tooling and learning GoLang
---
I've historically defaulted to the technology du jour to build and rebuild
my personal and company webpages. It's worked fine - there have been
moments of fantastic efficiency, and others of abysmal reverse
engineering of an errant Gatsby plugin.

In the spirit of {self-sufficient informatics}, and heavily inspired
by others on the {https://webring.xxiivv.com/|webring}, I sought out to
simplify my dependencies, build something myself, and learn some new
technologies in the process.

To start, I mapped out my hoped-for characteristics of the end result.
I landed on: 

- low friction (so as to promote more writing);
- low-level (no automagic frameworks or dependencies; and
- extendable (I want this to be somethign that evolves with me over the next 5, 10, 15 years).

{https://github.com/XXIIVV/oscean/|Devine Lu Linvega's Oscean} served
as the primary inspiration. I spent a weekend pouring over their wiki,
the underlying C-code, and the ecosystem of file formats and prorgams
they created. I could (and <strike>should</strike> will) write a whole
post on the great things I continue to learn from them. The act of
reading through their code and reverse engineering the site compilation
was one of the most fun weekends I've had in a while. I promise I am fun.
I gained an understanding of their technical approach to linking and
render disparate databases, while simultaneously exploring the content of
those databases themselves.

I opted for GoLang as a learning opportunity. Compared to my usual
Javascript, GoLang is miles closer to C - the last time I thought about
pointers was my senior year of college, if that. I never imagined saying
this, but I missed pointers! My impression of GoLang so far is fine. I
hardly tapped into the features it is known for, such as concurrancy. But
this project served as a gentle introduction.

Another reason I wanted to create the compiler myself was so that I could add
features over time that are typically only available on "platforms," such as
bidirectional linking and other memex-style data graph functions.
I noticed that Oscean as well as others on the webring were able to do
this. How cool.

In a moment of doubt, I played with {https://www.11ty.dev/|11ty} as well as
{https://github.com/mathieudutour/gatsby-digital-garden|gatsby plugins}
that promised functionality I sought. But after speedbumps with
each, always grappling with the obfuscation that make them "magic," I felt
confident that I was on the right path.

For my pages, I opted to use a modified version of the human readable
{https://wiki.xxiivv.com/site/indental.html|Indental} format, and started
by porting a portion of my site to it, and writing a parser to convert it
to Go structs. From there, I linked the structs into a hierarchical tree
and created various rendering functions that converted the ndtl entries
into html, resting a bit on Go's native templating engine. It works
remarkably well, and I am able to easily add features and improvements
over time. For example, I recently added the ability for an entry's body
to contain markdown.

Other features of established web frameworks are replicatable with bash,
usually. For instance, "live reloading" was achieved with the following
bash:

```
#!/bin/sh
bash build.sh
while inotifywait -qqre modify ./src ./links ./data; do
  bash build.sh
done
```
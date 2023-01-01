2021-06-24
# This site is called mmx

![A person building something with bamboo, from The Red Turtle](img/redturtle020.jpg)

mmx was built at the start of 2021. Largely inspired by others on [webring](https://webring.xxiivv.com/), especially [Devine](https://wiki.xxiivv.com/site/home.html), the site aims to be a long term repository for my writing, notes, and research. Entries on my previous site (a precarious tower of javascript dependencies) were ported over.

The site's compiler is written in [Go](https://golang.org). It generates the static site you are reading by building a graph of entries from a bespoke markup language [[mmxup]].

Go's [templating engine](https://golang.org/pkg/text/template/) is lightly used, but it's often more straightforward to concat the html directly into a string.

If you want to learn more, see the [source code](https://github.com/mrshll/mmx).

# Building mmx

I've historically defaulted to the technology du jour to build and rebuild my personal and company webpages. It's worked fine - there have been moments of fantastic efficiency, and others of abysmal reverse engineering of an errant configuration or plugin.

In the spirit of [[Low tech]], and heavily inspired by others on the [webring](https://webring.xxiivv.com/), I sought out to simplify my dependencies, build something myself, and learn some new technologies in the process.

To start, I mapped out my hoped-for characteristics of the end result. I landed on:

+ low friction, so as to promote more writing;
+ low-level - no automagic frameworks or dependencies
+ long-term, with limited dependencies and render to straight html; and
+ extendable (I want this to be something that evolves with me over the next 5, 10, 15 years)

[Devine Lu Linvega's Oscean](https://github.com/XXIIVV/oscean/) served as the primary inspiration. I spent a weekend pouring over their wiki, the underlying C-code, and the ecosystem of file formats and tools they created. The act of reading through their code and reverse engineering the site compilation was one of the most fun weekends I've had in a while. I gained an understanding of their technical approach to linking and render disparate databases, while simultaneously exploring the content of those databases themselves.

I opted for GoLang as a learning opportunity. Compared to my day-to-day work languages of Javascript and Python, GoLang is miles closer to C. The last time I thought about pointers was my senior year of college, if that. I never imagined saying this, but I missed pointers! My impression of GoLang so far is fine. I hardly tapped into the features it is known for, such as concurrency. But this project served as a gentle introduction.

Another reason I wanted to create the compiler myself was so that I could add features over time that are typically only available on "hosted platforms," such as bidirectional linking and other memex-style data graph functions. I noticed that Oscean as well as others on the webring were able to do this. How cool.

The compiler is able to pull context from inbound links to pages. This is achieved by building a node tree when compiling the templating language to HTML. The correct `a` tag is located in the tree, and the node's parent content is pulled in as html.

In a moment of doubt, I played with [11ty](https://www.11ty.dev/) as well as [gatsby plugins](https://github.com/mathieudutour/gatsby-digital-garden) that promised functionality I sought. But after speed bumps with each, always grappling with the obfuscation that make them "magic," I felt confident that growing something from seed was the right path.

I've updated mmx to no longer rely on external dependencies (like markdown) outside of Golang's standard library. In order to do that while maintaining a nice writing experience (as much as I love html...), I created my own markdown language, {mmxup}.

Other features of established web frameworks are replicatable with bash, usually. For instance, "live reloading" is achieved with the following bash:

    #!/bin/sh
    bash build.sh
    while inotifywait -qqre modify ./src ./links ./data; do
      bash build.sh
    done

The site is hosted using Github pages. This is great, because there is no build step. I check in the built html files (in /doc) and they are served within seconds. CNAME setup was a breeze.

---

> The wood thrush, it is! Now I know
> who sings that clear arpeggio,
> three far notes weaving
> into the evening
> among leaves
> and shadow;
> 
> or at dawn in the woods, I've heard
> the sweet ascending triple word
> echoing over
> the silent river —
> but never
> seen the bird.

_Learning the Name by Ursula K. Le Guin_

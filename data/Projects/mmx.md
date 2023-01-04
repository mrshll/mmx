2023-01-03
# This site is called mmx

![A person building something with bamboo, from The Red Turtle](img/redturtle020.jpg)

**mmx** was started early 2021. Largely inspired by others on [webring](https://webring.xxiivv.com/), especially [Devine](https://wiki.xxiivv.com/site/home.html), the site aims to be a long term repository for my writing, notes, and research. Entries on my previous site (a precarious tower of javascript dependencies) were ported over.

Originally, the site's compiler is written in [Go](https://golang.org). At the start of 2023, it was simplified and rewritten in [[Lua]]. It generates the static site you are reading by building a graph of entries. During this rewrite, entries were moved from a bespoke markup language [[mmxup]] to markdown so that it was more interoperable with other tools and maintainable going forward. mmx depends only on Lua, [GNU date binary](https://www.gnu.org/software/coreutils/manual/html_node/date-invocation.html), and `find`.

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

Another reason I wanted to create the compiler myself was so that I could add features over time that are typically only available on "hosted platforms," such as bidirectional linking and other memex-style data graph functions.

The compiler is able to pull context from inbound links to pages. This is achieved by building a node tree when compiling the templating language to HTML. The correct `a` tag is located in the tree, and the node's parent content is pulled in as html.

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

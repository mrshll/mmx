## Battery

Framework laptops with Ubuntu don't have that great of battery performance out of the box, but there are a few ways you can improve it:

1. [Install TLP](https://linrunner.de/tlp/installation/ubuntu.html). TLP will set good defaults on some kernel settings that effect battery life, without you having to tweak them yourself. After you install, make sure to double check that it's not getting stomped on by the [power-profiles-daemon](https://linrunner.de/tlp/faq/installation.html#faq-ppd-conflict).
2. [Change your sleep mode](https://devnull.land/laptop-s2idle-to-deep) from `s2idle` to `deep`. By default laptops continue to use a lot of battery even when they're asleep. Switching to deep sleep will increase the time you can leave your laptop unplugged and shut. This will make your computer take longer to wake up.
3. [enable hibernate or sleep-then-hibernate](https://luisartola.com/solving-the-framework-laptop-battery-drain/). This is even better than deep sleep, you can probably leave your computer unplugged and asleep for weeks without issue. Depending on how much swap space you allocated when you installed Ubuntu, this might require resizing your swap partition, and it'll also make your computer take longer to wake up.

[More suggestions if you want to go further down the rabbit hole](https://community.frame.work/t/linux-battery-life-tuning/6665)

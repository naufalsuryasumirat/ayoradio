# ayoradio

Ayoradio is an always on radio that can be ran on a raspberry pi to detect a person's presence by detecting registered devices on the same network and plays radio on a playback device (speaker).

Also will send magic packets given the target device has Wake-on-LAN feature enabled. Will send this magic packet once per day when a registered device's presence is detected.

Running this will also run a simple page to control the playback device and the current playing song.

Visit [the site](https://radio.naufalsuryasumirat.uk/).

Tested to be running on RaspberryPi 5 (Debian Bookworm).

## configuration/prerequisites

Set the correct network interface in the .env file.

Make sure mpv, arp-scan, latest yt-dlp is installed in the system itself.

## footnote

Please don't randomly change the song I'm playing, thank you!

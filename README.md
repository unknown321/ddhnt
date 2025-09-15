ddhnt
====

Pack/unpack [Deadhunt (2005)](https://store.steampowered.com/app/435250) pk5 archives.

### Usage

Unpack: `./ddhnt-linux-amd64 Data.pk5`

Pack: `./ddhnt-linux-amd64 Data.pk5.json`

### Modding notes

You don't have to repack the game, move pk5 contents to Data directory, so it would be "Deadhunt/Data/Base"
and "Deadhunt/Data/Sound". Edit GtEngine.ini, set GAME->UsePaqFile to 0. The game will read files from disk as is.

A bunch of csvs, scenarios (missions) are... complicated: changed amount of runes from 5 to 55 in
first mission, gui updated with new amount, but they still stopped spawning after 5 pickups.

Weapons can be partially modified - swapped UZI sound id to pistol one, but magazine size never changed - hardcoded in exe?

There is a dev mode: edit PrivateProfile.csv, set "DeveloperMode" to 1, press F1-F12 for some developer
commands (F4 turns on stats overlay). You can also somehow turn on engine log (GtEngine.log); I did it by patching the exe:
replace `0f 86 d7 02 00 00` at 0x0040cbfb with `66 90 66 90 66 90` and create `GtEngine.log`. Log contents are encoded in CP1251.

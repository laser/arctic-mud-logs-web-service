# arctic-logs-webservice

> An ArcticMUD log parser and server

This project contains the log parser and web service for [https://arctic-mud-logs.com/](https://arctic-mud-logs.com/).

## About

These logs were collected from a variety of places:

- [http://www.angelfire.com/az/jdawg/clan.html](http://web.archive.org/web/20020328062804/http://www.angelfire.com/az/jdawg/clan.html)
- [http://tirannon.tky.hut.fi/arctic/](http://web.archive.org/web/20061214213442/http://tirannon.tky.hut.fi/arctic/)
- [http://davidwees.com/arcticlogs/](http://davidwees.com/arcticlogs/)
- [http://www.normstorm.com/arctic/](http://web.archive.org/web/20061216125832/http://www.normstorm.com/arctic/)

## Contributing

If you have a log file which you'd like added, create a GitHub issue and attach
said log file to the issue. Same goes for UX or bug reports - just create an
issue and I'll review at some point in the future.

## Development

### Full Usage

1. Drop logs into the `logs` directory
2. `make build` to build the parser (requires Go)
3. `make parse` to generate metadata files for each log
4. `make server` to run the web server
5. `open localhost:5555` to view your handywork

### Log Parser

```txt
The Planning Room
<AMA> Bolsvhik the Ballsy Blundering Gnome is standing here. (flying)
<BG> Drood Dood, Dewed Druid is standing here. (flying)
. . . is surrounded by a shimmering light.
Lord Atin von Oalbor plans his next attack upon the Garnet Mountains.
An elite draconian guard looks around for any vile intruders.
Lord Atin von Oalbor misses you with his pierce.
An elite sivak guard tries to punch you, but his fist passes right
through your image.
387H 264V 1X 23.15% 613C [Eclipz:Perfect] [Lord Atin von
Oalbor:Perfect] Mem:2 Exits:E>
An elite sivak guard gives you a solid punch in the face.
379H 264V 1X 23.15% 613C [Eclipz:V.Good] [Lord Atin von Oalbor:Perfect]
Mem:2 Exits:E>
Lord Atin von Oalbor focuses harshly on you and utters the words,
'dispel magic'.
You feel weaker.
379H 264V 1X 23.15% 613C [Eclipz:V.Good] [Lord Atin von Oalbor:Perfect]
Mem:2 Exits:E>
Doct tells your group 'hehe'
```

```
$ cat file.txt | parser
{"player_names":["Doct"],"clan_names":["AMA","BG"]}
```

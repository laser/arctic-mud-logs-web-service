# arctic-logs-webservice

> An ArcticMUD log parser and server

## Full Usage

1. Drop logs into the `logs` directory
2. `make build` to build the parser (requires Go)
3. `make parse` to generate metadata files for each log
4. `make server` to run the web server
5. `open localhost:5555` to view your handywork

## Log Parser

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

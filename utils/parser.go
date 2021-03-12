package main

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"sort"

	"github.com/laser/arctic-logs-webservice/types"
)

var compiled = []*regexp.Regexp{
	// communication / emotes
	regexp.MustCompile(`^(\w+) issues the order .*\.`),
	regexp.MustCompile(`^(\w+) kneels down and studies the ground`),
	regexp.MustCompile(`^(\w+) looks at .*\.`),
	regexp.MustCompile(`^(\w+) puts away .* tablets with a content`),
	regexp.MustCompile(`^(\w+) recites a`),
	regexp.MustCompile(`^(\w+) reluctantly stands up`),
	regexp.MustCompile(`^(\w+) says .*`),
	regexp.MustCompile(`^(\w+) tells you .*`),
	regexp.MustCompile(`^(\w+) tells your group.*`),
	regexp.MustCompile(`^(\w+) utters some strange words`),
	regexp.MustCompile(`^You tell (\w+) .*`),
	// moving around
	regexp.MustCompile(`^You follow (\w+) \w+\.`),
	regexp.MustCompile(`^(\w+) flies in from .*\.`),
	regexp.MustCompile(`^(\w+) dismounts\.`),
	regexp.MustCompile(`^(\w+) flies (up|down|west|east|south|north)\.`),
	regexp.MustCompile(`^(\w+) arrives from the (west|east|south|north)\.`),
	regexp.MustCompile(`^(\w+) arrives from (below|above)\.`),
	regexp.MustCompile(`^(\w+) leaves .*\.`),
	regexp.MustCompile(`^(\w+) stands up.*\.`),
	regexp.MustCompile(`^(\w+) panics, and attempts to flee\.`),
	// status
	regexp.MustCompile(`^(\w+) is dead`),
	regexp.MustCompile(`^(\w+) seems to be blinded`),
	regexp.MustCompile(`^(\w+) is paralyzed`),
	regexp.MustCompile(`^(\w+) freezes in place`),
	regexp.MustCompile(`^(\w+) comes out of hiding`),
	regexp.MustCompile(`^(\w+) slowly fades into existence`),
	regexp.MustCompile(`^(\w+) sinks gently to the ground`),
	// miscellaneous combat
	regexp.MustCompile(`^(\w+) aims a magic missile at .*, who `),
	regexp.MustCompile(`^(\w+) barrels into .*, knocking`),
	regexp.MustCompile(`^(\w+) charges at .*\.`),
	regexp.MustCompile(`^(\w+) crashes into .*`),
	regexp.MustCompile(`^(\w+) forces .* to the ground`),
	regexp.MustCompile(`^(\w+) kneels with a look of humility and prays to the gods`),
	regexp.MustCompile(`^(\w+) yells and leaps into the fray`),
	regexp.MustCompile(`^(\w+) cries out in pain as .* grabs .*\.`),
	regexp.MustCompile(`^(\w+) focuses harshly .*`),
	regexp.MustCompile(`^(\w+) forces .* to the ground harshly\.`),
	regexp.MustCompile(`^(\w+) is crushed by a wall of fire`),
	regexp.MustCompile(`^(\w+) is shredded by shards of ice`),
	regexp.MustCompile(`^(\w+) is stunned momentarily, recovering quickly`),
	regexp.MustCompile(`^(\w+) steps aside as .*`),
	regexp.MustCompile(`^(\w+) stops fighting .*\.`),
	regexp.MustCompile(`^(\w+) tries to barrel into .*`),
	regexp.MustCompile(`^cries out in pain as (\w+) grabs .*\.`),
	regexp.MustCompile(`^forces (\w+) to the ground harshly\.`),
	regexp.MustCompile(`^You attempt to assist (\w+)\.`),
	regexp.MustCompile(`^You crash into (\w+) in a bone.*`),
	regexp.MustCompile(`^You miss (\w+) and destroy an image instead`),
	regexp.MustCompile(`^You stop fighting (\w+)\.`),
	regexp.MustCompile(`^You try to bash (\w+), but.*`),
	regexp.MustCompile(`^Huge flames burn (\w+) from above`),
	// melee damage
	regexp.MustCompile(`^(\w+) .* .* (?:very|extremely) hard\.`),
	regexp.MustCompile(`^(\w+) (?:misses|massacres|obliterates|annihilates) .*\.`),
	regexp.MustCompile(`^(\w+) bludgeons .*`),
	regexp.MustCompile(`barely (?:crushes|cleaves|stabs|slashes|pierces) (\w+)\.`),
	regexp.MustCompile(`(?:massacres|obliterates|annihilates) (\w+) with .*\.`),
	regexp.MustCompile(`(?:crushes|cleaves|stabs|slashes|pierces) (\w+) very hard\.`),
	regexp.MustCompile(`You (?:crush|cleave|stab|slash|pierce) (\w+) very hard\.`),
	regexp.MustCompile(`You (?:crush|cleave|stab|slash|pierce) (\w+) hard\.`),
	regexp.MustCompile(`You (?:crush|cleave|stab|slash|pierce) (\w+)\.`),
	regexp.MustCompile(`You (?:miss|massacre|obliterate|annihilate) (\w+) with .*\.`),
	// looting
	regexp.MustCompile(`^(\w+) gets .* from the corpse of`),
	regexp.MustCompile(`^(\w+) is zapped`),
	regexp.MustCompile(`gets a .* from the corpse of (\w+)\.`),
	regexp.MustCompile(`gets an .* from the corpse of (\w+)\.`),
	regexp.MustCompile(`You get a .* from the corpse of (\w+)\.`),
	regexp.MustCompile(`You get an .* from the corpse of (\w+)\.`),
	// manual add
	regexp.MustCompile(`Manual add: (\w+)\.`),
}

var clan = regexp.MustCompile(`^<(\w+)> .*here`)

var blacklist = []*regexp.Regexp{
	regexp.MustCompile(`^A$`),
	regexp.MustCompile(`^someone$`),
	regexp.MustCompile(`^Someone$`),
	regexp.MustCompile(`^YOU$`),
	regexp.MustCompile(`^You$`),
	regexp.MustCompile(`^your$`),
	regexp.MustCompile(`\d`),
	regexp.MustCompile(`^[a-z].*`),
}

func main() {
	playerNamesSet := make(map[string]interface{})
	clanNamesSet := make(map[string]interface{})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		m1 := clan.FindStringSubmatch(line)
		if len(m1) > 0 {
			clanNamesSet[m1[1]] = struct{}{}
		}

		for idx := range compiled {
			m2 := compiled[idx].FindStringSubmatch(line)
			if len(m2) > 0 {
				ok := true
				for idx := range blacklist {
					m3 := blacklist[idx].FindStringSubmatch(m2[1])
					if !ok || len(m3) != 0 {
						ok = false
					}
				}

				if ok {
					playerNamesSet[m2[1]] = struct{}{}
				}
			}
		}
	}

	check(scanner.Err())

	m := types.Meta{
		PlayerNames: []string{},
		ClanNames:   []string{},
	}

	playerNamesSorted := make([]string, 0)
	clanNamesSorted := make([]string, 0)

	for idx := range clanNamesSet {
		clanNamesSorted = append(clanNamesSorted, idx)
	}

	for idx := range playerNamesSet {
		playerNamesSorted = append(playerNamesSorted, idx)
	}

	sort.Strings(clanNamesSorted)
	sort.Strings(playerNamesSorted)

	m.ClanNames = clanNamesSorted
	m.PlayerNames = playerNamesSorted

	b, err := json.Marshal(&m)
	check(err)

	_, err = os.Stdout.Write(b)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

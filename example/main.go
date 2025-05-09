package main

import (
	"fmt"

	"github.com/legamerdc/gas"
	"github.com/legamerdc/gas/ds"
)

func main() {
	var (
		u1 = unit{
			Id:      1,
			Hp:      100,
			Gas:     gas.NewGAS[*world, *unit, *event, int64](),
			Faction: 1,
		}
		u2 = unit{
			Id: 2,
			X:  10, Y: 0,
			Hp:      100,
			Gas:     gas.NewGAS[*world, *unit, *event, int64](),
			Faction: 2,
		}
		u3 = unit{
			Id: 3,
			X:  0, Y: 10,
			Hp:      100,
			Gas:     gas.NewGAS[*world, *unit, *event, int64](),
			Faction: 1,
		}
		u4 = unit{
			Id: 4,
			X:  10, Y: 10,
			Hp:      100,
			Gas:     gas.NewGAS[*world, *unit, *event, int64](),
			Faction: 2,
		}
		w = world{
			units: map[int64]*unit{1: &u1, 2: &u2, 3: &u3, 4: &u4},
		}
	)
	u1.Gas.AddAbility(&w, &u1, ds.LookupPtr[abilityConfigI](1).Activate(&w, &u1))
	u2.Gas.AddAbility(&w, &u2, ds.LookupPtr[abilityConfigI](2).Activate(&w, &u2))
	for i := 0; i < 2000; i += 50 {
		w.SetNow(int64(i))
		fmt.Println("tick ", i, ":")
		if i == 200 {
			printErr(u1.Gas.Cast(&w, &u1, 0, 1))
			printErr(u2.Gas.Cast(&w, &u2, 0, 2))
		}
		w.Tick()
	}
}

func printErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

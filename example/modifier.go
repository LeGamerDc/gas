package main

import (
	"fmt"

	"github.com/legamerdc/gas/ds"
)

type modifierType int32

const (
	modifierNone = modifierType(iota)
	modifierDamage
	modifierHeal
	modifierAddRunning
)

type ModifierConfig struct {
	ID      int32                    `json:"id"`
	Type    modifierType             `json:"type"`
	Value   int                      `json:"value"`
	Running ds.Proxy[runningConfigI] `json:"running"`
}

func (c *ModifierConfig) Id() int32 {
	return c.ID
}

func (c *ModifierConfig) Apply(w *world, u *unit) {
	switch c.Type {
	case modifierDamage:
		fmt.Printf("    unit %d: damage %d->%d\n", u.Id, u.Hp, u.Hp-c.Value)
		u.Hp -= c.Value
	case modifierHeal:
		fmt.Printf("    unit %d: heal %d->%d\n", u.Id, u.Hp, u.Hp+c.Value)
		u.Hp += c.Value
	case modifierAddRunning:
		r := c.Running.Get().Activate(w, u)
		u.Gas.AddRunning(w, u, r)
		fmt.Printf("    unit %d: add running %d\n", u.Id, r.Id())
	default:
	}
}

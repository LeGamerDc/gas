package main

import (
	"math"
	"slices"
)

var _ selectorConfigI = (*SelectorConfig)(nil)

type SelectorConfig struct {
	ID          int32   `json:"id"`
	Range       float32 `json:"range"`
	Count       int     `json:"count"`
	AllyOrEnemy bool    `json:"ally"`
	Self        bool    `json:"self"`
}

func (c *SelectorConfig) Id() int32 {
	return c.ID
}

func (c *SelectorConfig) Select(w *world, u *unit, exclude []int64) (pick []*unit) {
	if c.Self {
		return []*unit{u}
	}
	if c.Count <= 0 {
		return
	}
	for _, v := range w.units {
		if ((c.AllyOrEnemy && u.Faction == v.Faction) ||
			(!c.AllyOrEnemy && u.Faction != v.Faction)) &&
			distance(u.X, u.Y, v.X, v.Y) < c.Range &&
			!slices.Contains(exclude, v.Id) {
			pick = append(pick, v)
		}
		if len(pick) >= c.Count {
			return
		}
	}
	return
}

func distance(x1, y1, x2, y2 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)))
}

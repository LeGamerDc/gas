package main

import (
	"errors"

	"github.com/legamerdc/gas"
	"github.com/legamerdc/gas/ds"
)

var (
	errAbilityCd       = errors.New("ability in cd")
	errAbilityPickNone = errors.New("ability no target picked")
)

type SimpleAbilityConfig struct {
	ID       int32                     `json:"id"`
	Cd       int64                     `json:"cd"`
	Selector ds.Proxy[selectorConfigI] `json:"selector"`
	Modifier ds.Proxy[modifierConfigI] `json:"modifier"`
}

func (c *SimpleAbilityConfig) Id() int32 {
	return c.ID
}

func (c *SimpleAbilityConfig) Activate(w *world, u *unit) gas.AbilityI[*world, *unit, *event, int64] {
	return &SimpleAbility{
		c:    c,
		next: 0,
	}
}

var _ gas.AbilityI[*world, *unit, *event, int64] = (*SimpleAbility)(nil)

type SimpleAbility struct {
	c    *SimpleAbilityConfig
	next int64
}

func (c *SimpleAbility) Cast(w *world, u *unit, _target int64) error {
	now := w.Now()
	if now < c.next {
		return errAbilityCd
	}
	if targets := c.c.Selector.Get().Select(w, u, nil); len(targets) > 0 {
		for _, v := range targets {
			c.c.Modifier.Get().Apply(w, v)
		}
		return nil
	}
	return errAbilityPickNone
}

func (c *SimpleAbility) ListenEvent() []gas.EventKind {
	return nil
}

func (c *SimpleAbility) OnCreate(w *world, u *unit) {

}

func (c *SimpleAbility) OnEvent(w *world, u *unit, e *event) {}

func (c *SimpleAbility) Id() int32 {
	return c.c.ID
}

type ChainLikeRunningConfig struct {
	ID       int32                     `json:"id"`
	Gap      int64                     `json:"gap"`
	Duration int64                     `json:"duration"`
	Selector ds.Proxy[selectorConfigI] `json:"selector"`
	Modifier ds.Proxy[modifierConfigI] `json:"modifier"`
	Repeat   bool                      `json:"repeat"`
}

func (c *ChainLikeRunningConfig) Id() int32 {
	return c.ID
}

func (c *ChainLikeRunningConfig) Activate(w *world, u *unit) gas.RunningI[*world, *unit, *event] {
	return &ChainLikeRunning{
		cfg: c,
		x:   u.X, y: u.Y,
	}
}

var _ gas.RunningI[*world, *unit, *event] = (*ChainLikeRunning)(nil)

type ChainLikeRunning struct {
	cfg       *ChainLikeRunningConfig
	x, y      float32
	next, end int64
	visited   []int64
}

func (c *ChainLikeRunning) Id() int32 {
	return c.cfg.ID
}

func (c *ChainLikeRunning) ListenEvent() []gas.EventKind {
	return nil
}

func (c *ChainLikeRunning) Stack() (int64, int64) {
	return 0, c.end
}

func (c *ChainLikeRunning) OnStack(_, end int64) {
	c.end = max(c.end, end)
}

func (c *ChainLikeRunning) Think(w *world, u *unit) int64 {
	now := w.Now()
	if now > c.end {
		return gas.Never
	}
	if now >= c.next {
		c.next += c.cfg.Gap
		targets := c.cfg.Selector.Get().Select(w, u, c.visited)
		for _, v := range targets {
			c.cfg.Modifier.Get().Apply(w, v)
			if !c.cfg.Repeat {
				c.visited = append(c.visited, v.Id)
			}
		}
	}
	return min(c.next, c.end)
}

func (c *ChainLikeRunning) OnBegin(w *world, u *unit) int64 {
	now := w.Now()
	c.next = now
	c.end = now + c.cfg.Duration
	return c.next
}

func (c *ChainLikeRunning) OnEnd(w *world, u *unit) {}

func (c *ChainLikeRunning) OnEvent(w *world, u *unit, e *event) {}

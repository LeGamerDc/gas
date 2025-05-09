package gas

import (
	"errors"
	"slices"
)

func NewGAS[W WI, U UI, E EI, T any]() *GAS[W, U, E, T] {
	gas := &GAS[W, U, E, T]{
		watchEvent: map[EventKind]struct{}{},
	}
	gas.Abilities.Init(4)
	gas.Running.Init(4)
	gas.Buff.Init(4)
	return gas
}

func (g *GAS[W, U, E, T]) Think(w W, u U) int64 {
	var (
		now  = w.Now()
		next = now + ThinkLater
	)

	for g.Running.Size() > 0 {
		i, _, v, when := g.Running.Top()
		if when > now {
			next = when
			break
		}
		if next1 := v.Think(w, u); next1 >= 0 {
			g.Running.Update(i, max(now+MinThinkGap, next1))
		} else {
			v.OnEnd(w, u)
			g.Running.Remove(i)
		}
	}
	next = min(next, g.thinkBuff(now, u))
	return next
}

func (g *GAS[W, U, E, T]) OnEvent(w W, u U, e E) {
	if _, ok := g.watchEvent[e.Kind()]; !ok {
		return
	}
	g.Abilities.Iter(func(a AbilityI[W, U, E, T]) {
		if slices.Contains(a.ListenEvent(), e.Kind()) {
			a.OnEvent(w, u, e)
		}
	})
	g.Running.Iter(func(v RunningI[W, U, E]) {
		if slices.Contains(v.ListenEvent(), e.Kind()) {
			v.OnEvent(w, u, e)
		}
	})
}

func (g *GAS[W, U, E, T]) Cast(w W, u U, t T, id int32) error {
	if i, v := g.Abilities.Get(id); i >= 0 {
		return v.Cast(w, u, t)
	}
	return errors.New("ability not found")
}

func (g *GAS[W, U, E, T]) AddAbility(w W, u U, a AbilityI[W, U, E, T]) bool {
	if i, _ := g.Abilities.Get(a.Id()); i >= 0 {
		return false
	}
	a.OnCreate(w, u)
	g.Abilities.Put(a.Id(), a)
	for _, x := range a.ListenEvent() {
		g.watchEvent[x] = struct{}{}
	}
	return true
}

func (g *GAS[W, U, E, T]) AddRunning(w W, u U, r RunningI[W, U, E]) {
	if i, v := g.Running.Get(r.Id()); i >= 0 {
		a, b := r.Stack()
		v.OnStack(a, b)
		return
	}
	if next := r.OnBegin(w, u); next >= 0 {
		g.Running.Put(r.Id(), r, next)
	}
}

func (g *GAS[W, U, E, T]) Watch(kind EventKind, id int32) {

}

func (g *GAS[W, U, E, T]) UnWatch(kind EventKind, id int32) {

}

package gas

import "slices"

func (g *GAS[W, U, E]) Think(w W, u U) int64 {
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

func (g *GAS[W, U, E]) OnEvent(w W, u U, e E) {
	if _, ok := g.watchEvent[e.Kind()]; !ok {
		return
	}
	g.Abilities.Iter(func(_ int32, a AbilityI[W, U, E]) (stop bool) {
		if slices.Contains(a.ListenEvent(), e.Kind()) {
			a.OnEvent(w, u, e)
		}
		return false
	})
	g.Running.Iter(func(_ int32, v RunningI[W, U, E]) (stop bool) {
		if slices.Contains(v.ListenEvent(), e.Kind()) {
			v.OnEvent(w, u, e)
		}
		return false
	})
}

func (g *GAS[W, U, E]) AddAbility(w W, u U, a AbilityI[W, U, E]) bool {
	if i, _ := g.Abilities.Get(a.Id()); i >= 0 {
		return false
	}
	a.OnCreate(w, u)
	g.Abilities.Push(a.Id(), a)
	for _, x := range a.ListenEvent() {
		g.watchEvent[x] = struct{}{}
	}
	return true
}

func (g *GAS[W, U, E]) AddRunning(w W, u U, r RunningI[W, U, E]) {
	if i, v := g.Running.Get(r.Id()); i >= 0 {
		a, b := r.Stack()
		v.OnStack(a, b)
		return
	}
	next := r.OnBegin(w, u)
	g.Running.Push(r.Id(), r, next)
}

func (g *GAS[W, U, E]) refreshWatch() {
	clear(g.watchEvent)
	g.Abilities.Iter(func(_ int32, a AbilityI[W, U, E]) (stop bool) {
		for _, x := range a.ListenEvent() {
			g.watchEvent[x] = struct{}{}
		}
		return false
	})
	g.Running.Iter(func(_ int32, v RunningI[W, U, E]) (stop bool) {
		for _, x := range v.ListenEvent() {
			g.watchEvent[x] = struct{}{}
		}
		return false
	})
}

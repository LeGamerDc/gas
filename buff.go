package gas

type (
	BuffKind    int32
	BuffCompose int32
	BuffStack   int32
)

const (
	BuffComposeNone BuffCompose = iota
	BuffComposeAdd
	BuffComposePercent
	BuffComposeMagnify
)

const (
	BuffStackNone BuffStack = iota
	BuffStackGreater
	BuffStackLonger
	BuffStackSeparate
)

type BuffNode struct {
	Source string
	Value  float64
	Expire int64
	Kind   BuffKind
}

type BuffList struct {
	Nodes   []BuffNode
	kind    BuffKind
	compose BuffCompose
	stack   BuffStack
}

func (b *BuffList) mergeBuff(n BuffNode) int64 {
	if b.stack != BuffStackSeparate {
		for i := range b.Nodes {
			if b.Nodes[i].Source == n.Source {
				switch b.stack {
				case BuffStackGreater:
					if b.Nodes[i].Value < n.Value {
						b.Nodes[i] = n
					}
				case BuffStackLonger:
					if b.Nodes[i].Expire < n.Expire {
						b.Nodes[i] = n
					}
				default:
				}
				return b.next()
			}
		}
	}
	b.Nodes = append(b.Nodes, n)
	return b.next()
}

func (b *BuffList) updateNone(now int64) (v float64) {
	n := len(b.Nodes) - 1
	for i := 0; i <= n; {
		if b.Nodes[i].Expire > now {
			v = 1
			i++
		} else {
			b.Nodes[i] = b.Nodes[n]
			n--
		}
	}
	b.Nodes = b.Nodes[:n]
	return
}

func (b *BuffList) updateAdd(now int64, base float64) (v float64) {
	n := len(b.Nodes) - 1
	for i := 0; i <= n; {
		if b.Nodes[i].Expire > now {
			v += b.Nodes[i].Value
			i++
		} else {
			b.Nodes[i] = b.Nodes[n]
			n--
		}
	}
	b.Nodes = b.Nodes[:n]
	return base + v
}

func (b *BuffList) updatePercent(now int64, base float64) (v float64) {
	n := len(b.Nodes) - 1
	for i := 0; i <= n; {
		if b.Nodes[i].Expire > now {
			v += b.Nodes[i].Value
			i++
		} else {
			b.Nodes[i] = b.Nodes[n]
			n--
		}
	}
	b.Nodes = b.Nodes[:n]
	return base * (1 + v)
}

func (b *BuffList) updateMagnify(now int64, base float64) (v float64) {
	n := len(b.Nodes) - 1
	v = 1
	for i := 0; i <= n; {
		if b.Nodes[i].Expire > now {
			v *= 1 + b.Nodes[i].Value
			i++
		} else {
			b.Nodes[i] = b.Nodes[n]
			n--
		}
	}
	b.Nodes = b.Nodes[:n]
	return base * v
}

func (b *BuffList) next() int64 {
	x := Never
	for _, n := range b.Nodes {
		if x == Never || n.Expire < x {
			x = n.Expire
		}
	}
	return x
}

func (g *GAS[W, U, E]) calculateBuff(now int64, v *BuffList, u U) {
	base := u.GetBuffBase(v.kind)
	newV := base
	switch v.compose {
	case BuffComposeNone:
		newV = v.updateNone(now)
	case BuffComposeAdd:
		newV = v.updateAdd(now, base)
	case BuffComposePercent:
		newV = v.updatePercent(now, base)
	case BuffComposeMagnify:
		newV = v.updateMagnify(now, base)
	default: // do nothing
	}
	u.SetBuff(v.kind, newV)
}

func (g *GAS[W, U, E]) thinkBuff(now int64, u U) int64 {
	for g.Buff.Size() > 0 {
		i, _, v, when := g.Buff.Top()
		if when > now {
			return when
		}
		g.calculateBuff(now, v, u)
		if next := v.next(); next >= 0 {
			g.Buff.Update(i, next)
		} else {
			g.Buff.Remove(i)
		}
	}
	return now + ThinkLater
}

func (g *GAS[W, U, E]) AddBuff(w W, u U, b BuffNode) {
	i, v := g.Buff.Get(b.Kind)
	if i >= 0 {
		next := v.mergeBuff(b)
		g.Buff.Update(i, next)
		g.calculateBuff(w.Now(), v, u)
		return
	}
	c, s := w.DescribeBuffKind(b.Kind)
	v = &BuffList{
		Nodes:   []BuffNode{b},
		kind:    b.Kind,
		compose: c,
		stack:   s,
	}
	g.Buff.Push(b.Kind, v, v.next())
	g.calculateBuff(w.Now(), v, u)
}

package main

import "github.com/legamerdc/gas"

type (
	world struct {
		nowTs int64 // millisecond
		units map[int64]*unit
	}

	unit struct {
		Id      int64
		X, Y    float32
		Hp      int
		Gas     *gas.GAS[*world, *unit, *event, int64]
		Faction int32
	}

	event struct{}

	abilityConfigI interface {
		Id() int32
		Activate(w *world, u *unit) gas.AbilityI[*world, *unit, *event, int64]
	}

	runningConfigI interface {
		Id() int32
		Activate(w *world, u *unit) gas.RunningI[*world, *unit, *event]
	}

	selectorConfigI interface {
		Id() int32
		Select(w *world, u *unit, exclude []int64) (pick []*unit)
	}

	modifierConfigI interface {
		Id() int32
		Apply(w *world, u *unit)
	}
)

func (w *world) SetNow(nowTs int64) {
	w.nowTs = nowTs
}

func (w *world) Tick() {
	for _, u := range w.units {
		u.Gas.Think(w, u)
	}
}

func (w *world) Now() int64 {
	return w.nowTs
}

func (w *world) DescribeBuffKind(kind gas.BuffKind) (gas.BuffCompose, gas.BuffStack) {
	// 这个例子不涉及buff
	return 0, 0
}

func (w *world) hurt() {

}

func (u *unit) GetBuffBase(kind gas.BuffKind) float64 {
	return 0
}

func (u *unit) SetBuff(kind gas.BuffKind, f float64) {}

func (e *event) Kind() gas.EventKind {
	return 0
}

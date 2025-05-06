package gas

import "github.com/legamerdc/gas/ds"

/*
type parameter explain:
W: game world
U: game unit which attached to
E: game event
*/

const (
	// ThinkLater 保底think，避免技能实体长时间不think
	ThinkLater int64 = 3000
	// MinThinkGap 最低think间隔，避免技能实体无限think
	MinThinkGap int64 = 10
	// Never 表示再也不会think (aka 销毁了)
	Never int64 = -1
)

type (
	EventKind int32

	// WI 世界抽象，开发实现时能通过 WI 影响整个世界，如对单位施加效果，添加buff/running等。
	WI interface {
		Now() int64
		DescribeBuffKind(BuffKind) (BuffCompose, BuffStack)
	}
	// UI 单位抽象，开发实现时能通过 UI 影响自己或获取信息，如提取攻击目标。
	UI interface {
		GetBuffBase(BuffKind) float64
		SetBuff(BuffKind, float64)
	}
	// EI 事件抽象，我们无法在库中知道事件的具体类型。
	EI interface {
		Kind() EventKind
	}

	// AbilityI 技能抽象，玩家拥有的技能，开发者可以在OnEvent实现中去管理cd、判断时机和影响世界。
	AbilityI[W WI, U UI, E EI] interface {
		Id() int32
		ListenEvent() []EventKind
		OnCreate(W, U)
		OnEvent(W, U, E)
	}

	// RunningI 运行时实体，主要包括3种接口
	// Stack/OnStack 堆叠方案，一些运行时实体不同的堆叠层数造成的效果不一样，开发者可以在这里实现这个实体如何堆叠。
	// OnEvent 事件触发，一些运行时实体可能受到外部事件触发，开发者可以在这里实现。
	// Think 主动控制，一些运行时实体可能周期性地去触发一些行为，检查状态影响自身，开发者可以在这里实现，Think 返回的是下次进入Think的时间
	RunningI[W WI, U UI, E EI] interface {
		Id() int32
		ListenEvent() []EventKind
		Stack() (int64, int64)
		OnStack(_, _ int64)
		Think(W, U) int64
		OnBegin(W, U) int64
		OnEnd(W, U)
		OnEvent(W, U, E)
	}

	// GAS game ability system: 一个单位身上的技能管理框架
	// ArrayMap 和 HeapArrayMap 可以看做使用slice实现的map。在元素较少时，访问性能与map接近，但遍历效率大幅优于map。
	GAS[W WI, U UI, E EI] struct {
		Abilities ds.ArrayMap[int32, AbilityI[W, U, E]]
		Running   ds.HeapArrayMap[int32, int64, RunningI[W, U, E]]
		Buff      ds.HeapArrayMap[BuffKind, int64, *BuffList]

		// watchEvent TODO 后续使用count维护watch
		watchEvent map[EventKind]struct{}
	}
)

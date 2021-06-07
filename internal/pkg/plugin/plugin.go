package plugin

type Lifecycle interface {
	Init()
	Destroy()
}

var plugins = make([]Lifecycle, 0)

func AddPlugin(plugin Lifecycle) {
	plugins = append(plugins, plugin)
}

func Get() []Lifecycle {
	return plugins
}

package factory

type Bean interface {
	Init()
	Destroy()
}

var (
	beanFactory = make(map[string]Bean)
)

func Register(name string, bean Bean) {
	beanFactory[name] = bean
}

func GetAllBeans() []Bean {
	beans := make([]Bean, len(beanFactory))
	var i = 0
	for _, value := range beanFactory {
		beans[i] = value
		i++
	}
	return beans
}

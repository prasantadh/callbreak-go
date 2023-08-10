package callbreak

import "fmt"

var registry = make(map[string]func() Strategy)

func RegisterStrategy(name string, f func() Strategy) error {
	if _, ok := registry[name]; ok {
		return fmt.Errorf("strategy %s already exists", name)
	}
	registry[name] = f
	return nil
}

func GetStrategy(name string) (Strategy, error) {
	if _, ok := registry[name]; !ok {
		return nil, fmt.Errorf("strategy %s does not exist", name)
	}
	return registry[name](), nil
}

package resources

import (
	"easy-lyric/EasyLyric/resources/base_resources"
	"easy-lyric/EasyLyric/resources/kidung"
)

func Get(resourceName string) base_resources.Source {
	var resource base_resources.Source = nil

	switch resourceName {
	case kidung.ResourceName:
		resource = kidung.Kidung
	default:
		return nil
	}
	return resource
}

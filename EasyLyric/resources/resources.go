package resources

import (
	"easy-lyric/EasyLyric/resources/base_resources"
	"easy-lyric/EasyLyric/resources/kidung"
	"easy-lyric/EasyLyric/resources/liriklagurohani"
	"easy-lyric/EasyLyric/resources/unlimited"
)

func Get(resourceName string) base_resources.Source {
	var resource base_resources.Source = nil

	switch resourceName {
	case kidung.ResourceName:
		resource = kidung.Kidung
	case unlimited.ResourceName:
		resource = unlimited.Unlimited
	case liriklagurohani.ResourceName:
		resource = liriklagurohani.LirikLagu
	default:
		return nil
	}
	return resource
}

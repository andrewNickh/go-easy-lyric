package config

import "easy-lyric/util/log"

func Log() {
	if Instance.LogLevel != "" {
		l, err := log.New(
			Instance.LogLevel,
			Instance.LogDir,
			log.LstdFlags|log.Lshortfile,
			"",
			true,
			Instance.Env,
			"easy_lyric_1")

		if err != nil {
			panic(err)
		}
		log.Export(l)
	}
}

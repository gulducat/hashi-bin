package util

import "os"

type ExtraOptions struct {
	All  bool
	Beta bool
	Ent  bool
}

func GetOptions() ExtraOptions {
	// feels a bit jank but it works..
	// TODO: add these args to -h somehow...
	return ExtraOptions{
		All:  os.Getenv("HASHI_ALL") != "" || InArray(os.Args, "-all"),
		Beta: os.Getenv("HASHI_BETA") != "" || InArray(os.Args, "-with-beta"),
		Ent:  os.Getenv("HASHI_ENTERPRISE") != "" || InArray(os.Args, "-only-enterprise"),
	}
}

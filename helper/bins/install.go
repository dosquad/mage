package bins

import "github.com/princjef/mageutil/bintool"

func Install(in *bintool.BinTool) error {
	if err := in.Ensure(); err != nil {
		return in.Install()
	}

	return nil
}

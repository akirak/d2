package d2cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"oss.terrastruct.com/util-go/xdefer"

	"oss.terrastruct.com/util-go/xmain"

	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2parser"
)

func fmtCmd(ctx context.Context, ms *xmain.State) (err error) {
	defer xdefer.Errorf(&err, "failed to fmt")

	ms.Opts = xmain.NewOpts(ms.Env, ms.Opts.Flags.Args()[1:])
	if len(ms.Opts.Args) == 0 {
		return xmain.UsageErrorf("fmt must be passed at least one file to be formatted")
	}

	for _, inputPath := range ms.Opts.Args {
		if inputPath != "-" {
			inputPath = ms.AbsPath(inputPath)
			d, err := os.Stat(inputPath)
			if err == nil && d.IsDir() {
				inputPath = filepath.Join(inputPath, "index.d2")
			}
		}

		input, err := ms.ReadPath(inputPath)
		if err != nil {
			return err
		}

		m, err := d2parser.Parse(inputPath, bytes.NewReader(input), nil)
		if err != nil {
			return err
		}

		output := []byte(d2format.Format(m))
		if !bytes.Equal(output, input) || inputPath == "-" {
			if err := ms.WritePath(inputPath, output); err != nil {
				return err
			}
		}
	}
	return nil
}

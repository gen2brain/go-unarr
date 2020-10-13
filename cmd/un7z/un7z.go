package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-unarr"
)

// version info as per https://goreleaser.com/customization/build/
var (
	version = "0.0.0-LOCAL"
	commit  = "g0000000"
	date    = "0000-00-00T00:00:00Z"
)

func usage() {
	fmt.Println("un7z xvf <filename>")
}

func main() {
	if 2 == len(os.Args) && "version" == os.Args[1] {
		fmt.Printf("un7z v%s %s (%s)\n", version, commit, date)
	}

	if len(os.Args) <= 1 {
		usage()
		os.Exit(1)
	}

	// We've got some funny business here to allow
	// this to have tar-like flags, while still (mostly)
	// using the go 'flag' package for parsing

	// look for first file
	var fileIndex = 0
	args := os.Args[1:]
	for i, arg := range args {
		if "-" == arg {
			fileIndex = i
			break
		}
		stat, err := os.Stat(arg)
		if nil != err {
			continue
		}
		ext := filepath.Ext(stat.Name())
		if stat.Mode().IsRegular() && len(ext) > 0 {
			fileIndex = i
			break
		}
	}

	// expand arguments like xvf into -x -v -f
	// (with a special exception for -C)
	var captureNext bool
	newArgs := make([]string, len(os.Args[:1]))
	_ = copy(newArgs, os.Args[:1])
	for i, arg := range args {
		if i >= fileIndex {
			break
		}

		if captureNext {
			captureNext = false
			newArgs = append(newArgs, arg)
			continue
		}

		// handle --extract
		if strings.HasPrefix(arg, "--") {
			newArgs = append(newArgs, arg)
			continue
		}

		// handle -xvf (turn into xvf)
		if "-C" == arg {
			newArgs = append(newArgs, arg)
			captureNext = true
			continue
		}

		// handle -xvf (turn into xvf)
		if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
		}

		if len(arg) <= 0 {
			newArgs = append(newArgs, "")
		}

		// handle xvf (turn into -x, -v, -f)
		// (special exception for -C)
		for _, c := range arg {
			newArgs = append(newArgs, "-"+string(c))
			if "C" == arg {
				captureNext = true
				break
			}
		}
	}
	os.Args = append(newArgs, os.Args[fileIndex+1:]...)

	// normal 'flag' parsing
	var extract bool
	var extractX bool
	var extractE bool
	var extractExtract bool
	var list bool
	var listT bool
	var listList bool
	var filename string
	var verboseV bool
	var verboseVerbose bool
	var verbose bool
	wd, err := os.Getwd()
	if nil != err {
		fmt.Fprintf(os.Stderr, "could not get working directory: %v\n", err)
		os.Exit(1)
	}
	flag.StringVar(&wd, "C", wd, "change destination directory")
	flag.BoolVar(&listT, "t", false, "list files in the archive")
	flag.BoolVar(&listList, "list", false, "alias for -t (list)")
	flag.BoolVar(&extractX, "x", false, "extract the given file")
	flag.BoolVar(&extractE, "e", false, "alias for -x (extract)")
	flag.BoolVar(&extractExtract, "extract", false, "alias for -x (extract)")
	flag.BoolVar(&verboseV, "v", false, "verbose mode")
	flag.BoolVar(&verboseVerbose, "verbose", false, "alias for -v (verbose)")
	flag.StringVar(&filename, "f", "", "path to archive file")
	flag.Parse()

	// if any of the verbose arguments are specified, verbose should be true
	verbose = verboseV || verboseVerbose

	// if any of the list arguments are specified, list should be true
	list = listT || listList
	extract = extractX || extractE || extractExtract

	if len(filename) <= 0 {
		filename = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		// TODO handle extracting specific files
		usage()
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("@%s\n", strings.Join(os.Args, " "))
	}

	if list {
		if extract {
			fmt.Fprintf(os.Stderr,
				"error: must list (-t) or extract (-x), not both")
			os.Exit(1)
		}
		if err := listArchive(filename); nil != err {
			fmt.Fprintf(os.Stderr, "failed to list %q: %v\n", filename, err)
			os.Exit(1)
		}
		return
	}

	// default is to extract
	if err := extractArchive(filename, wd, verbose); nil != err {
		fmt.Fprintf(os.Stderr, "failed to extract %q: %v\n", filename, err)
		os.Exit(1)
	}
	return
}

func listArchive(src string) error {
	a, err := unarr.NewArchive(src)
	if err != nil {
		return err
	}
	defer a.Close()

	filenames, err := a.List()
	if err != nil {
		return err
	}

	for _, name := range filenames {
		fmt.Println(name)
	}
	fmt.Println("total ", len(filenames))
	return nil
}

func extractArchive(src, dst string, verbose bool) error {
	a, err := unarr.NewArchive(src)
	if err != nil {
		return err
	}
	defer a.Close()

	filenames, err := a.Extract(dst)
	if nil != err {
		return err
	}
	if verbose {
		total := len(filenames)
		for _, name := range filenames {
			fmt.Println(name)
		}
		fmt.Println("total", total)
	}
	return nil
}

/*
func extractArchive(src, dst string, verbose bool) error {
	a, err := unarr.NewArchive(src)
	if err != nil {
		return err
	}
	defer a.Close()

	if verbose {
		fmt.Printf("extracting %q to %q ...\n", src, dst)
	}

	defer func() {
		// print trailing newline
		if verbose {
			fmt.Printf("\n")
		}
	}()

	total := 0
	for {
		if err := a.Entry(); nil != err {
			if io.EOF == err {
				break
			}
			return err
		}

		name := a.Name()
		sizeb := a.Size()
		if verbose {
			// x    0b name
			// x 1000b name2
			size := byteCountIEC(int64(sizeb))
			fmt.Printf("x %s\t%s ...", leftPad2Len(size, " ", 10), name)
		}

		dirname, err := secureJoin(dst, filepath.Dir(name))
		if nil != err {
			return err
		}
		if err := os.MkdirAll(dirname, 0755); nil != err {
			return err
		}

		dst := filepath.Join(dirname, filepath.Base(name))
		w, err := os.Create(dst)
		if nil != err {
			return err
		}
		defer w.Close()

		if err := copyFile(a, sizeb, w); nil != err {
			return err
		}
		total++

		if verbose {
			fmt.Printf("\n")
		}
	}

	if verbose {
		// no trailing \n, that will be added
		fmt.Printf("total %d", total)
	}

	return nil
}

func secureJoin(base, path string) (string, error) {
	abs, err := filepath.Abs(base)
	if nil != err {
		return "", err
	}

	cleanpath := filepath.Join(abs, path)
	// /path/to/evil => /path/to/evil/
	// /path/to/evildoer =>  /path/to/evildoer/
	if !strings.HasPrefix(cleanpath+"/", abs+"/") {
		return "", errors.New("insecure path")
	}

	return cleanpath, nil
}

// leftPad2Len https://github.com/DaddyOh/golang-samples/blob/master/pad.go
func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

// taken from https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

// taken from https://opensource.com/article/18/6/copying-files-go
func copyFile(r io.Reader, size int, w io.Writer) error {
	fmt.Println("got here")
	buf := make([]byte, 2048)
	for {
		n, err := r.Read(buf)
		fmt.Println("read n, err:", n, err)
		if nil != err {
			if io.EOF != err {
				return err
			}
		}
		if n == 0 {
			fmt.Println("n break")
			break
		}

		if size < n {
			n = size
		}
		size -= n
		if _, err := w.Write(buf[:n]); err != nil {
			fmt.Println("w break")
			return err
		}
		if io.EOF == err {
			break
		}
	}
	return nil
}
*/

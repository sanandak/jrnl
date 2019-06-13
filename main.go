// Program `jrnl` is a command line journaling tool that parses a string with
// a natural-language time, title, and entry and writes out an org-mode
// entry to an org file.
// Usage: jrnln [-f jrnl-file] [when.] [title:] [entry]
// `when` is a natural-language time (today, next wednesday, march 17th); ends with `.`
// `title` is the title for the entry; ends with `:`
// `entry` is the text of the diary entry.
// The entry is written to either environment variable JRNLFILE
// or to the file specified in the command line [-f jrnl-file]
// or to `./jrnl.org` if neither of those options is present
//
// The form of the org file is with entries grouped by day.
// A new headline is started if no headline exists for today
// * Entries for 2019/06/10
// ** An Entry...
// ** Another Later Entry The Same Day...
// ** And So On
// * Entries for 2019/06/11
//
// If there are no arguments to `jrnl`, the program will open an external editor
// `emacs` by default, but set environment variable EDITOR to use another one.
//
// If `when` is blank, then `today` is used for when
// If title is blank, then the entry has no title (warning: if the entry text
// has a colon in it, then the title is everything up to the colon!)
// Tags can be defined with `@tag` and will be appended to the headline in org-mode
// fashion as :@tag:

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/sanandak/jrnl/entry"

	"log"
	"os"
	"regexp"
	"strings"
)

var (
	defaultEditor     = "emacs"
	defaultEditorArgs = []string{"-nw", "-Q"}
	defaultOrgFile    = "./jrnl.org"
	// command line usage function
	usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s [flags] [[when.] [title:] text]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "  where 'when' is a time (today, next wed); period is required")
		fmt.Fprintln(flag.CommandLine.Output(), "  where 'title' is the entry title; colon is required")
		fmt.Fprintln(flag.CommandLine.Output(), "  where 'text' is the entry text; tags @tag allowed")
		fmt.Fprintln(flag.CommandLine.Output(), "  With no args, external editor is opened\nFlags:")
		flag.PrintDefaults()
	}
)

func main() {
	flag.Usage = usage
	var orgfile = flag.String("f", "", "jrnl org file to save entry into")
	flag.Parse()
	if len(*orgfile) == 0 {
		*orgfile = os.Getenv("JRNLFILE")
		if len(*orgfile) == 0 {
			*orgfile = defaultOrgFile
		}
	}
	args := flag.Args()
	//fmt.Println("output to", *orgfile)

	var raw string
	//fmt.Println("args", args, len(args))
	if len(args) == 0 { // open editor
		raw = useEditor()
		//fmt.Println(raw)
	} else { // get entry for command line
		raw = strings.Join(flag.Args(), " ")
	}
	//fmt.Println("raw", raw)
	ent, err := entry.NewEntry(raw)
	if err != nil {
		log.Fatal("err creating new entry, ", err)
	}
	//fmt.Printf("%+v, %+v\n", ent, err)

	// convert entry to 2nd level headline, timestamp, and paragraph
	out := ent.Print()
	writeOrgFile(*orgfile, out)
}

// writeOrgFile takes a filename and a byte slice as function args
// 1. read the file and search for a top-level headline with today's date
// 2. if not present, add a top-level headline
// 3. save the byte slice (a 2nd-level headline, timestamp and paragraph)
func writeOrgFile(outf string, out []byte) {
	// read the file
	orgContents, err := ioutil.ReadFile(outf)
	if err != nil {
		log.Println("new file...", outf)
	}

	// open it again for writing
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(outf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// search for today
	// * Entries for 2019/6/11 <-- top level headline
	todayTimeStamp := time.Now().Format("2006/01/02")
	todaySearchStr := fmt.Sprintf("Entries for %s", todayTimeStamp)
	//fmt.Println(todaySearchStr)
	entryFound, _ := regexp.Match(todaySearchStr, orgContents)
	//fmt.Printf("%v", entryFound)
	if !entryFound {
		headline := fmt.Sprintf("* Entries for %s\n", todayTimeStamp)
		f.Write([]byte(headline))
	}

	if _, err := f.Write(out); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// useEditor spawns an external editor and reads the contents when the
// user saves and exits it.
func useEditor() string {
	tmpDir := os.TempDir()
	tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "jrnl")
	if tmpFileErr != nil {
		log.Printf("Error %s while creating tempFile", tmpFileErr)
		return ""
	}

	editor := os.Getenv("EDITOR")
	var editorArgs = []string{}
	if len(editor) == 0 {
		editor = defaultEditor
		editorArgs = defaultEditorArgs
	}
	path, err := exec.LookPath(editor)
	if err != nil {
		log.Fatal("Error no path for editor: ", path, editor)
	}
	//fmt.Printf("%s is available at %s\nCalling it with file %s \n", defaultEditor, path, tmpFile.Name())

	editorArgs = append(editorArgs, tmpFile.Name())

	cmd := exec.Command(path, editorArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Fatal("External editor start failed", err)
	}
	//fmt.Printf("Waiting for command to finish.\n")
	err = cmd.Wait()
	//fmt.Printf("Command finished with error: %v\n", err)
	raw, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		log.Fatal("err reading tmp file", err)
	}
	//fmt.Println("emacs",raw)
	return string(raw)
}

# jrnl - golang-based journaling app

A go command-line program to add entries to a journal.

## Usage

`jrnl [-f jrnl-file] [[when.] [title:] journal entry]`

  - `when` is a natural language time (today, next wednesday, 3/14), terminated by a period (`.`)
  - `title` is the title for the entry, terminated by a colon (`:`)
  - `journal entry` is the text of the journal entry

  - `jrnl-file` is where the entry is saved.
    This file is in org-mode (https://orgmode.org) format, with one top-level
    headline per day, and an individual second-level headline per journal entry.

If no `jrnl-file` is specified, use the environment variable $JRNLFILE.
If $JRNLFILE is also not specified, save to `./jrnl.org`

If `jrnl` is invoked without arguments (except, perhaps `-f`), then open an
external editor (default `emacs`, or from environment variable $EDITOR).

The title may start with TODO or DONE (special keywords in org).  The title and text
may include tags thus: `@tag1 and @tag2`, which will be added to the headline as 
`:tag1:tag2:

Use a period alone for when to indicate now/today.  This has the added benefit that subsequent periods are part of the title/text and not accidentally part of the `when`

`Title: First sentence.  And the second`
will be intrepred wrongly as `when`: "Title: First Sentence" and `text` : "And the second."

But `. Title: First sentence.  And the second`
Will be interpreted as `when`: "now" and `title`:"Title" etc

## Build

  `go get github.com/sanandak/jrnl`  
  `go build`  
  `go install`  

# See also

I wrote this mainly to learn go - so this is a minimal program.  Please see `jrnl.sh` (https://jrnl.sh) for
a more-complete implementation, with filtering, search, etc. 

This program doesn't have any of that because I use emacs to search the org file.

Time is parsed with WolframAlpha.  
  - go to https://developer.wolframalpha.com and log in
  - create a new app and get the app id
  - create an environment variable WOLFRAMAPPID with that string


# TODO
Includes but doesn't use fantastic library `github.com/olebedev/when` to parse the times.  Maybe a command line flag?

I want to learn about channels in go - maybe use channels for calling WolframAlpha?

Use a config file instead of env vars?
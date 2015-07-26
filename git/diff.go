package git

/*
The delta compression algorithm that is used in git was originally based on [xdelta](http://xdelta.org/) and [LibXDiff](http://www.xmailserver.org/xdiff-lib.html)
 but was further simplified for the git use case (see the [“diff’ing files”](http://git.661346.n2.nabble.com/diff-ing-files-td6446460.html) thread on the git mailinglist).
The git delta encoding algorithm is a copy/insert based algorithm (this is apparent in patch-delta.c). The delta representation contains a delta header and a series of opcodes for either copy or insert instructions.
*/

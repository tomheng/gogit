# Gogit

Gogit is a go git in golang from bottom-up. It is just a golang programing ***excercise***.


# Quick Start

it is very easy to have a try, just use it as git, replacing git to gogit in your cmd.

~~~bash
go get github.com/tomheng/gogit
gogit ls git://github.com/tomheng/gogit
~~~

# Commands

1. Supported

	* gogit ls-remote (partial function)

2. on the way

	* git cat-file
	* git ls-tree
	* git clone
	* git gc
	* git daemon
	* git hash-object
	* git write-tree
	* git checkout
	* git branch
	* git show-branch
	* git unpack-objects
	* git reset
	* git add
	* git commit
	* git pull
	* git push
	* git symbolic-ref
	* git update-ref
	* git commit-tree
	* git unpack-objects
	* git rev-list
	* git rev-parse

#Reference

* [git technical](https://github.com/git/git/tree/master/Documentation/technical)
* [unpacking git packfiles](https://codewords.recurse.com/issues/three/unpacking-git-packfiles/)
* [git clone in haskell from the bottom up](http://stefan.saasen.me/articles/git-clone-in-haskell-from-the-bottom-up)
* [git source code](https://github.com/git/git)
* [File System Support for Delta Compression](http://mail.xmailserver.net/xdfs.pdf)
* [remyoudompheng gigot](https://github.com/remyoudompheng/gigot)
* [ChimeraCoder Gitgo](https://github.com/ChimeraCoder/gitgo/)
* [Decentralized, peer-to-peer Git repositories aka "Git meets Bitcoin"](https://github.com/gitchain/gitchain)

#License
Gogit is primarily distributed under the terms of both the MIT license.
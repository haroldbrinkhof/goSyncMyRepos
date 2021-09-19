# goSyncMyRepos

## Description
small tool that syncs multiple related git repositories to the current commit in the current directory
if this directory is defined as being in a group with them in the ~/.goSyncMyRepos.config file.

### Config file format
The format of this file is as follows:

    [groupname]
    /dir/to/repo1
    /dir/to/repo2
    /dir/to/repo3
    
    [new group]
    /dir/to/repoA
    /dir/to/repoB

    [another group]
    /dir/to/repoC
    /dir/to/repoB

#### Effect 
using this app in /dir/to/repo2 would sync the repositories at /dir/to/repo1 and /dir/to/repo2 (the ones in 'groupname')
to the timestamp of the current commit in repo2 (first commit with timestamp <= timestamp current commit).

using this app in /dir/to/repoB would do the same but for /dir/to/repoA and /dir/to/repoC since both of these are in a group with repoB


## Purpose
The purpose of this application is to more easily do full-app compiles where your app is a multi-repo (so spread over multiple repositories) without easy way to synchronize between them. General use-case I made it for was to sync the companion repos when doing git bisect.

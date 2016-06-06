[![Circle CI](https://circleci.com/gh/asarturas/paire/tree/master.svg?style=svg)](https://circleci.com/gh/asarturas/paire/tree/master)
# paire - simple release manager for projects, hosted at github

Features:
 - [x] push pre-release packages for current commit
 - [x] push release packages for current tag
 - [x] pull packages from specified tag/commit
 - [ ] parallel upload of multiple packages for quicker release
 - [ ] pull of specific asset only
 - [ ] option to disable pushing pre-release packages
 - [ ] warn customer then pushing dirty branch

## What is paire for?

Paire is built with a goal to simplify built assets of projects, hosted on github.

Example use case is dodin project, which is binary.
Asking customers to build their own paire executables would not be user friendly or cost effective.
Creating releases on github and manually uploading pre-built executables is mundane, repetitive and error prone.
The optimal solution is to add paire to deployment step, so that built packages would be tagged and uploaded to github automatically on each release. 

## How to install paire?

The latest version is 0.1.0.

Paire consists of two commands:
 - `paire-push` to upload packages (pre-built 64bin binaries for [Linux](l), [Mac](m),  [Windows](w)).
 - `paire-pull` to download packages for github (pre-built 64bin binaries for [Linux](l), [Mac](m),  [Windows](w)).

## How to use paire-push?

There are two ways you could use paire:

1. You could run it on local environment to upload package, built on your dev environment:
   
   ```
   curl https://github.com/asarturas/paire/releases/download/0.2.2/paire-push -o paire-push
   chmod +x paire-push
   ./paire-push -package your-built-package.zip
   ```
   
2. You could run it on deployment step in your continuous integration server. Example for Circle CI:
    
   ```
   deployment:
     github:
       branch: master
       commands:
         - curl https://github.com/asarturas/paire/releases/download/0.2.2/paire-push -o paire-push
         - chmod +x paire-push
         - ./paire-push -package your-built-package.zip
   ```

Where `your-built-package.zip` is your built package
If you have more than one package to upload, then just add another -package attribute.
Paire will recognise if current commit is tagged.
It will push pre-release if it's not (the version string will match the commit sha).

## How to use paire-pull?

Paire pull is useful when you need to download all built packages for certain release.
It is particularly useful if repository is private and there is no direct access via curl.
For example, these commands would download paire executable to current directory:
```
curl https://github.com/asarturas/paire/releases/download/0.2.2/paire-pull -o paire-pull
chmod +x paire-pull
./paire-pull -destination . -version 0.2.2
```

## Why use paire?

To simplify and streamline release of built packages and save time for yourself and others.

## Disclaimer

The software is provided as is. I use it and am happy to help you setup or to resolve any issues with it, but you are taking all the responsility for using it.

# Conget
Conget is a CLI app which is a concurrent file downloader  that download the file data by splitting into several chunks and fetch the data asynchronously.

## Install
Get the dependency with 'dep'

`$ dep ensure`

If you don't have `dep` installed, this app has only one dependency right now `github.com/urfave/cli` so you can just get it with 

`go get github.com/urfave/cli`.

And run this  

`$ go install`

If you have added go bin diretory to your $GOPATH, Just like this: `PATH="$GOPATH/bin:$PATH"` then you can just run `conget` command from anywhere on your terminal.
Or you may have to move the compiled binary file from `$GOPATH/bin` to `/usr/local/bin` to run it globally.   

**Installing or running it without Go**

Clone the repo or just get the `conget` file and run it like this `./conget` or if you want to run it as global command just move the file to your `/usr/local/bin`.


## Usage
`$ conget -u http://example.com/example.mp4`

The default concurrent number is set to 10 if no number is provided. To set the concurrent number provide the number with -c option. Ex:
`$ conget -c 12 -u  http://example.com/example.mp4`

**Note**
It downloads the file in current directory where you run the command.

### Contribution
This is just in experimental phase, so if you think anything could be improved or fixed just send the PR.  

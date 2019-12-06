# Conget
Conget is a CLI app which is a concurrent file downloader  that download the file data by splitting into several chunks and fetch the data asynchronously.

## Install & Build
Get the dependency with `go mod`

**Build with Makefile** 
To compile and build the binary file, run this command. it will build the binary in bin/ directory with the name conget.

```sh
$ make build  
```  
Or for  Mac 

```sh
$ make build-mac
```

To clean the compiled binary file 

```sh 
$ make clean 
```
It will remove all the binary file 

For running tests. 

```sh
$ make test 
```

**Installing or running it without Go**

Clone the repo or just get the `conget` file and run it like this `./conget` or if you want to run it as global command just move the file to your `/usr/local/bin`.


## Usage
```sh
$ conget -u http://example.com/example.mp4
```

The default concurrent number is set to 5 if -c flag is not set. To set the concurrent number provide the number with -c flag. Ex:
```sh 
$ conget -c 12 -u  http://example.com/example.mp4
```

**Note**
It downloads the file in current directory where you run the command.

### Contribution
It's still in experimental phase, so if you think anything could be improved or fixed just send the PR.  

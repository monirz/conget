# Conget
Conget is a CLI application, which downloads file data concurrently by splitting the data into several chunks and fetch those asynchronously.

## Install & Build

Conget uses [Cobra](https://github.com/spf13/cobra) for handling command line interaction.     
Get the dependency with `go mod`

### Build with Makefile 

To compile and build the binary file, run this command. it will build the binary in bin/ directory with the name conget.

```sh
$ make build  
```  
Or for  Mac 

```sh
$ make build-mac
```

To remove the compiled binary file. 

```sh 
$ make clean 
```
It will remove all the binary file. 

For running tests. 

```sh
$ make test 
```

###Installing or running it without Go

Binary file of latest build is pushed in the bin/ directory, download the file according to your OS and run the command `./conget`    


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

# Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Make changes and add them (`git add .`)
4. Commit your changes (`git commit -m 'Added some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new pull request

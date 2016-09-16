# Installation

* Install Golang.
 * Ubuntu users can do `sudo apt install golang`
* Set the `GOPATH` variable.
 * Open up `~/.bashrc` and add the following lines (you can set your own location, this is what I used)

```
export GOPATH=$HOME/.go
export PATH=$PATH:$GOPATH/bin
```

* Create a folder called `src` inside of `$GOPATH`.
 * `mkdir $GOPATH/src`
* Clone this repo in that folder

```
cd $GOPATH/src
git clone https://github.com/tkshnwesper/niec.git
```

* Install the various Golang dependencies

```
cd $GOPATH/src/niec
go get
```

* Set up a MariaDB/MySql database (I'm assuming you have one of those installed already)
 * Execute the queries from the `query` file
 * Ensure that the user who runs the server has access to that database
* Set up a `config.json` file in at the root level of the project `$GOPATH/src/niec/config.json`
 * It should consist of the following

```
{
    "DB": {
        "Name": "name_of_the_database",
        "User": "database_username",
        "Password": "database_password"
    }
}
```

* Set up a `static` folder in the project root and install the frontend libraries

```
mkdir $GOPATH/src/niec/static
cd $GOPATH/src/niec/static
```

 * I hope you have `bower`. No problem even if you don't. Do `npm i -g bower`. You might need to prepend `sudo`.
 * `bower install` the following inside the `static` folder
  1. `jquery`
  2. `bootstrap`
  3. Also `highlight.js`, but you can skip that for now.

* We are done! `cd` to the project root directory and do `go clean; go install; niec` to start the server.
# Description
The goal of tsql-cli is to provide a simple interface to mysql databases without using mysql-cli
# Why
Locked in terminal/REPL with mysql-cli
Utilize CLI history for quick access to common queries ( zsh history auto fill )
mysqlcli history sucks
Provide a structure to save and load queries from directories
pretty print to terminal, because.

##
tsql-cli -g gateway 'select duck from duck;'

tsql-cli -g gateway -p 0 'select duck from duck;'

export TSQL_DEFAULT_GROUP="gateway"
export TSQL_DEFAULT_SERVER_POSITION=0
tsql-cli 'select duck from duck;'
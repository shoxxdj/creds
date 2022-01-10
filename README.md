# creds

A go binary to store the creds of a pentest 

## Setup
```
go get github.com/shoxxdj/creds
```
## Usage
```
Creds : a binary to store creds for attackers. v:0.1
	-d: The Creds ID to delete from the database
	-dl: Database location
	-full: Get full details (Location, creds and id)
	-l: Login to add in the database
	-p: Password to add in the database
	-reset: Reset configuration to default
	-save: Save configuration (need dbLocation to be defined to be efficient)
```

# Container Graph project
Hamza Motiwalla

Shivakant Mishra

University of Colorado Boulder


## Pre-requisites:
1. Docker
2. Golang


### To launch the controller
`$ go run controller.go dockerfunc.go helper.go`

The controller will now listen for requests on localhost:8000


### To test the deployment (run in a different terminal)
`$ cd tests/`

`$ go run creategraph.go`

It take a minute or so to get all vertices and edges up and running..

`$ docker logs vertex3`

### To add vertex4
`$ cd tests/`

`$ go run addvertex4.go`

`$ docker logs vertex3`

### To remove vertex2
`$ cd tests/`

`$ go run remvertex2.go`

`$ docker logs vertex3`

### To clean up
Kill the controller.go process

`$ bash clean.sh`


## Analysis and Evaluation 
Coming soon...

Also checkout hamza15comp/vertex

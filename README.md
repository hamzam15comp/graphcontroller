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
`$ go run tests/creategraph.go`

`$ docker logs vertex3`

### To add vertex4
`$ go run tests/addvertex4.go`

`$ docker logs vertex3`

### To remove vertex2
`$ go run tests/remvertex2.go`

`$ docker logs vertex3`


## Analysis and Evaluation 
Coming soon...


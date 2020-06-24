package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	"github.com/hamzam15comp/vertex"
)

type Graph struct {
	VertexArr []Vertex `json:"vertexarr"`
	EdgeArr   []Edge   `json:"edgearr"`
	GraphId   int      `json:"graphid"`
	GraphName string   `json:"graphname"`
}

type Edge struct {
	EdgeName string `json:"edgename"`
	EdgeId   int    `json:"edgeid"`
}

type Vertex struct {
	VertexName string `json:"vertexname"`
	VertexId   int    `json:"vertexid"`
	Pub        []Edge `json:"pub"`
	Sub        []Edge `json:"sub"`
}

var graph = Graph{
	GraphName: "graphy",
	GraphId: 42,
}

func loadEdge(w http.ResponseWriter, r *http.Request) (Edge){
	var ed Edge
	vererr := getFile(r, "edge")
	if vererr != nil {
		http.Error(
			w,
			"failed to get vertex json file",
			http.StatusBadRequest,
		)
		return Edge{}
	}
	temp, rerr := ioutil.ReadFile("edge.json")
	err := json.Unmarshal((temp), &ed)
	if err != nil || rerr != nil {
		http.Error(
			w,
			"failed to unmarshall vertex json file",
			http.StatusBadRequest,
		)
		return Edge{}
	}
	graph.EdgeArr = append(graph.EdgeArr, ed)
	return ed
}

func loadVertex(w http.ResponseWriter, r *http.Request) (Vertex){
	var vert Vertex
	vererr := getFile(r, "vertex")
	if vererr != nil {
		http.Error(
			w,
			"failed to get vertex json file",
			http.StatusBadRequest,
		)
		return Vertex{}
	}
	temp, rerr := ioutil.ReadFile("vertex.json")
	err := json.Unmarshal((temp), &vert)
	if err != nil || rerr != nil {
		http.Error(
			w,
			"failed to unmarshall vertex json file",
			http.StatusBadRequest,
		)
		return Vertex{}
	}
	graph.VertexArr = append(graph.VertexArr, vert)
	return vert
}

func loadGraph(w http.ResponseWriter, r *http.Request) {
	grapherr := getFile(r, "graph")
	if grapherr != nil {
		http.Error(
			w,
			"failed to get graph json file",
			http.StatusBadRequest,
		)
		return
	}
	temp, _ := ioutil.ReadFile("graph.json")
	err := json.Unmarshal((temp), &graph)
	if err != nil {
		http.Error(
			w,
			"failed to unmarshall graph json file",
			http.StatusBadRequest,
		)
		return
	}
	logger.Println("Graph loaded successfully")
	graphJson, _ := json.MarshalIndent(graph, "", "  ")
	logger.Println(string(graphJson))
}

func getVertexObj(vname string) (v Vertex) {
	for _, v := range graph.VertexArr {
		if v.VertexName == vname {
			return v
		}
	}
	return Vertex{}

}
func getEdgeObj(ename string) (v Edge) {
	for _, e := range graph.EdgeArr {
		if e.EdgeName == ename {
			return e
		}
	}
	return Edge{}
}

func addConnections(v Vertex) (error) {
	for _, p := range v.Pub {
		vertex.SendToVagent(
			vertex.ControlMsg{
				Edge: p.EdgeId,
				Vertexno: v.VertexId,
				Vertextype: "pub",
				Cmd: "add",
				Msgid: 1,
			},
			"localhost",
			7000+v.VertexId,
		)
	}
	for _, s := range v.Sub {
		vertex.SendToVagent(
			vertex.ControlMsg{
				Edge: s.EdgeId,
				Vertexno: v.VertexId,
				Vertextype: "sub",
				Cmd: "add",
				Msgid: 1,
			},
			"localhost",
			7000+v.VertexId,
		)
	}
	return nil
}

func removeConnections(v Vertex) (error) {
	for _, p := range v.Pub {
		vertex.SendToVagent(
			vertex.ControlMsg{
				Edge: p.EdgeId,
				Vertexno: v.VertexId,
				Vertextype: "pub",
				Cmd: "rem",
				Msgid: 1,
			},
			"localhost",
			7000+v.VertexId,
		)
	}
	for _, s := range v.Sub {
		vertex.SendToVagent(
			vertex.ControlMsg{
				Edge: s.EdgeId,
				Vertexno: v.VertexId,
				Vertextype: "sub",
				Cmd: "rem",
				Msgid: 1,
			},
			"localhost",
			7000+v.VertexId,
		)
	}
	return nil
}

func createGraph(w http.ResponseWriter, r *http.Request) {
	parseMultiPart(w, r)
	loadGraph(w, r)
	for _, e := range graph.EdgeArr {
		getFile(r, e.EdgeName)
		BuildImage(e.EdgeName)
		CreateContainer(e.EdgeName)
	}
	time.Sleep(20*time.Second)
	for _, v := range graph.VertexArr {
		getFile(r, v.VertexName)
		BuildImage(v.VertexName)
		CreateContainer(v.VertexName)
		time.Sleep(10*time.Second)
		addConnections(v)
	}
}


func addVertex(w http.ResponseWriter, r *http.Request) {
	parseMultiPart(w, r)
	v := loadVertex(w, r)
	getFile(r, v.VertexName)
	BuildImage(v.VertexName)
	CreateContainer(v.VertexName)
	time.Sleep(10*time.Second)
	addConnections(v)
}


func removeVertex(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	vname := params["name"]
	v := getVertexObj(vname)
	if v.VertexId == 0 {
		http.Error(
			w,
			"Failed to find vertex",
			http.StatusBadRequest,
		)
		return
	}
	removeConnections(v)
	err := removeContainer(vname)
	if err != nil {
		http.Error(
			w,
			"Failed to remove vertex",
			http.StatusBadRequest,
		)
		return
	}
}

func addConn (w http.ResponseWriter, r *http.Request) {}
func remConn (w http.ResponseWriter, r *http.Request) {}

func addEdge(w http.ResponseWriter, r *http.Request) {
	parseMultiPart(w, r)
	e := loadEdge(w, r)
	getFile(r, e.EdgeName)
	BuildImage(e.EdgeName)
	CreateContainer(e.EdgeName)
}

func removeEdge(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	ename := params["name"]
	//remove all associated connections first
	removeContainer(ename)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Container Graph home!")
}

func main() {
	logInit()
	pwd, _ = os.Getwd()
	pwd = pwd + "/"
	err := CreateNetwork("graph")
	if err != nil {
		logger.Println("Network error", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods(http.MethodGet)
	router.HandleFunc("/createGraph", createGraph).Methods(http.MethodPost)
	//router.HandleFunc("/getGraph", getGraph).Methods(http.MethodPost)
	router.HandleFunc("/addVertex", addVertex).Methods(http.MethodPost)
	router.HandleFunc("/addConn/{edge}/{vertex}", addConn).Methods(http.MethodPost)
	router.HandleFunc("/remConn/{edge}/{vertex}", remConn).Methods(http.MethodPost)
	router.HandleFunc("/removeVertex/{name}", removeVertex).Methods(http.MethodGet)
	router.HandleFunc("/addEdge", addEdge).Methods(http.MethodPost)
	router.HandleFunc("/removeEdge/{name}", removeEdge).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8000", router))

}

/*

e1 := Edge{EdgeName: "edge1", EdgeId: 1}
v1 := Vertex{VertexName: "vertex1", VertexId: 1, Pub: []Edge{e1}, Sub: []Edge{}}
v2 := Vertex{VertexName: "vertex2", VertexId: 2, Pub: []Edge{}, Sub: []Edge{e1}}
graf := Graph{
		VertexArr:[]Vertex{v1, v2},
		EdgeArr: []Edge{e1},
		GraphId: 1,
		GraphName:"graph1"
	}

*/

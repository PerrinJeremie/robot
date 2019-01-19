package main

import (
	"encoding/gob"
	"os"
	"polygon"
)

type Env struct {
	VertexBegin    polygon.Vertex
	VertexEnd      polygon.Vertex
	Space          polygon.Space
	SpaceConf      polygon.Space
	SpaceAcc       polygon.Space
	RobSpace       polygon.Space
	SpaceInit      polygon.Space
	AllEdge        []polygon.Edge
	AllVert        []polygon.Vertex
	RobVert        []polygon.Vertex
	RobEdge        []polygon.Edge
	PolRobot       polygon.Polygon
	PolRobotVect   []polygon.Vertex
	ShowBorder     bool
	ShowCells      bool
	ShowPath       bool
	ShowVisibility bool
	Mink           bool
	Setr           bool
	Seto           bool
	ListOfEdges    []polygon.Edge
	RobDec         []polygon.Polygon
	AccPol         []polygon.Polygon
	Cells          []polygon.Polygon
	G              polygon.Graph
	Vg             polygon.VGraph
	Lpath          []polygon.Vertex
	Turns          int
}

func SaveCurrentEnv(path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rdec := []polygon.Polygon{}
	for _, polptr := range *robDec {
		rdec = append(rdec, *polptr)
	}

	cel := []polygon.Polygon{}
	for _, polptr := range cells {
		cel = append(cel, *polptr)
	}

	env := Env{
		VertexBegin:    vertexBegin,
		VertexEnd:      vertexEnd,
		Space:          space,
		SpaceConf:      spaceConf,
		SpaceAcc:       spaceAcc,
		RobSpace:       robSpace,
		SpaceInit:      spaceInit,
		AllEdge:        allEdge,
		AllVert:        allVert,
		RobVert:        robVert,
		RobEdge:        robEdge,
		PolRobot:       polRobot,
		PolRobotVect:   polRobotVect,
		ShowBorder:     showBorder,
		ShowCells:      showCells,
		ShowPath:       showPath,
		ShowVisibility: showVisibility,
		Mink:           mink,
		Setr:           setr,
		Seto:           seto,
		ListOfEdges:    *listOfEdges,
		RobDec:         rdec,
		AccPol:         *accPol,
		Cells:          cel,
		G:              g,
		Vg:             vg,
		Lpath:          lpath,
		Turns:          turns,
	}

	enc := gob.NewEncoder(f)
	err = enc.Encode(env)
	if err != nil {
		panic(err)
	}
}

func LoadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rdec := []polygon.Polygon{}
	for _, polptr := range *robDec {
		rdec = append(rdec, *polptr)
	}

	cel := []polygon.Polygon{}
	for _, polptr := range cells {
		cel = append(cel, *polptr)
	}

	env := Env{}
	dec := gob.NewDecoder(f)
	err = dec.Decode(&env)
	if err != nil {
		panic(err)
	}

	vertexBegin = env.VertexBegin
	vertexEnd = env.VertexEnd
	space = env.Space
	spaceConf = env.SpaceConf
	spaceAcc = env.SpaceAcc
	robSpace = env.RobSpace
	spaceInit = env.SpaceInit
	allEdge = env.AllEdge
	allVert = env.AllVert
	robVert = env.RobVert
	robEdge = env.RobEdge
	polRobot = env.PolRobot
	polRobotVect = env.PolRobotVect
	showBorder = env.ShowBorder
	showCells = env.ShowCells
	showPath = env.ShowPath
	showVisibility = env.ShowVisibility
	mink = env.Mink
	setr = env.Setr
	seto = env.Seto
	listOfEdges = &(env.ListOfEdges)
	accPol = &(env.AccPol)
	g = env.G
	vg = env.Vg
	lpath = env.Lpath
	turns = env.Turns

	robDec = &([]*polygon.Polygon{})
	for _, pol := range env.RobDec {
		*robDec = append(*robDec, &pol)
	}
	cells = []*polygon.Polygon{}
	for _, pol := range env.Cells {
		cells = append(cells, &pol)
	}

}

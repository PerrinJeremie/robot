package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/veandco/go-sdl2/gfx"
import "github.com/veandco/go-sdl2/ttf"
import "polygon"
import "time"

//import "math"
import "fmt"

var (
	MAIN_WINDOW_ID uint32
	INSTWINDOW     *sdl.Window
	PLWINDOW       *sdl.Window
	INSTRENDERER   *sdl.Renderer
	PLRENDERER     *sdl.Renderer
	FONT           *ttf.Font
	err            error
	WHITE          sdl.Color
	spaceInit      polygon.Space
	space          polygon.Space
	spaceConf      polygon.Space
	robSpace       polygon.Space
	spaceAcc       polygon.Space
	listOfVert     []polygon.Vertex
	polRobot       polygon.Polygon
	polRobotVect   []polygon.Vertex
	allVert        []polygon.Vertex
	allEdge        []polygon.Edge
	robVert        []polygon.Vertex
	robEdge        []polygon.Edge
	quit           bool
	setr           bool
	seto           bool
	mink           bool
	vertexBegin    polygon.Vertex
	vertexEnd      polygon.Vertex
	frame_polygon1 = polygon.Polygon{map[uint32]polygon.Vertex{0: polygon.Vertex{0, 0, 0}, 1: polygon.Vertex{1, -5, 1400}, 2: polygon.Vertex{2, 2005, 1400}, 3: polygon.Vertex{3, 2000, 0}}, map[uint32]([]uint32){0: []uint32{1, 3}, 1: []uint32{0, 2}, 2: []uint32{1, 3}, 3: []uint32{2, 0}}}
	frame_polygon2 = polygon.Polygon{map[uint32]polygon.Vertex{0: polygon.Vertex{0, -10, -10}, 1: polygon.Vertex{1, -15, 1410}, 2: polygon.Vertex{2, 2015, 1410}, 3: polygon.Vertex{3, 2010, -10}}, map[uint32]([]uint32){0: []uint32{1, 3}, 1: []uint32{0, 2}, 2: []uint32{1, 3}, 3: []uint32{2, 0}}}
)

const (
	INSTRUCTION_MENU = "\n Menu \n\nR - Set Robot \nO - Set Obstacles\nM - Draw Minkowski\nC - Clear Obstacles\nQ - Quit \n"

	INSTRUCTION_ROBOT     = "\n Robot \n\nDraw the robot\n\nD - Finish and draw last edge\nQ - Quit"
	INSTRUCTION_PLACE     = "\n Robot \n\nPlace start\nPlace end\nQ - Quit"
	INSTRUCTION_OBSTACLES = "\n Obstacles \n\nDraw the obstacle\nD - Finish obstacle\nM - Goto Menu\nQ - Quit\n"
	FRAME_RATE            = 24
)

func Init() {

	if err = ttf.Init(); err != nil {
		panic(err)
	}

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	INSTWINDOW, err = sdl.CreateWindow("Instructions", 0,
		0, 400, 800, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	INSTRENDERER, err = sdl.CreateRenderer(INSTWINDOW, 0, 0)
	if err != nil {
		panic(err)
	}

	PLWINDOW, err = sdl.CreateWindow("Playground", 450,
		0, 2000, 1400, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	MAIN_WINDOW_ID, err = PLWINDOW.GetID()
	if err != nil {
		panic(err)
	}

	PLRENDERER, err = sdl.CreateRenderer(PLWINDOW, 0, 0)
	if err != nil {
		panic(err)
	}

	FONT, err = ttf.OpenFont("courier.ttf", 30)
	if err != nil {
		panic(err)
		fmt.Printf("Error loading font\n")
	}

	WHITE = sdl.Color{255, 255, 255, 0}

	PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
	PLRENDERER.FillRect(nil)
	PLRENDERER.Present()
	PLRENDERER.SetDrawColor(0, 0, 0, 0)

	setr = false
	mink = false
	quit = true
	seto = false

	space = polygon.NewEmptySpace()
	spaceInit = polygon.NewEmptySpace()
	spaceConf = polygon.NewEmptySpace()
	robSpace = polygon.NewEmptySpace()
	spaceAcc = polygon.NewEmptySpace()

	polRobot = polygon.Polygon{}
	polRobotVect = []polygon.Vertex{}
	vertexBegin = polygon.Vertex{}
	vertexEnd = polygon.Vertex{}

	PrintInstruction(INSTRUCTION_MENU)

}

func Quit() {
	ttf.Quit()
	sdl.Quit()
	INSTWINDOW.Destroy()
	PLWINDOW.Destroy()
}

func DrawPolygon(pol *polygon.Polygon, c1, c2, c3 uint8, width int32) {
	for k, v := range pol.Edges {
		drawn := gfx.ThickLineRGBA(PLRENDERER, int32(pol.Vertices[k].X), int32(pol.Vertices[k].Y),
			int32(pol.Vertices[v[0]].X), int32(pol.Vertices[v[0]].Y), width, c1, c2, c3, 0xff)
		if !drawn {
			fmt.Printf("Not Drawn")
		}
		drawn = gfx.ThickLineRGBA(PLRENDERER, int32(pol.Vertices[k].X), int32(pol.Vertices[k].Y),
			int32(pol.Vertices[v[1]].X), int32(pol.Vertices[v[1]].Y), width, c1, c2, c3, 0xff)
		if !drawn {
			fmt.Printf("Not Drawn")
		}
	}
}

func DrawSpace(space *polygon.Space, c1, c2, c3 uint8, width int32) {
	for _, p := range space.Polygons {
		DrawPolygon(&p, c1, c2, c3, width)
	}
}

func DrawListVectFromVert(listOfVects *[]polygon.Vertex, x, y float64, c1, c2, c3 uint8, width int32) {
	xtmp := x
	ytmp := y
	for _, v := range *listOfVects {
		//PLRENDERER.DrawLine(int32(xtmp), int32(ytmp), int32(xtmp+v.X), int32(ytmp+v.Y))
		drawn := gfx.ThickLineRGBA(PLRENDERER, int32(xtmp), int32(ytmp), int32(xtmp+v.X), int32(ytmp+v.Y), width, c1, c2, c3, 0xff)
		if !drawn {
			fmt.Printf("Not Drawn")
		}
		xtmp = xtmp + v.X
		ytmp = ytmp + v.Y
	}
}

func DrawListVert(listOfVerts *[]polygon.Vertex, c1, c2, c3 uint8, width int32) {
	for i := 0; i < len(*listOfVerts)-1; i++ {
		x1 := (*listOfVerts)[i].X
		y1 := (*listOfVerts)[i].Y
		x2 := (*listOfVerts)[i+1].X
		y2 := (*listOfVerts)[i+1].Y
		drawn := gfx.ThickLineRGBA(PLRENDERER, int32(x1), int32(y1), int32(x2), int32(y2), width, c1, c2, c3, 0xff)
		if !drawn {
			fmt.Printf("Not Drawn")
		}
	}
}

func GetRobot() {
	running := true
	listOfClicks := []polygon.Vertex{}
	robEdge = []polygon.Edge{}
	robVert = []polygon.Vertex{}
	last := time.Now().UnixNano()
	estimate := time.Now().UnixNano() - last
	i := 0
	PrintInstruction(INSTRUCTION_ROBOT)
	for running && quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				quit = false
				break
			case *sdl.MouseMotionEvent:
				if ev.WindowID == MAIN_WINDOW_ID {
					if i > 0 {
						PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0xff)
						PLRENDERER.FillRect(nil)
						_ = gfx.ThickLineRGBA(PLRENDERER, int32(listOfClicks[i-1].X),
							int32(listOfClicks[i-1].Y), ev.X, ev.Y, 5, 0, 0xff, 0, 0xff)

					}
					vertexBegin = polygon.Vertex{0, float64(ev.X), float64(ev.Y)}
				}
				break
			case *sdl.MouseButtonEvent:
				if ev.State == sdl.PRESSED && ev.WindowID == MAIN_WINDOW_ID {
					vclick := polygon.Vertex{uint32(i), float64(ev.X), float64(ev.Y)}
					if !vclick.SameX(&robVert) {
						fmt.Printf("%v\n\n", []int32{ev.X, ev.Y})
						if i > 0 {
							eclick := polygon.Edge{listOfClicks[i-1], vclick}
							if !eclick.IntersectsOne(&robEdge) {
								listOfClicks = append(listOfClicks, vclick)
								drawn := gfx.ThickLineRGBA(PLRENDERER, int32(listOfClicks[i-1].X),
									int32(listOfClicks[i-1].Y), ev.X, ev.Y, 5, 0, 0xff, 0, 0xff)
								if !drawn {
									fmt.Printf("Not Drawn")
								}
								robEdge = append(robEdge, eclick)
								robVert = append(robVert, vclick)
								i++
							}
						} else {
							listOfClicks = append(listOfClicks, vclick)
							robVert = append(robVert, vclick)
							i++
						}
						vertexBegin = polygon.Vertex{0, float64(ev.X), float64(ev.Y)}
					}
				}
				break
			case *sdl.KeyboardEvent:
				if ev.Type == sdl.KEYDOWN {
					switch ev.Keysym.Sym {
					case sdl.GetKeyFromName("Q"):
						quit = false
						break
					case sdl.GetKeyFromName("D"):
						if i < 3 {
							break
						}
						e := polygon.Edge{polygon.Vertex{0, listOfClicks[i-1].X, listOfClicks[i-1].Y},
							polygon.Vertex{0, listOfClicks[0].X, listOfClicks[0].Y}}
						if !e.IntersectsOne(&robEdge) {
							drawn := gfx.ThickLineRGBA(PLRENDERER, int32(listOfClicks[i-1].X),
								int32(listOfClicks[i-1].Y), int32(listOfClicks[0].X),
								int32(listOfClicks[0].Y), 5, 0, 0xff, 0, 0xff)
							if !drawn {
								fmt.Printf("Not Drawn")
							}
							polRobot = *(polygon.PolygonOfList(&listOfClicks))
							fmt.Printf("%v\n\n", polRobot)

							polRobotVect = polygon.GetFirstVect(listOfClicks)

							running = false
						}
						break
					}
				}
			}
			if time.Now().UnixNano()-last > estimate && running {
				tmp := time.Now().UnixNano()
				DrawSpace(&space, 0x66, 0x66, 0x66, 5)
				DrawListVert(&listOfClicks, 0, 0xff, 0, 5)
				PLRENDERER.Present()
				estimate = time.Now().UnixNano() - tmp
				last = time.Now().UnixNano()
			}
		}

	}
}

func BeginPolygon() *[]polygon.Vertex {
	l := []polygon.Vertex{vertexBegin}
	for _, v := range polRobotVect[:len(polRobotVect)-1] {
		l = append(l, polygon.Vertex{0, l[len(l)-1].X + v.X, l[len(l)-1].Y + v.Y})
	}
	return &l
}

func DecomposeRobot() *[](*polygon.Polygon) {
	l := BeginPolygon()
	polygon.AddListToSpaceAsPoly(l, &robSpace)
	fmt.Printf("\n%v\n", robSpace)
	res := polygon.Decompose(&robSpace)
	return &res
}

func SetLines() {
	n := len(space.Polygons)
	fmt.Printf("\nNumber Of Polygons : %v\n", n)
	l := polygon.Decompose(&space)
	fmt.Printf("\nNumber Of Polygons: %v \n", len(l))
	space.Polygons = space.Polygons[n:]
	fmt.Printf("\nNumber Of Polygons: %v\n", len(space.Polygons))
	for _, v := range l {
		polygon.AddPolygonToSpace(v, &space)
	}
	DrawSpace(&space, 0xff, 0, 0, 5)
	PLRENDERER.Present()
	fmt.Printf("\nOk\n")
}

func isOkay(v polygon.Vertex, polRobotVect *[]polygon.Vertex) bool {
	if v.SameX(&allVert) {
		return false
	}
	for _, ve := range *polRobotVect {
		v1 := polygon.Vertex{0, v.X + ve.X, v.Y + ve.Y}
		e := polygon.Edge{v, v1}
		if e.IntersectsOne(&allEdge) || v1.SameX(&allVert) {
			return false
		}
		v = v1
	}
	return true
}

func RobotEdge(place int) {
	if place == 0 {
		robEdge = []polygon.Edge{}
		tmpv := vertexBegin
		for _, v := range polRobotVect {
			robEdge = append(robEdge, polygon.Edge{tmpv, polygon.Vertex{0, tmpv.X + v.X, tmpv.Y + v.Y}})
			tmpv = polygon.Vertex{0, tmpv.X + v.X, tmpv.Y + v.Y}
		}
	} else {
		tmpv := vertexEnd
		for _, v := range polRobotVect {
			robEdge = append(robEdge, polygon.Edge{tmpv, polygon.Vertex{0, tmpv.X + v.X, tmpv.Y + v.Y}})
			tmpv = polygon.Vertex{0, tmpv.X + v.X, tmpv.Y + v.Y}
		}
	}
}

func SetBeginEnd() {
	running := true
	place := 0
	last := time.Now().UnixNano()
	estimate := time.Now().UnixNano() - last
	PrintInstruction(INSTRUCTION_PLACE)
	for running && quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				quit = false
				break
			case *sdl.MouseButtonEvent:
				if ev.State == sdl.PRESSED && ev.WindowID == MAIN_WINDOW_ID {
					vclick := polygon.Vertex{0, float64(ev.X), float64(ev.Y)}
					if isOkay(vclick, &polRobotVect) {
						if place == 0 {
							vertexBegin = polygon.Vertex{0, float64(ev.X), float64(ev.Y)}
							RobotEdge(place)
							place = 1
						} else if place == 1 {
							vertexEnd = polygon.Vertex{1, float64(ev.X), float64(ev.Y)}
							RobotEdge(place)
							fmt.Printf("%v\n\n", vertexEnd)
							running = false
						}
					}

				}
				break
			case *sdl.MouseMotionEvent:
				if ev.WindowID == MAIN_WINDOW_ID {
					if time.Now().UnixNano()-last > estimate {
						tmp := time.Now().UnixNano()
						PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
						PLRENDERER.FillRect(nil)
						DrawSpace(&space, 0x66, 0x66, 0x66, 5)
						if place == 0 {
							DrawListVectFromVert(&polRobotVect, float64(ev.X), float64(ev.Y), 0, 0xff, 0, 5)
						} else if place == 1 {
							DrawListVectFromVert(&polRobotVect, vertexBegin.X, vertexBegin.Y, 0, 0xff, 0, 5)
							DrawListVectFromVert(&polRobotVect, float64(ev.X), float64(ev.Y), 0, 0, 0xff, 5)
						}
						PLRENDERER.Present()
						estimate = time.Now().UnixNano() - tmp
						last = time.Now().UnixNano()
					}
				}
				break
			}
		}
	}
}

func PrintInstruction(Inst string) {
	INSTRENDERER.SetDrawColor(0, 0, 0, 0)
	INSTRENDERER.FillRect(nil)
	surface, err := FONT.RenderUTF8BlendedWrapped(Inst, WHITE, 400)
	texture, err := INSTRENDERER.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	messageRect := sdl.Rect{0, 0, surface.W, surface.H}
	INSTRENDERER.Copy(texture, nil, &messageRect)
	INSTRENDERER.Present()
}

func LetsWatch(lpath []polygon.Vertex) {
	splittedPath := polygon.SplitPathIntoPixels(lpath, 10)
	for _, v := range splittedPath {
		PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
		PLRENDERER.FillRect(nil)
		DrawSpace(&spaceInit, 0x66, 0x66, 0x66, 5)
		DrawListVectFromVert(&polRobotVect, v.X, v.Y, 0x66, 0x66, 0x66, 5)
		PLRENDERER.Present()
		//time.Sleep(40000 * time.Nanosecond)
	}
}

func Menu() {
	PrintInstruction(INSTRUCTION_MENU)
	showBorder := false
	showCells := false
	showPath := false
	showVisibility := false
	listOfEdges := &([]polygon.Edge{})
	robDec := &([]*polygon.Polygon{})
	accPol := &([]polygon.Polygon{})
	cells := []*polygon.Polygon{}
	g := polygon.NewGraph()
	vg := polygon.VGraph{}.NewGraph()
	lpath := []polygon.Vertex{}
	turns := 0
	for quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch ev := event.(type) {
			case *sdl.QuitEvent:
				quit = false
				break
			case *sdl.KeyboardEvent:
				if ev.Type == sdl.KEYDOWN {
					switch ev.Keysym.Sym {
					case sdl.GetKeyFromName("P"):

						break
					case sdl.GetKeyFromName("Q"):
						quit = false
						break
					case sdl.GetKeyFromName("W"):
						LetsWatch(lpath)
						break
					case sdl.GetKeyFromName("E"):
						space = polygon.NewEmptySpace()
						spaceConf = polygon.NewEmptySpace()
						spaceAcc = polygon.NewEmptySpace()
						robSpace = polygon.NewEmptySpace()
						spaceInit = polygon.NewEmptySpace()
						allEdge = []polygon.Edge{}
						allVert = []polygon.Vertex{}
						robVert = []polygon.Vertex{}
						robEdge = []polygon.Edge{}
						polRobot = polygon.Polygon{}
						polRobotVect = []polygon.Vertex{}
						showBorder = false
						showCells = false
						showPath = false
						showVisibility = false
						mink = false
						setr = false
						seto = false
						listOfEdges = &([]polygon.Edge{})
						robDec = &([]*polygon.Polygon{})
						accPol = &([]polygon.Polygon{})
						cells = []*polygon.Polygon{}
						g = polygon.NewGraph()
						vg = polygon.VGraph{}.NewGraph()
						lpath = []polygon.Vertex{}
						turns = 0
						break
					case sdl.GetKeyFromName("N"):
						switch turns {
						case 0:
							break
						case 1:
							spaceInit = space
							polygon.AddPolygonToSpace(&frame_polygon1, &space)
							polygon.AddPolygonToSpace(&frame_polygon2, &space)
							SetLines()
							turns++
							break
						case 2:
							if !mink {
								spaceConf = polygon.NewEmptySpace()
								robDec = DecomposeRobot()
								listOfPol := polygon.MinkowskiGeneral(vertexBegin, robDec, &space)
								for _, pol := range *listOfPol {
									polygon.AddPolygonToSpace(pol, &spaceConf)
								}
								mink = true
							} else {
								mink = false
							}
							turns++
							break
						case 3:
							listOfEdges = polygon.FasterBorderEdges(&spaceConf)
							fmt.Printf("\n%v\n", *listOfEdges)
							accPol = polygon.FromEdgesToPoly(listOfEdges)
							fmt.Printf("\n%v\n", *accPol)
							showBorder = true
							for _, pol := range *accPol {
								polygon.AddPolygonToSpace(&pol, &spaceAcc)
							}
							polygon.DeleteFirstPol(&spaceAcc)
							turns++
							break

						}
					case sdl.GetKeyFromName("R"):
						space = spaceInit
						spaceConf = polygon.NewEmptySpace()
						spaceAcc = polygon.NewEmptySpace()
						robSpace = polygon.NewEmptySpace()
						showBorder = false
						showVisibility = false
						showCells = false
						showPath = false
						mink = false
						polRobot = polygon.Polygon{}
						polRobotVect = []polygon.Vertex{}
						robVert = []polygon.Vertex{}
						robEdge = []polygon.Edge{}
						PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
						PLRENDERER.FillRect(nil)
						DrawSpace(&space, 0x66, 0x66, 0x66, 5)
						PLRENDERER.Present()
						GetRobot()
						PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
						PLRENDERER.FillRect(nil)
						DrawSpace(&space, 0x66, 0x66, 0x66, 5)
						DrawListVectFromVert(&polRobotVect, vertexBegin.X, vertexBegin.Y, 0, 0xff, 0, 5)
						PLRENDERER.Present()
						SetBeginEnd()
						setr = true
						if seto {
							turns = 1
						}
						spaceConf = polygon.NewEmptySpace()
						break
					case sdl.GetKeyFromName("O"):
						space = spaceInit
						spaceConf = polygon.NewEmptySpace()
						spaceAcc = polygon.NewEmptySpace()
						showBorder = false
						showVisibility = false
						showCells = false
						showPath = false
						mink = false
						PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
						PLRENDERER.FillRect(nil)
						DrawSpace(&space, 0x66, 0x66, 0x66, 5)
						PLRENDERER.Present()
						SetObstacles()
						seto = true
						if setr {
							turns = 1
						}
						spaceInit = space
						break
					case sdl.GetKeyFromName("S"):
						DrawSpace(&spaceInit, 0x66, 0x66, 0x66, 5)
						DrawListVectFromVert(&polRobotVect, vertexBegin.X, vertexBegin.Y, 0, 0xff, 0, 5)
						PLRENDERER.Present()
						SetBeginEnd()
						break
					case sdl.GetKeyFromName("C"):
						if turns == 4 {
							showVisibility = false
							cells, g = polygon.DecomposeAndBuildGraph(&spaceAcc)
							lpath, showPath = g.FindPath(vertexBegin, vertexEnd)
							showCells = true
							showPath = true
						}
						break
					case sdl.GetKeyFromName("V"):
						if turns == 4 {
							vg = polygon.VisibilityGraphFrom(&spaceAcc, &spaceConf)
							fmt.Printf("\n%v\n", vertexBegin)
							vg = vg.AddBeginEnd(&vertexBegin, &vertexEnd, &spaceAcc, &spaceConf)
							fmt.Printf("\n%v\n", vertexBegin)
							lpath = vg.SPDijkstra(vertexBegin.Id, vertexEnd.Id)
							fmt.Printf("\nYES\n")
							showVisibility = true
							showPath = true
						}
						break
					}
					PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
					PLRENDERER.FillRect(nil)
					DrawSpace(&space, 0x66, 0x66, 0x66, 5)
					if setr {
						DrawListVectFromVert(&polRobotVect, vertexBegin.X, vertexBegin.Y, 0, 0xff, 0, 5)
						DrawListVectFromVert(&polRobotVect, vertexEnd.X, vertexEnd.Y, 0, 0, 0xff, 5)
					}
					if mink {
						DrawSpace(&spaceConf, 0xff, 0, 0, 1)
						for _, polptr := range *robDec {
							DrawPolygon(polptr, 0xff, 0, 0, 5)
						}
					}
					if showBorder {
						for _, v := range *listOfEdges {
							bx, by := int32(v.Begin().X), int32(v.Begin().Y)
							ex, ey := int32(v.End().X), int32(v.End().Y)
							_ = gfx.ThickLineRGBA(PLRENDERER, bx, by, ex, ey, 5, 0, 0xff, 0, 0xff)

						}
						for _, pol := range *accPol {
							DrawPolygon(&pol, 0, 0, 0xff, 5)
						}
					}
					if showCells {
						for _, pol := range cells {
							DrawPolygon(pol, 0xff, 0, 0xff, 5)
						}
					}
					if showPath {
						for i := 0; i < len(lpath)-1; i++ {
							bx, by := int32(lpath[i].X), int32(lpath[i].Y)
							ex, ey := int32(lpath[i+1].X), int32(lpath[i+1].Y)
							_ = gfx.ThickLineRGBA(PLRENDERER, bx, by, ex, ey, 5, 0xff, 0, 0, 0xff)
						}
					}
					if showVisibility {
						for i, v1 := range vg.Vertices {
							for _, v2 := range vg.AdjList[i] {
								v2x, v2y := int32(vg.Vertices[v2].X), int32(vg.Vertices[v2].Y)
								_ = gfx.ThickLineRGBA(PLRENDERER, int32(v1.X), int32(v1.Y), v2x, v2y, 3, 0x11, 0xff, 0x11, 0xff)
							}
						}
					}

					PLRENDERER.Present()

				}
				PrintInstruction(INSTRUCTION_MENU)
				break
			}

		}
	}
}

func handleClickObstacles(ev *sdl.MouseButtonEvent, listOfClicks *[]polygon.Vertex, i *int) {
	fmt.Printf("\n%v\n", *i)
	if ev.State == sdl.PRESSED && ev.WindowID == MAIN_WINDOW_ID {
		vclick := polygon.Vertex{uint32(*i), float64(ev.X), float64(ev.Y)}
		l := append(allVert, robVert...)
		if !vclick.SameX(&l) {
			fmt.Printf("%v\n\n", []int32{ev.X, ev.Y})
			if *i > 0 {
				eclick := polygon.Edge{(*listOfClicks)[len(*listOfClicks)-1], vclick}
				l2 := append(allEdge, robEdge...)
				if !eclick.IntersectsOne(&l2) {
					*listOfClicks = append(*listOfClicks, vclick)
					allVert = append(allVert, vclick)
					allEdge = append(allEdge, eclick)
					drawn := gfx.ThickLineRGBA(PLRENDERER,
						int32((*listOfClicks)[*i-1].X),
						int32((*listOfClicks)[*i-1].Y),
						ev.X, ev.Y, 5, 0x66, 0x66, 0x66, 0xff)
					if !drawn {
						fmt.Printf("Not Drawn")
					}
					(*i)++
				}
			} else {
				*listOfClicks = append(*listOfClicks, vclick)
				allVert = append(allVert, vclick)
				(*i)++
			}
		}
	}
}

func handleKeysObstacles(ev *sdl.KeyboardEvent, listOfClicks *[]polygon.Vertex, i *int, runningptr *bool, running *bool) {
	if ev.Type == sdl.KEYDOWN {
		switch ev.Keysym.Sym {
		case sdl.GetKeyFromName("D"):
			if *i < 3 {
				break
			}
			e := polygon.Edge{polygon.Vertex{0, (*listOfClicks)[*i-1].X, (*listOfClicks)[*i-1].Y}, polygon.Vertex{0, (*listOfClicks)[0].X, (*listOfClicks)[0].Y}}
			l := append(allEdge, robEdge...)
			if !e.IntersectsOne(&l) {
				polygon.AddListToSpaceAsPoly(listOfClicks,
					&space)
				PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0)
				PLRENDERER.FillRect(nil)
				DrawSpace(&space, 0x66, 0x66, 0x66, 5)
				DrawListVectFromVert(&polRobotVect, vertexBegin.X, vertexBegin.Y, 0, 0xff, 0, 5)
				DrawListVectFromVert(&polRobotVect, vertexEnd.X, vertexEnd.Y, 0, 0, 0xff, 5)
				*listOfClicks = []polygon.Vertex{}
				*i = 0
				allEdge = append(allEdge, e)
			}
			break
		case sdl.GetKeyFromName("M"):
			*running = false
			break
		case sdl.GetKeyFromName("Q"):
			*runningptr = false
			break
		}
	}
}

func SetObstacles() {
	running := true
	listOfClicks := []polygon.Vertex{}
	i := 0
	PrintInstruction(INSTRUCTION_OBSTACLES)
	last := time.Now().UnixNano()
	estimate := time.Now().UnixNano() - last
	for running && quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				quit = false
				break
			case *sdl.MouseButtonEvent:
				handleClickObstacles(ev, &listOfClicks, &i)
				break
			case *sdl.KeyboardEvent:
				handleKeysObstacles(ev, &listOfClicks, &i, &quit, &running)
				break
			case *sdl.MouseMotionEvent:
				if ev.WindowID == MAIN_WINDOW_ID {
					if i > 0 {
						PLRENDERER.SetDrawColor(0xff, 0xff, 0xff, 0xff)
						PLRENDERER.FillRect(nil)
						_ = gfx.ThickLineRGBA(PLRENDERER, int32(listOfClicks[i-1].X),
							int32(listOfClicks[i-1].Y), ev.X, ev.Y, 5, 0x66, 0x66, 0x66, 0xff)
					}
				}
				break
			}
		}
		if time.Now().UnixNano()-last > estimate {
			tmp := time.Now().UnixNano()
			DrawSpace(&space, 0x66, 0x66, 0x66, 5)
			DrawListVert(&listOfClicks, 0x66, 0x66, 0x66, 5)
			DrawListVectFromVert(&polRobotVect, vertexBegin.X, vertexBegin.Y, 0, 0xff, 0, 5)
			DrawListVectFromVert(&polRobotVect, vertexEnd.X, vertexEnd.Y, 0, 0, 0xff, 5)
			PLRENDERER.Present()
			estimate = time.Now().UnixNano() - tmp
			last = time.Now().UnixNano()
		}
	}

}

func main() {

	Init()

	Menu()

	Quit()
}

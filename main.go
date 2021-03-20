package main

import "fmt"

type Cubie struct {
	indexOrigin uint8
	orientation uint8
}

const (
	_90  = 3
	_180 = 2
	_270 = 1
)

const (
	_U = 0
	_F = 1
	_R = 2
	_B = 3
	_L = 4
	_D = 5
)

type Cubies_4 [4]*Cubie

func (cubie *Cubie) checkOrientation() bool {
	return cubie.orientation == 0
}

func (cubie *Cubie) rotateEdge() {
	cubie.orientation ^= 1
}

func (cubie *Cubie) rotateCornerClockwise() {
	cubie.orientation++
	cubie.orientation %= 3
}

func (cubie *Cubie) rotateCornerCounterClockwise() {
	cubie.orientation += 3
	cubie.orientation--
	cubie.orientation %= 3
}

//UF UR UB UL DF DR DB DL FR FL BR BL UFR URB UBL ULF DRF DFL DLB DBR

func (cubies *Cubies_4) rotate(rot uint8) {
	var tmpCubies [4]Cubie
	for i := uint8(0); i < 4; i++ {
		tmpCubies[i] = *cubies[(i+rot)%4]
	}
	for i := 0; i < 4; i++ {
		*cubies[i] = tmpCubies[i]
	}
}

type Face struct {
	edges   Cubies_4
	corners Cubies_4
}

func (face *Face) rotate(rot uint8) {
	face.edges.rotate(rot)
	face.corners.rotate(rot)
}

type Cube struct {
	edges   [12]Cubie
	corners [8]Cubie
	faces   [6]Face //U F R B L D
}

func (cube *Cube) rotateFace(faceIndex uint8, rot uint8) {
	switch faceIndex {
	case _U, _D:
		if rot != _180 {
			for _, v := range cube.faces[faceIndex].edges {
				v.rotateEdge()
			}
		}
	case _F, _B:
		if rot != _180 {
			for i, v := range cube.faces[faceIndex].corners {
				if i%2 == 0 {
					v.rotateCornerClockwise()
				} else {
					v.rotateCornerCounterClockwise()
				}
			}
		}
	}
	cube.faces[faceIndex].rotate(rot)
}

func (cube *Cube) getEdges(indexes []uint8) (cubies Cubies_4) {
	for i, v := range indexes {
		cubies[i] = &cube.edges[v]
	}
	return cubies
}
func (cube *Cube) getCorners(indexes []uint8) (cubies Cubies_4) {
	for i, v := range indexes {
		cubies[i] = &cube.corners[v]
	}
	return cubies
}

func (cube *Cube) initialize() {
	cube.faces[0].edges = cube.getEdges([]uint8{3, 2, 1, 0})
	cube.faces[1].edges = cube.getEdges([]uint8{0, 8, 4, 9})
	cube.faces[2].edges = cube.getEdges([]uint8{1, 10, 5, 8})
	cube.faces[3].edges = cube.getEdges([]uint8{2, 11, 6, 10})
	cube.faces[4].edges = cube.getEdges([]uint8{3, 9, 7, 11})
	cube.faces[5].edges = cube.getEdges([]uint8{4, 5, 6, 7})

	cube.faces[0].corners = cube.getCorners([]uint8{3, 2, 1, 0})
	cube.faces[1].corners = cube.getCorners([]uint8{3, 0, 4, 5})
	cube.faces[2].corners = cube.getCorners([]uint8{0, 1, 6, 4})
	cube.faces[3].corners = cube.getCorners([]uint8{1, 2, 7, 6})
	cube.faces[4].corners = cube.getCorners([]uint8{2, 3, 5, 7})
	cube.faces[5].corners = cube.getCorners([]uint8{4, 6, 7, 5})

	for i := range cube.edges {
		cube.edges[i].indexOrigin = uint8(i)
	}
	for i := range cube.corners {
		cube.corners[i].indexOrigin = uint8(i)
	}
}

func main() {
	cube := Cube{}
	cube.initialize()

	fmt.Println(cube.edges)
	fmt.Println(cube.corners)
	fmt.Println("--")
	cube.rotateFace(_F, _90)
	fmt.Println(cube.edges)
	fmt.Println(cube.corners)
	fmt.Println("--")
}

//EDGES
//0 U
// 3 2 1 0

//1 F
//0 8 4 9

//2 R
//1 10 5 8

//3 B
//2 11 6 10

//4 L
//3 9 7 11

//5 D
//4 5 6 7

//CORNERS
//0 U
// 3 2 1 0

//1 F
//3 0 4 5

//2 R
//0 1 6 4

//3 B
//1 2 7 6

//4 L
//2 3 5 7

//5 D
//4 6 7 5

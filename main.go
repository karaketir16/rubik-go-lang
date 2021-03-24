package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Cubie struct {
	indexOrigin byte
	orientation byte
}

const (
	_90         = 0
	_180        = 1
	_270        = 2
	_MAX_DEGREE = 3
)

const (
	_U        = 0
	_F        = 1
	_R        = 2
	_B        = 3
	_L        = 4
	_D        = 5
	_MAX_FACE = 6
)

type faceRotPair struct {
	face byte
	rot  byte
}

var phase1grp []faceRotPair
var phase2grp []faceRotPair
var phase3grp []faceRotPair
var phase4grp []faceRotPair

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

type byte_4 [4]byte

type Face struct {
	edges   byte_4
	corners byte_4
}

func (arr byte_4) rotate(cubies []Cubie, rot byte) {

	switch rot {
	case _90:
		rot = 3
	case _180:
		rot = 2
	case _270:
		rot = 1
	}

	result := [4]Cubie{}

	for i := byte(0); i < 4; i++ {
		result[i] = cubies[arr[(i+rot)%4]]
	}
	for i := byte(0); i < 4; i++ {
		cubies[arr[i]] = result[i]
	}
	return
}

func (face *Face) rotate(cube *Cube, rot byte) {
	face.edges.rotate(cube.edges[:], rot)
	face.corners.rotate(cube.corners[:], rot)
}

type Cube struct {
	edges   [12]Cubie
	corners [8]Cubie
}

func (cube Cube) rotateFace(faceIndex byte, rot byte) Cube {
	switch faceIndex {
	case _U, _D:
		if rot != _180 {
			for _, v := range faces[faceIndex].edges {
				cube.edges[v].rotateEdge()
				//print("Hello")
			}
		}
	case _F, _B:
		if rot != _180 {
			for i, v := range faces[faceIndex].corners {
				if i%2 == 0 {
					cube.corners[v].rotateCornerClockwise()
				} else {
					cube.corners[v].rotateCornerCounterClockwise()
				}
			}
		}
	}
	faces[faceIndex].rotate(&cube, rot)
	return cube
}

type Faces [6]Face //U F R B L D

func (faces *Faces) initialize() {
	faces[0].edges = [4]byte{3, 2, 1, 0}
	faces[1].edges = [4]byte{0, 8, 4, 9}
	faces[2].edges = [4]byte{1, 10, 5, 8}
	faces[3].edges = [4]byte{2, 11, 6, 10}
	faces[4].edges = [4]byte{3, 9, 7, 11}
	faces[5].edges = [4]byte{4, 5, 6, 7}

	faces[0].corners = [4]byte{3, 2, 1, 0}
	faces[1].corners = [4]byte{3, 0, 4, 5}
	faces[2].corners = [4]byte{0, 1, 6, 4}
	faces[3].corners = [4]byte{1, 2, 7, 6}
	faces[4].corners = [4]byte{2, 3, 5, 7}
	faces[5].corners = [4]byte{4, 6, 7, 5}
}

func (cube *Cube) initialize() {
	for i := range cube.edges {
		cube.edges[i].indexOrigin = byte(i)
	}
	for i := range cube.corners {
		cube.corners[i].indexOrigin = byte(i)
	}
}

var visited map[Cube]bool
var visited_lock = sync.Mutex{}

type Qtype struct {
	cube  *Cube
	depth int
}

var queue chan Qtype

var faces Faces

func BFS(val Qtype, rots []faceRotPair) (Cube, int) {
	for _, v := range rots {
		i := v.face
		j := v.rot
		rotatedCube := val.cube.rotateFace(i, j)
		visited_lock.Lock()
		if !visited[rotatedCube] {

			select {
			case queue <- Qtype{&rotatedCube, val.depth + 1}:
				visited[*val.cube] = true
			default:
			}
		}
		visited_lock.Unlock()

	}
	return *val.cube, val.depth
}

func checkPhase1(cube *Cube) bool {
	for _, v := range cube.edges {
		if !v.checkOrientation() {
			return false
		}
	}
	return true
}

func byteInSlice(a byte, list []byte) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkPhase2(cube *Cube) bool {
	for _, v := range cube.corners {
		if !v.checkOrientation() {
			return false
		}
	}

	middleLayer := []byte{0, 2, 4, 6}
	for _, v := range middleLayer {
		if !byteInSlice(cube.edges[v].indexOrigin, middleLayer) {
			return false
		}
	}

	return true
}

func solvePhase(checker func(*Cube) bool, cube Cube, rots []faceRotPair) Cube {
	visited = make(map[Cube]bool)
	queue = make(chan Qtype, 100000)
	queue <- Qtype{&cube, 0}
	visited[cube] = true

	routineCount := 4

	stop := make(chan byte, routineCount)
	done := make(chan byte, routineCount)
	answer := make(chan Cube, 1)

	searcher := func() {
	out:
		for {
			select {
			case val := <-queue:
				result, _ := BFS(val, rots)
				if checker(&result) {
					for i := 0; i < routineCount; i++ {
						select {
						case done <- 0:
						default:
						}
					}
					select {
					case answer <- result:
					default:
						//someone else found
					}
					break out
				}
			case <-done:
				break out
			}
		}
		stop <- 0
		//fmt.Println("I am done")
	}
	for i := 0; i < routineCount; i++ {
		go searcher()
	}
	for i := 0; i < routineCount; i++ {
		<-stop
	}
	return <-answer
}

func main() {

	for i := byte(0); i < _MAX_FACE; i++ {
		for j := byte(0); j < _MAX_DEGREE; j++ {
			phase1grp = append(phase1grp, faceRotPair{i, j})
		}
	}
	for i := byte(0); i < _MAX_FACE; i++ {
		for j := byte(0); j < _MAX_DEGREE; j++ {
			if i == _U || i == _D {
				if j == _180 {
					phase2grp = append(phase2grp, faceRotPair{i, j})
				}
			} else {
				phase2grp = append(phase2grp, faceRotPair{i, j})
			}
		}
	}

	cube := Cube{}
	cube.initialize()
	faces.initialize()
	rand.Seed(time.Now().UnixNano())
	fmt.Println("start", cube)
	for i := 0; i < 50; i++ {
		cube = cube.rotateFace(byte(rand.Intn(_MAX_FACE)), byte(rand.Intn(_MAX_DEGREE)))
	}
	fmt.Println("scrambled", cube)
	phase1Solved := solvePhase(checkPhase1, cube, phase1grp)
	fmt.Println("phase1", phase1Solved)
	phase2Solved := solvePhase(checkPhase2, phase1Solved, phase2grp)
	fmt.Println("phase2", phase2Solved)
}

//func test() {
//	mx := 0
//	for {
//		cube := Cube{}
//		cube.initialize()
//		faces.initialize()
//		rand.Seed(time.Now().UnixNano())
//
//		for i := 0; i < 50; i++ {
//			cube = cube.rotateFace(byte(rand.Intn(_MAX_FACE)), byte(rand.Intn(_MAX_DEGREE)))
//		}
//
//		visited = make(map[Cube]bool)
//		queue = make(chan Qtype, 100000)
//		queue <- Qtype{&cube, 0}
//		visited[cube] = true
//
//		routineCount := 4
//
//		stop := make(chan byte, routineCount)
//		done := make(chan byte, routineCount)
//
//		searcher := func() {
//		out:
//			for {
//				select {
//				case val := <-queue:
//					result, depth := BFS(val)
//					if result.checkPhase1() {
//						//fmt.Println("Found")
//						if depth > mx {
//							mx = depth
//						}
//						fmt.Println("queue len ", len(queue), "|| depth ", depth, " || ", mx)
//
//						for i := 0; i < routineCount; i++ {
//							select {
//							case done <- 0:
//							default:
//							}
//						}
//						break out
//					}
//				case <-done:
//					break out
//				}
//			}
//			stop <- 0
//			//fmt.Println("I am done")
//		}
//		for i := 0; i < routineCount; i++ {
//			go searcher()
//		}
//		for i := 0; i < routineCount; i++ {
//			<-stop
//		}
//	}
//	//fmt.Print("Done")
//	//fmt.Println(cube)
//}

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

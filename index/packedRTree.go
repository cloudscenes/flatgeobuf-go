package index

import (
	"encoding/binary"
	"fmt"
	"math"
)

type NodeItem struct {
	minX   float64
	minY   float64
	maxX   float64
	maxY   float64
	offset uint64
}

func (ni *NodeItem) expand(n NodeItem) {
	if n.minX < ni.minX {
		ni.minX = n.minX
	}
	if n.minY < ni.minY {
		ni.minY = n.minY
	}
	if n.maxX > ni.maxX {
		ni.maxX = n.maxX
	}
	if n.maxY > ni.maxY {
		ni.maxY = n.maxY
	}
}

const nodeItemLen = 8*4 + 8

func CalcTreeSize(numItems uint64, nodeSize uint16) (uint64, error) {
	if nodeSize < 2 {
		return 0, fmt.Errorf("node size must be at least 2")
	}
	if numItems == 0 {
		return 0, fmt.Errorf("number of items must be greater than 0")
	}

	minNodeSize := uint64(nodeSize)
	n := numItems
	numNodes := n
	for {
		n = (n + minNodeSize - 1) / minNodeSize
		numNodes += n
		if n == 1 {
			break
		}
	}

	return numNodes * nodeItemLen, nil
}

type LevelBounds struct {
	start uint64
	end   uint64
}

type PackedRTree struct {
	extent      NodeItem
	items       []NodeItem
	numItems    uint64
	numNodes    uint64
	nodeSize    uint16
	levelBounds []LevelBounds
}

func NewPackedRTree(numItems uint64, nodeSize uint16, data []byte) (*PackedRTree, error) {
	if nodeSize < 2 {
		return nil, fmt.Errorf("node size must be at least 2")
	}
	if numItems == 0 {
		return nil, fmt.Errorf("number of items must be greater than 0")
	}

	numNodes, err := CalcTreeSize(numItems, nodeSize)
	if err != nil {
		return nil, err
	}

	prt := &PackedRTree{
		extent:      NodeItem{},
		items:       make([]NodeItem, 0, numNodes),
		numItems:    numItems,
		numNodes:    numNodes,
		nodeSize:    nodeSize,
		levelBounds: nil,
	}

	for i := 0; i < len(data); i += 8 * 5 {
		ni := NodeItem{}
		minxb := binary.LittleEndian.Uint64(data[i : i+8])
		ni.minX = math.Float64frombits(minxb)
		minyb := binary.LittleEndian.Uint64(data[i+8 : i+8*2])
		ni.minY = math.Float64frombits(minyb)
		maxxb := binary.LittleEndian.Uint64(data[i+8*2 : i+8*3])
		ni.maxX = math.Float64frombits(maxxb)
		maxyb := binary.LittleEndian.Uint64(data[i+8*3 : i+8*4])
		ni.maxY = math.Float64frombits(maxyb)
		ni.offset = binary.LittleEndian.Uint64(data[i+8*4 : i+8*5])
		prt.items = append(prt.items, ni)
	}

	prt.generateLevelBounds()

	return prt, nil
}

func (prt *PackedRTree) generateLevelBounds() {
	minNodeSize := uint64(prt.nodeSize)
	n := prt.numItems
	numNodes := n
	var levelNumNodes []uint64
	levelNumNodes = append(levelNumNodes, n)
	for {
		n = (n + minNodeSize - 1) / minNodeSize
		numNodes += n
		levelNumNodes = append(levelNumNodes, n)
		if n == 1 {
			break
		}
	}

	levelOffsets := make([]uint64, 0, len(levelNumNodes))
	n = numNodes
	for _, size := range levelNumNodes {
		levelOffsets = append(levelOffsets, n-size)
		n -= size
	}

	levelBounds := make([]LevelBounds, 0, len(levelNumNodes))
	for i := len(levelNumNodes) - 1; i >= 0; i-- {
		levelBounds = append(levelBounds, LevelBounds{levelOffsets[i], levelOffsets[i] + levelNumNodes[i]})
	}

	prt.levelBounds = levelBounds
}

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

func (ni *NodeItem) readFromBytes(b []byte) {
	minxb := binary.LittleEndian.Uint64(b[:8])
	ni.minX = math.Float64frombits(minxb)
	minyb := binary.LittleEndian.Uint64(b[8 : 8*2])
	ni.minY = math.Float64frombits(minyb)
	maxxb := binary.LittleEndian.Uint64(b[8*2 : 8*3])
	ni.maxX = math.Float64frombits(maxxb)
	maxyb := binary.LittleEndian.Uint64(b[8*3 : 8*4])
	ni.maxY = math.Float64frombits(maxyb)
	ni.offset = binary.LittleEndian.Uint64(b[8*4 : 8*5])
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

func (ni *NodeItem) intersects(n NodeItem) bool {
	if ni.maxX < n.minX {
		return false
	}
	if ni.maxY < n.minY {
		return false
	}
	if ni.minX > n.maxX {
		return false
	}
	if ni.minY > n.maxY {
		return false
	}
	return true
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
	extent      *NodeItem
	items       []NodeItem
	numItems    uint64
	numNodes    uint64
	nodeSize    uint16
	levelBounds []LevelBounds
}

func ReadPackedRTreeBytes(numItems uint64, nodeSize uint16, data []byte) (*PackedRTree, error) {
	numNodes, err := CalcTreeSize(numItems, nodeSize)
	if err != nil {
		return nil, err
	}

	prt := &PackedRTree{
		extent:      nil,
		items:       make([]NodeItem, 0, numItems),
		numItems:    numItems,
		numNodes:    numNodes,
		nodeSize:    nodeSize,
		levelBounds: nil,
	}

	for i := 0; i < len(data); i += 8 * 5 {
		ni := NodeItem{}
		ni.readFromBytes(data[i : i+8*5])
		prt.items = append(prt.items, ni)
		if prt.extent == nil {
			prt.extent = &ni
		}
		prt.extent.expand(ni)
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

type SearchResultItem struct {
	Offset uint32
	index  uint64
}

func (prt *PackedRTree) Search(minX float64, minY float64, maxX float64, maxY float64) []SearchResultItem {
	bounds := NodeItem{minX, minY, maxX, maxY, 0}

	toSearch := []LevelBounds{prt.levelBounds[0]}
	result := make([]SearchResultItem, 0)

	lastBound := prt.levelBounds[len(prt.levelBounds)-1]

	for len(toSearch) > 0 {
		v := toSearch[0]
		toSearch = toSearch[1:]
		for pos := v.start; pos < v.end; pos++ {
			if !bounds.intersects(prt.items[pos]) {
				continue
			}
			if pos >= lastBound.start && pos <= lastBound.end {
				item := SearchResultItem{
					Offset: uint32(prt.items[pos].offset),
					index:  pos,
				}
				result = append(result, item)
			} else {
				lb := LevelBounds{prt.items[pos].offset, prt.items[pos].offset + 16}
				toSearch = append(toSearch, lb)
			}
		}
	}
	return result
}

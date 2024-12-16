package hash

import(
	"hash/crc32"
	"sort"
	"strconv"
)

type HashRing struct {
	nodes       []int 
	nodeMap    map[int]string
	replication int
}

//NewHashRing creates a new hash ring
func NewHashRing(replication int) *HashRing {
	return &HashRing{
		nodes :     []int{},
		nodeMap:    make(map[int]string),
		replication: replication,
	}
}

//AddNode adds a node to the hash ring
func (hr *HashRing) AddNode(node string) {
	for i:=0; i< hr.replication; i++ {
		hash:= int(crc32.ChecksumIEEE([]byte(node+strconv.Itoa(i))))
		hr.nodes = append(hr.nodes, hash)
		hr.nodeMap[hash] = node
	}	
	sort.Ints(hr.nodes)
}

//GetNode returns the node for a given key
func (hr *HashRing) GetNode(key string) string {
	hash := int(crc32.ChecksumIEEE([]byte(key)))
	idx := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i] >= hash
	})
	if idx == len(hr.nodes) {
		idx = 0
	}
	return hr.nodeMap[hr.nodes[idx]]
}

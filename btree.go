package btree

import (
	"bytes"
	"fmt"
	"sort"
)

type Item struct {
	Value int
}

func (i *Item) Less(item *Item) bool {
	return i.Value < item.Value
}

func (i *Item) High(item *Item) bool {
	return i.Value > item.Value
}

func (i *Item) Equal(item *Item) bool {
	return i.Value == item.Value
}

type Items []*Item

func (i *Items) InsertAt(item *Item, index int) {
	// 枠を一個増やす
	*i = append(*i, nil)

	// 最後のindex以外の場合、新しいItemを入れるスペースを作る
	if index < len(*i) {
		copy((*i)[index+1:], (*i)[index:])
	}

	(*i)[index] = item
}

func (i *Items) RemoveAt(index int) *Item {
	item := (*i)[index]
	copy((*i)[index:], (*i)[index+1:])
	(*i)[len(*i)-1] = nil
	*i = (*i)[:len(*i)-1]

	return item
}

func (i *Items) findItem(item *Item) (index int, exist bool) {
	index = sort.Search(len(*i), func(in int) bool {
		return item.Less((*i)[in])
	})
	// あるItemと等価だった場合は、一つ前のindexで返す
	if index > 0 && !((*i)[index-1].Less(item)) {
		return index - 1, true
	}

	return index, false
}

// extract - 指定されたindexまでのItemのみを抽出する
func (i *Items) extract(index int) {
	for n := index; n < len(*i); n++ {
		(*i)[n] = nil
	}
}

func (i *Items) print() string {
	var b bytes.Buffer
	b.WriteString("[ ")
	for n, item := range *i {
		b.WriteString(fmt.Sprintf("%d", item.Value))
		if n != len(*i)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(" ]")

	return b.String()
}

type Node struct {
	Items  Items
	Branch Branch
}

type Branch []*Node

// FindItem - ノードの要素、Branchの中に指定されたItemがあるかをチェック
func (n *Node) FindItem(item *Item) bool {
	if _, ok := n.Items.findItem(item); ok {
		return true
	}

	for _, cn := range n.Branch {
		if cn.FindItem(item) {
			return true
		}
	}

	return false
}

func Min(n *Node) *Item {
	if n == nil {
		return nil
	}

	// nodeの下にさらにnodeがあれば一番左のNodeまで探索
	for len(n.Branch) > 0 {
		n = n.Branch[0]
	}

	// nodeにItemが一つもなければnil
	if len(n.Items) == 0 {
		return nil
	}

	return n.Items[0]
}

func Max(n *Node) *Item {
	if n == nil {
		return nil
	}

	for len(n.Branch) > 0 {
		n = n.Branch[len(n.Branch)-1]
	}

	if len(n.Items) == 0 {
		return nil
	}

	return n.Items[len(n.Items)-1]
}

func (n *Node) Length() int {
	return len(n.Items)
}

func (n *Node) InsertAt(index int, node *Node) {
	n.Branch = append(n.Branch, nil)
	if index < len(n.Branch) {
		copy(n.Branch[index+1:], n.Branch[index:])
	}

	n.Branch[index] = node
}

//
//func (n *Node) Split(index int) (Item, *Node) {
//	items := n.Items
//	newNode := new(Node)
//	newNode.Items = append(newNode.Items, items[index+1:]...)
//
//	n.Items.extract(index)
//	if len(n.Branch) > 0 {
//		newNode.Branch = append(newNode.Branch, n.Branch[index+1:]...)
//		n.Branch.(i + 1)
//	}
//}

var nilChildren = make(Branch, 16)

func (b *Branch) extract(index int) {
	*b = (*b)[:index]

	// index以降を初期化
	for i := index; i < len(*b); i++ {
		(*b)[i] = nil
	}
}

//func (n *Node) find(value *Item) *Item {
//	i, found := n.Items.findItem(value)
//	if found
//}

type Btree struct {
	m    int
	num  int
	Root *Node
}

func NewBTree(m int) (*Btree, error) {
	if m < 2 {
		return nil, fmt.Errorf("m must be greater than 3, can't be %d", m)
	}

	return &Btree{
		m: m,
	}, nil
}

func (b *Btree) maxItemNum() int {
	return b.m*2 - 1
}

//func (b *Btree) InsertItem(item *Item) (*Item, bool, error) {
//	if item == nil {
//		return nil, false, fmt.Errorf("nil item is unacceptable")
//	}
//
//	if b.Root == nil {
//		b.Root = new(Node)
//		b.Root.Items = append(b.Root.Items, item)
//		b.num++
//	} else {
//		if len(b.Root.Items) > b.maxItemNum() {
//			item2, second := b.Root.split(b.maxItemNum() / 2)
//			old := b.Root
//
//			b.Root = new(Node)
//			b.Root.Items = append(b.Root.Items, item2)
//			b.Root.Branch = append(b.Root.Branch, old, second)
//		}
//	}
//
//	res := b.Root.InsertItem(b.maxItemNum(), item)
//	if res == nil {
//		b.num++
//		return res, false, nil
//	}
//
//	return res, true, nil
//}

func (b *Btree) Find(item *Item) bool {
	return b.Root.FindItem(item)
}

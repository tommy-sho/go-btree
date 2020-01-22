package btree

import "sort"

type Item struct {
	Value int
}

func (i *Item) Less(item *Item) bool {
	return i.Value < item.Value
}

func (i *Item) High(item *Item) bool {
	return i.Value > item.Value
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

func (i *Items) SearchItems(item *Item) (int, bool) {
	index := sort.Search(len(*i), func(in int) bool {
		return item.Less((*i)[in])
	})
	// あるItemと等価だった場合は、一つ前のindexで返す
	if index > 0 && !((*i)[index-1].Less(item)) {
		return index - 1, true
	}

	return index, false
}

type Node struct {
	Items    *Items
	Children []*Node
}

// FindItem - ノードの要素、Childrenの中に指定されたItemがあるかをチェック
func (n *Node) FindItem(item *Item) bool {
	if _, ok := n.Items.SearchItems(item); ok {
		return true
	}

	for _, cn := range n.Children {
		if cn.FindItem(item) {
			return true
		}
	}

	return false
}

type Btree struct {
	Root *Node
}

func NewBTree() *Btree {
	return &Btree{Root: nil}
}

func (b *Btree) Find() {
}

package btree

import (
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

func (i *Items) insertAt(item *Item, index int) {
	// 枠を一個増やす
	*i = append(*i, nil)

	// 最後のindex以外の場合、新しいItemを入れるスペースを作る
	if index < len(*i) {
		copy((*i)[index+1:], (*i)[index:])
	}

	(*i)[index] = item
}

func (i *Items) removeAt(index int) *Item {
	item := (*i)[index]
	copy((*i)[index:], (*i)[index+1:])
	(*i)[len(*i)-1] = nil
	*i = (*i)[:len(*i)-1]

	return item
}

// findItem - Itemsの中から検索対象のItemを検索する。なかった場合は、もっとも近く、小さいItemのindexを返す
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

type Node struct {
	Items  Items
	Branch Branch
}

type Branch []*Node

func (b *Branch) InsertAt(index int, n *Node) {
	*b = append(*b, nil)
	if index < len(*b) {
		copy((*b)[index+1:], (*b)[index:])
	}

	(*b)[index] = n
}

// findItem - ノードの要素、Branchの中に指定されたItemがあるかをチェック
func (n *Node) findItem(item *Item) bool {
	if _, ok := n.Items.findItem(item); ok {
		return true
	}

	for _, cn := range n.Branch {
		if cn.findItem(item) {
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

//insert - nodeにItemを挿入する。limitの数を超えていた場合、新しくNodeを作成する
//Itemがもしあれば、Itemを入れ替え、入れ替えるまえのItemを返却する
func (n *Node) insert(item *Item, limit int) *Item {
	i, ok := n.Items.findItem(item)
	if ok {
		res := n.Items[i]
		n.Items[i] = item
		return res
	}

	if len(n.Branch) == 0 {
		n.Items.insertAt(item, i)
		return nil
	}

	if n.shouldSplit(i, limit) {
		inTree := n.Items[i]
		switch {
		case item.Less(inTree):
			// no change, we want first split node
		case inTree.Less(item):
			i++ // we want second split node
		default:
			out := n.Items[i]
			n.Items[i] = item
			return out
		}
	}

	return n.Branch[i].insert(item, limit)
}

func (n *Node) insertAt(index int, node *Node) {
	n.Branch = append(n.Branch, nil)
	if index < len(n.Branch) {
		copy(n.Branch[index+1:], n.Branch[index:])
	}

	n.Branch[index] = node
}

func (n *Node) split(index int) (*Item, *Node) {
	// ちょうど分割の真ん中にあるItem
	item := n.Items[index]

	newNode := new(Node)
	newNode.Items = append(newNode.Items, n.Items[index+1:]...)

	n.Items.extract(index)
	if len(n.Branch) > 0 {
		newNode.Branch = append(newNode.Branch, n.Branch[index+1:]...)
		n.Branch.extract(index + 1)
	}

	return item, newNode
}

func (n *Node) shouldSplit(i, limit int) bool {
	if len(n.Branch[i].Items) <= limit {
		return false
	}

	// limitを超えて分割すべきだった場合、分割してBranchに登録
	center, newNode := n.Branch[i].split(limit / 2)
	n.Items.insertAt(center, i)
	n.Branch.InsertAt(i+1, newNode)

	return true
}

func (b *Branch) extract(index int) {
	*b = (*b)[:index]

	// index以降を初期化
	for i := index; i < len(*b); i++ {
		(*b)[i] = nil
	}
}

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
		m:    m,
		Root: nil,
	}, nil
}

func (b *Btree) maxItemNum() int {
	return b.m*2 - 1
}

func (b *Btree) InsertOrUpdateItem(item *Item) (*Item, error) {
	if item == nil {
		return nil, fmt.Errorf("nil item is unacceptable")
	}

	// RootのNodeがまだなければ作成
	if b.Root == nil {
		b.Root = &Node{}
		b.Root.Items = append(b.Root.Items, item)
		b.num++
	} else {
		// Rootに収まらなくなったら分割する
		if len(b.Root.Items) >= b.maxItemNum() {
			item2, second := b.Root.split(b.maxItemNum() / 2)
			old := b.Root

			b.Root = new(Node)
			b.Root.Items = append(b.Root.Items, item2)
			b.Root.Branch = append(b.Root.Branch, old, second)
		}
	}

	res := b.Root.insert(item, b.maxItemNum())
	if res == nil {
		b.num++
	}

	return res, nil
}

func (b *Btree) Find(item *Item) bool {
	return b.Root.findItem(item)
}

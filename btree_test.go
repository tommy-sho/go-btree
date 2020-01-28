package btree

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestItems_InsertAt(t *testing.T) {
	type args struct {
		item  *Item
		index int
	}
	tests := []struct {
		name string
		i    Items
		args args
		want Items
	}{
		{
			name: "insert_4",
			i: Items{
				{Value: 1},
				{Value: 2},
				{Value: 3},
			},
			args: args{
				item: &Item{
					Value: 4,
				},
				index: 3,
			},
			want: Items{
				{Value: 1},
				{Value: 2},
				{Value: 3},
				{Value: 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.i.insertAt(tt.args.item, tt.args.index)
			if !reflect.DeepEqual(tt.i, tt.want) {
				t.Errorf("Items should be= %v(len=%d),but  %v(lend=%d)", tt.want, len(tt.want), tt.i, len(tt.i))
			}
		})
	}
}

func TestItems_RemoveAt(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name      string
		i         Items
		args      args
		want      *Item
		wantItems Items
	}{
		{
			name: "remove_index_1",
			i: Items{
				{Value: 1},
				{Value: 2},
				{Value: 3},
			},
			args: args{
				index: 1,
			},
			want: &Item{
				Value: 2,
			},
			wantItems: Items{
				{Value: 1},
				{Value: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.removeAt(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeAt() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.i, tt.wantItems) {
				t.Errorf("remain Items should be= %v,but remain %v", tt.i, tt.wantItems)
			}
		})
	}
}

func TestItems_findItem(t *testing.T) {
	type args struct {
		item *Item
	}
	tests := []struct {
		name  string
		i     Items
		args  args
		want  int
		want1 bool
	}{
		{
			name: "find_equal_item",
			i: Items{
				{Value: 0},
				{Value: 10},
				{Value: 100},
				{Value: 200},
				{Value: 1000},
			},
			args: args{
				item: &Item{Value: 10},
			},
			want:  1,
			want1: true,
		},
		{
			name: "not_found",
			i: Items{
				{Value: 0},
				{Value: 10},
				{Value: 100},
				{Value: 200},
				{Value: 1000},
			},
			args: args{
				item: &Item{Value: 11},
			},
			want:  2,
			want1: false,
		},
		{
			name: "not_found_more_than_maximum",
			i: Items{
				{Value: 0},
				{Value: 10},
				{Value: 100},
				{Value: 200},
				{Value: 1000},
			},
			args: args{
				item: &Item{Value: 10000},
			},
			want:  5,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.i.findItem(tt.args.item)
			if got != tt.want {
				t.Errorf("findItem() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findItem() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNode_FindItem(t *testing.T) {
	type fields struct {
		Items  Items
		Branch []*Node
	}
	type args struct {
		item *Item
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "exist",
			fields: fields{
				Items: Items{
					{Value: 10},
					{Value: 100},
				},
				Branch: []*Node{{
					Items: Items{{Value: 5}},
				}},
			},
			args: args{
				item: &Item{Value: 5},
			},
			want: true,
		},
		{
			name: "not_exist",
			fields: fields{
				Items: Items{
					{Value: 10},
					{Value: 100},
				},
				Branch: []*Node{{
					Items: Items{{Value: 5}},
				}},
			},
			args: args{
				item: &Item{Value: 5},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Items:  tt.fields.Items,
				Branch: tt.fields.Branch,
			}
			if got := n.findItem(tt.args.item); got != tt.want {
				t.Errorf("findItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_InsertAt(t *testing.T) {
	type fields struct {
		Items  Items
		Branch []*Node
	}
	type args struct {
		index int
		node  *Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Node
	}{
		{
			name: "",
			fields: fields{
				Items: Items{{Value: 5}},
				Branch: []*Node{
					{Items: Items{{Value: 1}}},
				},
			},
			args: args{
				index: 1,
				node: &Node{
					Items: Items{{Value: 10}},
				},
			},
			want: &Node{
				Items: Items{{Value: 5}},
				Branch: []*Node{
					{Items: Items{{Value: 1}}},
					{Items: Items{{Value: 10}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Items:  tt.fields.Items,
				Branch: tt.fields.Branch,
			}
			n.insertAt(tt.args.index, tt.args.node)
			if diff := cmp.Diff(n, tt.want); diff != "" {
				t.Errorf("TestNode_InsertAt() diff = %v", diff)
			}
		})
	}
}

func TestItems_extract(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		i    Items
		args args
		want Items
	}{
		{
			i: Items{
				{Value: 1},
				{Value: 2},
				{Value: 3},
				{Value: 4},
				{Value: 5},
			},
			args: args{
				index: 2,
			},
			want: Items{
				{Value: 1},
				{Value: 2},
				nil,
				nil,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.i.extract(tt.args.index)
			if !reflect.DeepEqual(tt.i, tt.want) {
				t.Errorf("TestNode_extract() got = %v, but want = %v", tt.i, tt.want)
			}
		})
	}
}

func TestBranch_extract(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		b    Branch
		args args
		want Branch
	}{
		{
			name: "",
			b: Branch{
				{Items: Items{{Value: 1}}},
				{Items: Items{{Value: 2}}},
				{Items: Items{{Value: 3}}},
				{Items: Items{{Value: 4}}},
				{Items: Items{{Value: 5}}},
			},
			args: args{
				index: 3,
			},
			want: Branch{
				{Items: Items{{Value: 1}}},
				{Items: Items{{Value: 2}}},
				{Items: Items{{Value: 3}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.extract(tt.args.index)
			if !reflect.DeepEqual(tt.b, tt.want) {
				t.Errorf("TestBranch_extract() got = %v, but want = %v", tt.b, tt.want)
			}
		})
	}
}

func TestNode_Split(t *testing.T) {
	type fields struct {
		Items  Items
		Branch Branch
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Item
		want1  *Node
		want2  *Node
	}{
		{
			name: "[5_10_15]",
			fields: fields{
				Items: Items{{Value: 5}, {Value: 10}, {Value: 15}},
				Branch: Branch{
					{Items: Items{{Value: 4}}},
					{Items: Items{{Value: 8}}},
					{Items: Items{{Value: 12}}},
					{Items: Items{{Value: 20}}},
				},
			},
			args: args{
				index: 1,
			},
			want: &Item{
				Value: 10,
			},
			want1: &Node{
				Items: Items{{Value: 15}},
				Branch: Branch{
					{Items: Items{{Value: 12}}},
					{Items: Items{{Value: 20}}},
				},
			},
			want2: &Node{
				Items: Items{{Value: 5}, nil, nil},
				Branch: Branch{
					{Items: Items{{Value: 4}}},
					{Items: Items{{Value: 8}}},
					nil,
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Items:  tt.fields.Items,
				Branch: tt.fields.Branch,
			}
			got, got1 := n.split(tt.args.index)
			tt.want2.Branch = tt.want2.Branch[:tt.args.index+1]
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("split() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(n, tt.want2) {
				t.Errorf("split() original node = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNode_ShouldSplit(t *testing.T) {
	type fields struct {
		Items  Items
		Branch Branch
	}
	type args struct {
		i     int
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "",
			fields: fields{
				Items: Items{{Value: 2}, {Value: 4}, nil},
				Branch: Branch{
					{Items: Items{{Value: 1}, nil, nil}},
					{Items: Items{{Value: 3}, nil, nil}},
					{Items: Items{{Value: 5}, {Value: 8}, nil}},
				},
			},
			args: args{
				i:     2,
				limit: 2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Items:  tt.fields.Items,
				Branch: tt.fields.Branch,
			}
			if got := n.shouldSplit(tt.args.i, tt.args.limit); got != tt.want {
				t.Errorf("shouldSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_AppendItem(t *testing.T) {
	type fields struct {
		Items  Items
		Branch Branch
	}
	type args struct {
		item  *Item
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Item
	}{
		{
			name: "",
			fields: fields{
				Items: Items{{Value: 5}, {Value: 10}, {Value: 15}},
				Branch: Branch{
					{Items: Items{{Value: 1}, {Value: 2}}},
					{Items: Items{{Value: 6}, {Value: 7}}},
					{Items: Items{{Value: 13}, {Value: 14}}},
					{Items: Items{{Value: 17}, {Value: 20}}},
				},
			},
			args: args{
				item: &Item{
					Value: 3,
				},
				limit: 3,
			},
			want: &Item{
				Value: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Items:  tt.fields.Items,
				Branch: tt.fields.Branch,
			}
			if got := n.AppendItem(tt.args.item, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

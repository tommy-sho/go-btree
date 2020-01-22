package btree

import (
	"reflect"
	"testing"
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
			i:     Items{
				{Value: 1},
				{Value:2},
				{Value:3},
			},
			args: args{
				item: &Item{
					Value: 4,
				},
				index: 3,
			},
			want: Items{
				{Value: 1},
				{Value:2},
				{Value:3},
				{Value:4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.i.InsertAt(tt.args.item,tt.args.index)
			if !reflect.DeepEqual(tt.i, tt.want) {
				t.Errorf("Items should be= %v(len=%d),but  %v(lend=%d)", tt.want, len(tt.want),tt.i, len(tt.i))
			}
		})
	}
}

func TestItems_RemoveAt(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		i    Items
		args args
		want *Item
		wantItems Items
	}{
		{
			name: "remove_index_1",
			i:    Items{
				{Value: 1},
				{Value:2},
				{Value:3},
			},
			args: args{
				index: 1,
			},
			want: &Item{
				Value: 2,
			},
			wantItems:Items{
				{Value: 1},
				{Value:3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.RemoveAt(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveAt() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.i, tt.wantItems) {
				t.Errorf("remain Items should be= %v,but remain %v", tt.i, tt.wantItems)
			}
		})
	}
}


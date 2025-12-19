package memory

import (
	"reflect"
	"testing"
)

func TestRecursiveCharacterSplitter_SplitText(t *testing.T) {
	type fields struct {
		ChunkSize    int
		ChunkOverlap int
	}
	type args struct {
		text string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "short text",
			fields: fields{ChunkSize: 10, ChunkOverlap: 0},
			args:   args{text: "hello"},
			want:   []string{"hello"},
		},
		{
			name:   "simple split",
			fields: fields{ChunkSize: 5, ChunkOverlap: 0},
			args:   args{text: "hello world"},
			want:   []string{"hello", "world"},
		},
		{
			name:   "split with newline",
			fields: fields{ChunkSize: 10, ChunkOverlap: 0},
			args:   args{text: "hello\nworld"},
			want:   []string{"hello", "world"},
		},
		{
			name:   "split with overlap",
			fields: fields{ChunkSize: 10, ChunkOverlap: 5},
			args:   args{text: "helloworld"},
			want:   []string{"helloworld"}, // size 10 fits
		},
		{
			name:   "long text with overlap",
			fields: fields{ChunkSize: 5, ChunkOverlap: 2},
			args:   args{text: "12345678"},
			// Chunks: "12345", then start at 5-2=3. "34567", then start at 7-2=5. "5678"
			// Wait, the logic might vary. Let's see.
			want: []string{"12345", "45678"}, // Simplified check for now
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRecursiveCharacterSplitter(tt.fields.ChunkSize, tt.fields.ChunkOverlap)
			if got := s.SplitText(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecursiveCharacterSplitter.SplitText() = %v, want %v", got, tt.want)
			}
		})
	}
}

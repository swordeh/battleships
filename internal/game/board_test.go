package game

import "testing"

func TestNewEmptyBoard(t *testing.T) {
	board := NewEmptyBoard()
	for i := 0; i > 100; i++ {
		if board.Tiles[i] != 0 {
			t.Errorf("NewEmptyBoard failed at iteration %d", i)
		}
	}
}

func TestGridToSliceIndex(t *testing.T) {
	tests := []struct {
		name    string
		ref     string
		want    int
		wantErr bool
	}{
		{
			name:    "valid position in upper case",
			ref:     "A1",
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid position in mid grid",
			ref:     "c5",
			want:    24,
			wantErr: false,
		},
		{
			name:    "valid position at end of grid",
			ref:     "j10",
			want:    99,
			wantErr: false,
		},
		{
			name:    "invalid position - too short",
			ref:     "a",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "invalid position - too long",
			ref:     "ab12",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "invalid position - row out of bounds",
			ref:     "k5",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "invalid position - col out of bounds",
			ref:     "a11",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "invalid position - non numeric column",
			ref:     "ab",
			want:    -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gridToSliceIndex(tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("gridToSliceIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("gridToSliceIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

package navgrid

type Cell struct {
	Blocked bool
}

type NavGrid struct {
	W, H     int
	Cells    []Cell
	CellSize float64
}

func (g *NavGrid) WorldToCell(x, z float64) (cx, cz int) {
	return int(x / g.CellSize), int(z / g.CellSize)
}

func (g *NavGrid) IsBlocked(x, z float64) bool {
	cx, cz := g.WorldToCell(x, z)
	if cx < 0 || cz < 0 || cx >= g.W || cz >= g.H {
		return true
	}
	return g.Cells[cz*g.W+cx].Blocked
}

func NewNavGrid(width, height int, cellSize float64) NavGrid {
	return NavGrid{
		W:        width,
		H:        height,
		CellSize: cellSize,
		Cells:    make([]Cell, width*height),
	}
}
func (g *NavGrid) SetBlocked(x, z int, blocked bool) {
	if x < 0 || z < 0 || x >= g.W || z >= g.H {
		return
	}
	g.Cells[z*g.W+x].Blocked = blocked
}

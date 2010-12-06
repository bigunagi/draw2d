package draw2d

type DashVertexConverter struct {
	command        VertexCommand
	next           VertexConverter
	x, y, distance float
	dash           []float
	currentDash    int
	dashOffset     float
}

func NewDashConverter(dash []float, dashOffset float, converter VertexConverter) *DashVertexConverter {
	var dasher DashVertexConverter
	dasher.dash = dash
	dasher.currentDash = 0
	dasher.dashOffset = dashOffset
	dasher.next = converter
	return &dasher
}

func (dasher *DashVertexConverter) NextCommand(cmd VertexCommand) {
	dasher.command = cmd
	if dasher.command == VertexStopCommand {
		dasher.next.NextCommand(VertexStopCommand)
	}
}

func (dasher *DashVertexConverter) Vertex(x, y float) {
	switch dasher.command {
	case VertexStartCommand:
		dasher.start(x, y)
	default:
		dasher.lineTo(x, y)
	}
	dasher.command = VertexNoCommand
}

func (dasher *DashVertexConverter) start(x, y float) {
	dasher.next.NextCommand(VertexStartCommand)
	dasher.next.Vertex(x, y)
	dasher.x, dasher.y = x, y
	dasher.distance = dasher.dashOffset
	dasher.currentDash = 0
}

func (dasher *DashVertexConverter) lineTo(x, y float) {
	rest := dasher.dash[dasher.currentDash] - dasher.distance
	for rest < 0 {
		dasher.distance = dasher.distance - dasher.dash[dasher.currentDash]
		dasher.currentDash = (dasher.currentDash + 1) % len(dasher.dash)
		rest = dasher.dash[dasher.currentDash] - dasher.distance
	}
	d := distance(dasher.x, dasher.y, x, y)
	for d >= rest {
		k := rest / d
		lx := dasher.x + k*(x-dasher.x)
		ly := dasher.y + k*(y-dasher.y)
		if dasher.currentDash%2 == 0 {
			// line
			dasher.next.Vertex(lx, ly)
		} else {
			// gap
			dasher.next.NextCommand(VertexStopCommand)
			dasher.next.NextCommand(VertexStartCommand)
			dasher.next.Vertex(lx, ly)
		}
		d = d - rest
		dasher.x, dasher.y = lx, ly
		dasher.currentDash = (dasher.currentDash + 1) % len(dasher.dash)
		rest = dasher.dash[dasher.currentDash]
	}
	dasher.distance = d
	if dasher.currentDash%2 == 0 {
		// line
		dasher.next.Vertex(x, y)
	} else {
		// gap
		dasher.next.NextCommand(VertexStopCommand)
		dasher.next.NextCommand(VertexStartCommand)
		dasher.next.Vertex(x, y)
	}
	if dasher.distance >= dasher.dash[dasher.currentDash] {
		dasher.distance = dasher.distance - dasher.dash[dasher.currentDash]
		dasher.currentDash = (dasher.currentDash + 1) % len(dasher.dash)
	}
	dasher.x, dasher.y = x, y
}

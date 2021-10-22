package internal

import (
	"image"
	icolor "image/color"
	"image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

var (
	_diffs = []diff1{
		{
			diff0: diff0{
				key:  "app.json",
				mode: diffMode_MODIFY,
			},
			absFilepath: "/var/apollo-synchronizer/app",
		},
		{
			diff0: diff0{
				key:  "app2.json",
				mode: diffMode_CREATE,
			},
			absFilepath: "/var/apollo-synchronizer/app",
		},
		{
			diff0: diff0{
				key:  "app3.json",
				mode: diffMode_DELETE,
			},
			absFilepath: "/var/apollo-synchronizer/app",
		},
	}

	_results = []*synchronizeResult{
		{
			key:       "app.json",
			mode:      diffMode_MODIFY,
			error:     "you failed",
			succeeded: false,
			published: false,
		},
		{
			key:       "app2.json",
			mode:      diffMode_CREATE,
			error:     "",
			succeeded: true,
			published: false,
		},
		{
			key:       "app3.json",
			mode:      diffMode_DELETE,
			error:     "",
			succeeded: true,
			published: true,
		},
	}
)

func Test_renderer_terminal(t *testing.T) {
	r := terminalRenderer{}

	r.renderingDiffs(_diffs)
	r.renderingResult(_results)
}

func Test_renderer_termui(t *testing.T) {
	r := newTermUI(&SynchronizeScope{
		Mode:              SynchronizeMode_DOWNLOAD,
		Path:              "",
		LocalFiles:        nil,
		ApolloSecret:      "",
		ApolloAppID:       "",
		ApolloEnv:         "",
		ApolloClusterName: "",
		ApolloPortalAddr:  "",
		ApolloAccount:     "",
		ApolloAutoPublish: false,
		Overwrite:         false,
		Force:             false,
	})

	r.renderingDiffs(_diffs)
	r.renderingResult(_results)
}

// draw arrow direction (left arrow, right arrow) image
func directionImage(x1, y1, x2, y2 int, mode SynchronizeMode) image.Image {
	img := image.NewGray(image.Rect(x1, y1, x2, y2))
	padding := 1
	w, h := x2-x1-2*padding, y2-y1-2*padding

	p1 := image.Point{X: x1 + padding, Y: y1 + padding + h/2}
	p2 := image.Point{X: x1 + padding + w/3, Y: y1 + padding}
	p3 := image.Point{X: x1 + padding + w/3, Y: y1 + padding + h/3}
	p4 := image.Point{X: x1 + padding + w, Y: y1 + padding + h/3}
	p5 := image.Point{X: x1 + padding + w, Y: y1 + padding + 2*h/3}
	p6 := image.Point{X: x1 + padding + w/3, Y: y1 + padding + 2*h/3}
	p7 := image.Point{X: x1 + padding + w/3, Y: y1 + padding + h}

	c := icolor.Gray{Y: 255}
	brush := func(x, y int) {
		img.SetGray(x, y, c)
	}

	drawPoly([]image.Point{p1, p2, p3, p4, p5, p6, p7}, brush)

	return img
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func lineTo(x0, y0, x1, y1 int, brush func(x, y int)) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		brush(x0, y0)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func drawPoly(pts []image.Point, brush func(x, y int)) {
	last := len(pts) - 1
	for i := 0; i < len(pts); i++ {
		lineTo(pts[last].X, pts[last].Y, pts[i].X, pts[i].Y, brush)
		last = i
	}
}

func Test_directionImages_down(t *testing.T) {
	img := directionImage(0, 0, 27, 11, SynchronizeMode_DOWNLOAD)
	fd, err := os.OpenFile("./direction_down.jpg", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	require.NoError(t, err)
	defer fd.Close()
	err = jpeg.Encode(fd, img, nil)
	assert.NoError(t, err)
}

func Test_directionImages_up(t *testing.T) {
	img := directionImage(0, 0, 100, 100, SynchronizeMode_UPLOAD)
	fd, err := os.OpenFile("./direction_up.jpg", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	require.NoError(t, err)
	defer fd.Close()
	err = jpeg.Encode(fd, img, nil)
	assert.NoError(t, err)
}

//
//const GOPHER_IMAGE = `iVBORw0KGgoAAAANSUhEUgAAAEsAAAA8CAAAAAALAhhPAAAFfUlEQVRYw62XeWwUVRzHf2+OPbo9d7tsWyiyaZti6eWGAhISoIGKECEKCAiJJkYTiUgTMYSIosYYBBIUIxoSPIINEBDi2VhwkQrVsj1ESgu9doHWdrul7ba73WNm3vOPtsseM9MdwvvrzTs+8/t95ze/33sI5BqiabU6m9En8oNjduLnAEDLUsQXFF8tQ5oxK3vmnNmDSMtrncks9Hhtt/qeWZapHb1ha3UqYSWVl2ZmpWgaXMXGohQAvmeop3bjTRtv6SgaK/Pb9/bFzUrYslbFAmHPp+3WhAYdr+7GN/YnpN46Opv55VDsJkoEpMrY/vO2BIYQ6LLvm0ThY3MzDzzeSJeeWNyTkgnIE5ePKsvKlcg/0T9QMzXalwXMlj54z4c0rh/mzEfr+FgWEz2w6uk8dkzFAgcARAgNp1ZYef8bH2AgvuStbc2/i6CiWGj98y2tw2l4FAXKkQBIf+exyRnteY83LfEwDQAYCoK+P6bxkZm/0966LxcAAILHB56kgD95PPxltuYcMtFTWw/FKkY/6Opf3GGd9ZF+Qp6mzJxzuRSractOmJrH1u8XTvWFHINNkLQLMR+XHXvfPPHw967raE1xxwtA36IMRfkAAG29/7mLuQcb2WOnsJReZGfpiHsSBX81cvMKywYZHhX5hFPtOqPGWZCXnhWGAu6lX91ElKXSalcLXu3UaOXVay57ZSe5f6Gpx7J2MXAsi7EqSp09b/MirKSyJfnfEEgeDjl8FgDAfvewP03zZ+AJ0m9aFRM8eEHBDRKjfcreDXnZdQuAxXpT2NRJ7xl3UkLBhuVGU16gZiGOgZmrSbRdqkILuL/yYoSXHHkl9KXgqNu3PB8oRg0geC5vFmLjad6mUyTKLmF3OtraWDIfACyXqmephaDABawfpi6tqqBZytfQMqOz6S09iWXhktrRaB8Xz4Yi/8gyABDm5NVe6qq/3VzPrcjELWrebVuyY2T7ar4zQyybUCtsQ5Es1FGaZVrRVQwAgHGW2ZCRZshI5bGQi7HesyE972pOSeMM0dSktlzxRdrlqb3Osa6CCS8IJoQQQgBAbTAa5l5epO34rJszibJI8rxLfGzcp1dRosutGeb2VDNgqYrwTiPNsLxXiPi3dz7LiS1WBRBDBOnqEjyy3aQb+/bLiJzz9dIkscVBBLxMfSEac7kO4Fpkngi0ruNBeSOal+u8jgOuqPz12nryMLCniEjtOOOmpt+KEIqsEdocJjYXwrh9OZqWJQyPCTo67LNS/TdxLAv6R5ZNK9npEjbYdT33gRo4o5oTqR34R+OmaSzDBWsAIPhuRcgyoteNi9gF0KzNYWVItPf2TLoXEg+7isNC7uJkgo1iQWOfRSP9NR11RtbZZ3OMG/VhL6jvx+J1m87+RCfJChAtEBQkSBX2PnSiihc/Twh3j0h7qdYQAoRVsRGmq7HU2QRbaxVGa1D6nIOqaIWRjyRZpHMQKWKpZM5feA+lzC4ZFultV8S6T0mzQGhQohi5I8iw+CsqBSxhFMuwyLgSwbghGb0AiIKkSDmGZVmJSiKihsiyOAUs70UkywooYP0bii9GdH4sfr1UNysd3fUyLLMQN+rsmo3grHl9VNJHbbwxoa47Vw5gupIqrZcjPh9R4Nye3nRDk199V+aetmvVtDRE8/+cbgAAgMIWGb3UA0MGLE9SCbWX670TDy1y98c3D27eppUjsZ6fql3jcd5rUe7+ZIlLNQny3Rd+E5Tct3WVhTM5RBCEdiEK0b6B+/ca2gYU393nFj/n1AygRQxPIUA043M42u85+z2SnssKrPl8Mx76NL3E6eXc3be7OD+H4WHbJkKI8AU8irbITQjZ+0hQcPEgId/Fn/pl9crKH02+5o2b9T/eMx7pKoskYgAAAABJRU5ErkJggg==`
//
//func image2() image.Image {
//	img, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(GOPHER_IMAGE)))
//	if err != nil {
//		log.Fatalf("failed to decode gopher image: %v", err)
//	}
//
//	return img
//}

//
//func Test_directionImages_gopher(t *testing.T) {
//	img := image2()
//	fd, err := os.OpenFile("./gopher.jpg", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
//	require.NoError(t, err)
//	defer fd.Close()
//	err = jpeg.Encode(fd, img, nil)
//	assert.NoError(t, err)
//}

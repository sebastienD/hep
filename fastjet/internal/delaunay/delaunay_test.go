// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"math"
	"math/rand"
	"testing"
)

const tol = 1e-3

func TestHierarchicalDelaunayInsertSmall(t *testing.T) {
	// NewPoint(x, y)
	p1 := NewPoint(0, 0)
	p2 := NewPoint(0, 2)
	p3 := NewPoint(1, 0)
	p4 := NewPoint(4, 4)
	ps := []*Point{
		p1,
		p2,
		p3,
		p4,
	}
	d := HierarchicalDelaunay()
	for _, p := range ps {
		d.Insert(p)
	}
	exp := []*Triangle{
		NewTriangle(p1, p2, p3),
		NewTriangle(p2, p3, p4),
	}
	ts := d.Triangles()
	got, want := len(ts), len(exp)
	if got != want {
		t.Errorf("got=%d delaunay triangles, want=%d", got, want)
	}
	for i := range ts {
		ok := false
		for j := range exp {
			if ts[i].Equals(exp[j]) {
				ok = true
				// remove triangles that have been matched from slice,
				// in case there are duplicate triangles.
				exp = append(exp[:j], exp[j+1:]...)
				break
			}
		}
		if !ok {
			t.Errorf("Triangle T%s not as expected", ts[i])
		}
	}
	var (
		nn []*Point
		nd []float64
	)
	for _, p := range ps {
		n, d := p.NearestNeighbor()
		nn = append(nn, n)
		nd = append(nd, d)
	}
	expN := []*Point{p3, p1, p1, p2}
	expD := []float64{1.0, 2.0, 1.0, 4.4721}
	got, want = len(nn), len(expN)
	if got != want {
		t.Errorf("got=%d nearest neighbors, want=%d", got, want)
	}
	for i := range nn {
		if !nn[i].Equals(expN[i]) {
			t.Errorf("got=N%s nearest neighbor, want=N%s", nn[i], expN[i])
		}
		if math.Abs(nd[i]-expD[i]) > tol {
			t.Errorf("got=%f distance, want=%f for point P%s with neighbour N%s", nd[i], expD[i], ps[i], nn[i])
		}
	}
}

func TestHierarchicalDelaunayInsertMedium(t *testing.T) {
	// NewPoint(x, y)
	p1 := NewPoint(-1.5, 3.2)
	p2 := NewPoint(1.8, 3.3)
	p3 := NewPoint(-3.7, 1.5)
	p4 := NewPoint(-1.5, 1.3)
	p5 := NewPoint(0.8, 1.2)
	p6 := NewPoint(3.3, 1.5)
	p7 := NewPoint(-4, -1)
	p8 := NewPoint(-2.3, -0.7)
	p9 := NewPoint(0, -0.5)
	p10 := NewPoint(2, -1.5)
	p11 := NewPoint(3.7, -0.8)
	p12 := NewPoint(-3.5, -2.9)
	p13 := NewPoint(-0.9, -3.9)
	p14 := NewPoint(2, -3.5)
	p15 := NewPoint(3.5, -2.25)
	ps := []*Point{p1, p2, p3, p4, p5, p6, p7, p8,
		p9, p10, p11, p12, p13, p14, p15}
	d := HierarchicalDelaunay()
	for _, p := range ps {
		d.Insert(p)
	}
	ts := d.Triangles()
	exp := []*Triangle{
		NewTriangle(p1, p3, p4),
		NewTriangle(p1, p4, p5),
		NewTriangle(p1, p5, p2),
		NewTriangle(p2, p5, p6),
		NewTriangle(p3, p4, p8),
		NewTriangle(p3, p8, p7),
		NewTriangle(p4, p8, p9),
		NewTriangle(p4, p9, p5),
		NewTriangle(p5, p9, p10),
		NewTriangle(p5, p10, p6),
		NewTriangle(p6, p10, p11),
		NewTriangle(p7, p8, p12),
		NewTriangle(p8, p12, p13),
		NewTriangle(p8, p13, p9),
		NewTriangle(p9, p13, p10),
		NewTriangle(p10, p13, p14),
		NewTriangle(p10, p14, p15),
		NewTriangle(p10, p15, p11),
	}
	got, want := len(ts), len(exp)
	if got != want {
		t.Errorf("got=%d delaunay triangles, want=%d", got, want)
	}
	for i := range ts {
		ok := false
		for j := range exp {
			if ts[i].Equals(exp[j]) {
				ok = true
				// remove triangles that have been matched from slice,
				// in case there are duplicate triangles.
				exp = append(exp[:j], exp[j+1:]...)
				break
			}
		}
		if !ok {
			t.Errorf("Triangle T%s not as expected", ts[i])
		}
	}
	var (
		nn []*Point
		nd []float64
	)
	for _, p := range ps {
		n, d := p.NearestNeighbor()
		nn = append(nn, n)
		nd = append(nd, d)
	}
	expN := []*Point{p4, p5, p4, p1, p9, p11, p8, p7, p5, p15, p15, p7, p12, p15, p11}
	expD := []float64{1.9, 2.326, 2.209, 1.9, 1.879, 2.335, 1.726, 1.726, 1.879, 1.677, 1.464, 1.965, 2.786, 1.953, 1.464}
	got, want = len(nn), len(expN)
	if got != want {
		t.Errorf("got=%d nearest neighbors, want=%d", got, want)
	}
	for i := range nn {
		if !nn[i].Equals(expN[i]) {
			t.Errorf("got=N%s nearest neighbor, want=N%s", nn[i], expN[i])
		}
		if math.Abs(nd[i]-expD[i]) > tol {
			t.Errorf("got=%f distance, want=%f for point P%s with neighbour N%s", nd[i], expD[i], ps[i], nn[i])
		}
	}
}

func grid(nx, ny int, angle float64) []*Point {
	s := math.Sin(angle)
	c := math.Cos(angle)
	var points []*Point
	for xi := 0; xi < nx; xi++ {
		tx := float64(xi)
		for yi := 0; yi < ny; yi++ {
			ty := float64(yi)
			x := tx*c - ty*s
			y := tx*s + ty*c
			points = append(points, NewPoint(x, y))
		}
	}
	return points
}

func TestHierarchicalDelaunayGrid(t *testing.T) {
	const degrees = math.Pi / 180
	const n = 10
	ps := grid(n, n, 10*degrees)
	d := HierarchicalDelaunay()
	for _, p := range ps {
		d.Insert(p)
	}
}

func TestHierarchicalDelaunayGridRotated(t *testing.T) {
	const degrees = math.Pi / 180
	const n = 10
	ps := grid(n, n, 60*degrees)
	d := HierarchicalDelaunay()
	for _, p := range ps {
		d.Insert(p)
	}
}

func benchmarkHierarchicalDelaunayInsertion(i int, b *testing.B) {
	ps := make([]*Point, i)
	for j := 0; j < i; j++ {
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		ps[j] = NewPoint(x, y)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d := HierarchicalDelaunay()
		for _, p := range ps {
			d.Insert(p)
		}
	}
}

func BenchmarkHierarchicalDelaunayInsertion50(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(50, b)
}

func BenchmarkHierarchicalDelaunayInsertion100(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(100, b)
}

func BenchmarkHierarchicalDelaunayInsertion150(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(150, b)
}

func BenchmarkHierarchicalDelaunayInsertion200(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(200, b)
}

func BenchmarkHierarchicalDelaunayInsertion250(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(250, b)
}

func BenchmarkHierarchicalDelaunayInsertion300(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(300, b)
}

func BenchmarkHierarchicalDelaunayInsertion350(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(350, b)
}

func BenchmarkHierarchicalDelaunayInsertion400(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(400, b)
}

func BenchmarkHierarchicalDelaunayInsertion450(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(450, b)
}

func BenchmarkHierarchicalDelaunayInsertion500(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(500, b)
}

func BenchmarkHierarchicalDelaunayInsertion550(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(550, b)
}

func BenchmarkHierarchicalDelaunayInsertion600(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(600, b)
}

func BenchmarkHierarchicalDelaunayInsertion650(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(650, b)
}

func BenchmarkHierarchicalDelaunayInsertion700(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(700, b)
}

func BenchmarkHierarchicalDelaunayInsertion750(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(750, b)
}

func BenchmarkHierarchicalDelaunayInsertion800(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(800, b)
}

func BenchmarkHierarchicalDelaunayInsertion850(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(850, b)
}

func BenchmarkHierarchicalDelaunayInsertion900(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(900, b)
}

func BenchmarkHierarchicalDelaunayInsertion950(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(950, b)
}

func BenchmarkHierarchicalDelaunayInsertion1000(b *testing.B) {
	benchmarkHierarchicalDelaunayInsertion(1000, b)
}

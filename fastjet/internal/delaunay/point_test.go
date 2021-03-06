// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestPointEquals(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		want           bool
	}{
		{5, 5, 5, 5, true},
		{0.5, 4, 1, 4, false},
		{-3, 20, 20, -3, false},
		{2.5, 4, 2.5, 4, true},
	}
	for _, test := range tests {
		p1 := NewPoint(test.x1, test.y1)
		p2 := NewPoint(test.x2, test.y2)
		got := p1.Equals(p2)
		if got != test.want {
			t.Errorf("%v == %v, got = %v, want = %v", p1, p2, got, test.want)
		}
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		want           float64
	}{
		{1, 1, 3, 3, 8},
		{0, 0, 0, 5, 25},
	}
	for _, test := range tests {
		p1 := NewPoint(test.x1, test.y1)
		p2 := NewPoint(test.x2, test.y2)
		got := p1.distance(p2)
		if got != test.want {
			t.Errorf("Distance between %v and %v, got = %v, want = %v", p1, p2, got, test.want)
		}
	}
}

func TestFindNearest(t *testing.T) {
	const tol float64 = 0.001
	tests := []struct {
		p                 *Point
		adjacentTriangles []*Triangle
		wantDist          float64
		wantNeighbor      *Point
	}{
		{
			NewPoint(0, 0),
			[]*Triangle{
				NewTriangle(NewPoint(0, 0), NewPoint(2, 0), NewPoint(0, 5)),
				NewTriangle(NewPoint(3, 3), NewPoint(0, 0), NewPoint(0, 5)),
			},
			2,
			NewPoint(2, 0),
		},
	}
	for _, test := range tests {
		test.p.adjacentTriangles = test.adjacentTriangles
		test.p.findNearest()
		gotNeighbor, gotDistance := test.p.NearestNeighbor()
		if !floats.EqualWithinAbs(gotDistance, test.wantDist, tol) || !test.wantNeighbor.Equals(gotNeighbor) {
			t.Errorf("Nearest Neighbor for P%v with adjacent triangles %v, \n gotDist = %v, wantDist = %v, gotNeighbor = %v, wantNeighbor = %v",
				test.p, test.adjacentTriangles, gotDistance, test.wantDist, gotNeighbor, test.wantNeighbor)
		}
	}
}

func TestInTriangle(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, x3, y3, x, y float64
		want                         location
	}{
		{2, 2, 5, 3, 6, 2, 5, 2.5, inside},
		{2, 2, 5, 3, 6, 2, 5, 2, onEdge},
		{2, 2, 5, 3, 6, 2, 4, 3, outside},
		{2, 2, 5, 3, 6, 2, 10, 2, outside},
		{5, 3, 2, 2, 6, 1.5, 5, 2.5, inside},
		{5, 3, 2, 2, 6, 1.5, 6, 2.5, outside},
	}
	for _, test := range tests {
		tri := NewTriangle(NewPoint(test.x1, test.y1), NewPoint(test.x2, test.y2), NewPoint(test.x3, test.y3))
		p := NewPoint(test.x, test.y)
		got := p.inTriangle(tri)
		if got != test.want {
			t.Fatalf("Point %v in Triangle %v, got = %v, want = %v", p, tri, got, test.want)
		}
	}
}

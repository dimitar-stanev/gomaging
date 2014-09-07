package main

import { 
	"fmt"
	"testing"
	"image/color"
}

func BenchmarkKmeans(b *testing.B) {
	var c *Canvas = CanvasFromFile("forest.png")
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("One iteration")
        c = kMeans(30, c)
    }
}

type testpairforcompare struct {
	firstArray []color.Color
	secondArray []color.Color
	isItTrue bool
}

var tests []testpairforcompare {
	{
		{ color.RGBA(20, 20, 20, 255), color.RGBA(20, 20, 20, 255), color.RGBA(20, 20, 20, 255) }, 
		{ color.RGBA(20, 20, 20, 255), color.RGBA(20, 20, 20, 255), color.RGBA(20, 20, 20, 255), 
			true }
	},
	{
		{ color.RGBA(20, 20, 30, 255), color.RGBA(20, 10, 20, 255), color.RGBA(15, 20, 20, 255) }, 
		{ color.RGBA(20, 20, 20, 255), color.RGBA(20, 20, 20, 255), color.RGBA(20, 20, 20, 255), 
			false }
	}
}

func TestCompare( t *testing.T ) {
	for _, pair := range tests {
	        v := compareCentroids(pair.firstArray, pair.secondArray)
	        if v != pair.isItTrue {
	            t.Error(
	                "For", pair.firstArray,
	                "and", pair.secondArray 
	                "expected", pair.isItTrue,
	                "got", v,
	            )
	        }
	    }
	
}
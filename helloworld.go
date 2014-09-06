package main

import "fmt"
import "image"
import "image/png"
import "os"
import "log"
import "image/color"
import "math/rand"
import "time"
import "math"

type Point struct {
	x, y int
}

type PixelsWithCentroid struct {
	centroidIndex int
	colour color.Color
}

func main() {

	rand.Seed( time.Now().UnixNano())
	var c *Canvas = CanvasFromFile("forest.png")
	var v1 = Vector{0, 0}
	var v2 = Vector{50,50}
	c.Blur(10, new(WeightFunctionDouble))
	c.DrawLine(color.RGBA{255,255,255,255}, v1, v2)
	c.DrawCircle(color.RGBA{255,255,255,255}, Vector{80,80}, 30)
	width := c.Bounds().Max.X
	height := c.Bounds().Max.Y
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	// var newCanvas = kMeans(100, c)
	size := m.Bounds().Size()
	for xPos:= 0; xPos < size.X; xPos++ {
		for yPos := 0; yPos < size.Y; yPos++ {
			m.Set(xPos, yPos, c.At(xPos, yPos))
		}
	}
	out_filename := "result3.png"
	out_file, err := os.Create(out_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out_file.Close()
	log.Print("Saving image to: ", out_filename)
	png.Encode(out_file, m)
}
func kMeans(k int, canv *Canvas) *Canvas {
  	// Assign the cluster points randomly
	var clusterPoints = make([]Point, k);
	for i:=0; i<k; i++ {
		clusterPoints[i] = Point{
			rand.Intn(canv.Bounds().Max.X),
			rand.Intn(canv.Bounds().Max.Y),
		}
	}

  	// Make the colors of the random centroids
	var centroidPoints = make([]color.Color, k)
	for j:=0; j<k; j++ {
		centroidPoints[j] = canv.At(clusterPoints[j].x, clusterPoints[j].y )
	}

  	// fmt.Println("CENTROIDS : ", centroidPoints[0])

	var dividedImage = make([][]PixelsWithCentroid, canv.Bounds().Max.X)
	for i := range dividedImage {
		dividedImage[i] = make([]PixelsWithCentroid, canv.Bounds().Max.Y)
	}

	dividedImage = calculateNearestCentroid(centroidPoints, canv)

	var means = make([]color.Color, k)

	means = calculateNewCentroid(dividedImage, k, canv)

	var lastMeans = make([]color.Color, k)
	lastMeans = means

	var times = 0
	for tries:=0; tries<100; tries++ {
  		// fmt.Println(times)
		dividedImage = calculateNearestCentroid(means, canv)
		means = calculateNewCentroid(dividedImage,k, canv)
  		 //  	fmt.Println(means)
  			// fmt.Println(lastMeans)
		if ( compareCentroids(means, lastMeans) ) {
  			// fmt.Println(means)
  			// fmt.Println(lastMeans)
			fmt.Println("SUCCESS BITCH !")
			fmt.Println("Took exactly ", times, "times to get it right.")
			break
		}
		lastMeans = means
		times++
	}
	for i:=0; i<k; i++ {
		fmt.Println(means[i])
	}

	canv = convertImage(dividedImage, means, canv)
	fmt.Println(canv.At(50,50), canv.At(0,0))
	return canv;
}

func lengthBetweenPixels(p1 color.Color, p2 color.Color) float64 {

	var r1, g1, b1, _ = p1.RGBA()
	var r2, g2, b2, _ = p2.RGBA()

	r1f := float64(r1)/257
	r2f := float64(r2)/257
	g1f := float64(g1)/257
	g2f := float64(g2)/257
	b1f := float64(b1)/257
	b2f := float64(b2)/257

	var a, b, c float64

	if r1f>r2f {
		a = r1f-r2f
	} else {
		a = r2f-r1f
	}
	if g1f>g2f {
		b = g1f-g2f
	} else {
		b = g2f-g1f
	}
	if b1f>b2f {
		c = b1f-b2f
	} else {
		c = b2f-b1f
	}


	return math.Sqrt(math.Pow(a,2)+math.Pow(b,2)+math.Pow(c,2))
}

func minimumLength (arr []float64) (float64,int) {
	var min = arr[0]
	var minIndex = 0
	for i:=0; i<len(arr); i++ {
		if arr[i] < min {
			min = arr[i]
			minIndex = i
		}
	}
	return min, minIndex
}

func calculateNearestCentroid(centroidPoints []color.Color, canvs *Canvas) [][]PixelsWithCentroid {

	var divdImg = make([][]PixelsWithCentroid, canvs.Bounds().Max.X)
	for i := range divdImg {
		divdImg[i] = make([]PixelsWithCentroid, canvs.Bounds().Max.Y)
	}
	var lengthsBetweenCurrentPixelAndCentroids = make([]float64, len(centroidPoints))

  	// calculate the nearest centroid for each point of the image
	for xPosition:=0; xPosition<canvs.Bounds().Max.X; xPosition++ {
		for yPosition:=0; yPosition<canvs.Bounds().Max.Y; yPosition++ {
			for currentCentroid:=0; currentCentroid < len(centroidPoints); currentCentroid++{
				lengthsBetweenCurrentPixelAndCentroids[currentCentroid] = lengthBetweenPixels(canvs.At(xPosition, yPosition), centroidPoints[currentCentroid])
    			if currentCentroid == len(centroidPoints)-1 { // here we have completed the length between a pixel and each centroid, so we check which one has the minimum "length"
    				_, b:= minimumLength(lengthsBetweenCurrentPixelAndCentroids)
    				// fmt.Println("FOR PIXEL : ", xPosition, yPosition, "closest centroid ", b, "is at ", a)
    				var currPix = PixelsWithCentroid{ b, canvs.At(xPosition, yPosition)}
    				divdImg[xPosition][yPosition] = currPix
    			}
    		}
    	}
    }
    return divdImg
}

func calculateNewCentroid(divdImg [][]PixelsWithCentroid, k int, canvs *Canvas) []color.Color {

	// calculate the new centroids and return them as means[]
	var newCentroids = make([]color.Color, k)
	for currCentrCalculated := 0; currCentrCalculated<k; currCentrCalculated++ {
		var centroidCounter uint32 = 0
		var sumR uint32
		var sumG uint32
		var sumB uint32
		for xPosition:=0; xPosition<canvs.Bounds().Max.X; xPosition++ {
			for yPosition:=0; yPosition<canvs.Bounds().Max.Y; yPosition++ {
				if divdImg[xPosition][yPosition].centroidIndex == currCentrCalculated {
					var rToAdd, _, _, _ =  divdImg[xPosition][yPosition].colour.RGBA()	
					sumR += rToAdd
					var _,gToAdd, _, _ =  divdImg[xPosition][yPosition].colour.RGBA()	
					sumG += gToAdd
					var _, _, bToAdd, _ =  divdImg[xPosition][yPosition].colour.RGBA()	
					sumB += bToAdd
					centroidCounter++
				}
			}
		}
    	// fmt.Println("FOR CENTROID ", currCentrCalculated," THERE ARE ", centroidCounter)
		if centroidCounter == 0 {
			centroidCounter = 1
		}
		newCentroids[currCentrCalculated] = color.RGBA{ uint8(sumR/centroidCounter), uint8(sumG/centroidCounter), uint8(sumB/centroidCounter), 255}
	}
	return newCentroids
}

func compareCentroids ( a []color.Color, b []color.Color ) bool {
	for i:=0; i<len(a); i++ {
		if !containsPixel(a[i], b){
			return false
		}
	}
	return true
}
func containsPixel (pxl color.Color, slice []color.Color) bool {
	for i:=0; i<len(slice); i++ {
		if slice[i] == pxl {
			return true
		}
	}
	return false
}
func convertImage(divdImg [][]PixelsWithCentroid, centroids []color.Color, canvs *Canvas ) *Canvas {
	for currCentrCalculated:=0; currCentrCalculated<len(centroids); currCentrCalculated++ {
		for xPosition:=0; xPosition<canvs.Bounds().Max.X; xPosition++ {
			for yPosition:=0; yPosition<canvs.Bounds().Max.Y; yPosition++ {
				if divdImg[xPosition][yPosition].centroidIndex == currCentrCalculated {
					canvs.Set(xPosition, yPosition, centroids[currCentrCalculated])
				}
			}
		}
	}

	return canvs
}
package main

import (
	"fmt"
	"slices"
)

func main() {
    var i int
	var avg float32
	var mySlice [] int;
	fmt.Print("Type some numbers: ('-1' to exit)\n")
	for{
		fmt.Scan(&i)
		if(i == -1){
			break
		}
		mySlice = append(mySlice, i)
	}
	
	avg = average(mySlice)
	var sum int = sum(mySlice)
	fmt.Println("Sum = ", sum, "\nAverage = ", avg )
	
	fmt.Println("Max = ", slices.Max(mySlice), "\nMin = ", slices.Min(mySlice))
}

func average(nums[] int) float32{

	var sum int = 0
	var size int = len(nums)
	// code to be executed
	for i := 0; i < size; i++ {
		sum+=nums[i]
	}

	var avg float32 = float32(sum / size)
	return avg;

  }

  func sum(nums[] int) int{

	var sum int = 0
	var size int = len(nums)

	for i := 0; i < size; i++ {
		sum+=nums[i]
	}
	return sum;

  }
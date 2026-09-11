package quicksort

import (
	"cmp"
)

func Quicksort[S []I, I cmp.Ordered](arr S, low int, high int) {
	if low < high {
		pivot := getPivot(arr, low, high)

		// Sort the low section
		Quicksort(arr, low, pivot-1)
		// Sort the high section
		Quicksort(arr, pivot+1, high)
	}
}

func getPivot[S []I, I cmp.Ordered](arr S, low int, high int) int {
	pivot := arr[high]

	// i starts 1 behind j and keeps track of which index value to swap with j index value when j index value is less than the pivot
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	// Swap the Pivot and the index value after i. Dont really understand why yet but I have it memorized this way
	arr[i+1], arr[high] = arr[high], arr[i+1]

	// return the index of the pivot
	return i + 1
}

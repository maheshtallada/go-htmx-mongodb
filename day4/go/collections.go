package main

// Filter returns only the even numbers, preserving order.
// func Filter(nums []int) []int {
//     int[] filtered = []
//     for (int i=0, i<nums.length, i++) {
//         if (i%2==0) filtered[filtered.length] = i
//     }
// }

func Filter(nums []int) []int {
    filtered := make([]int, 0)

    for _, v := range nums {
        if v%2 == 0 {
            filtered = append(filtered, v)
        }
    }

    return filtered
}

// 1. range loops through each value in nums.
// 2. v%2 == 0 keeps only even numbers.
// 3. append adds those values to a new slice while preserving the original order.
// •
//    for , v := range nums:
//    ◦
//    This means: iterate through every value in the slice nums
//    ◦
//     ignores the index, v stores each number


// make([]int, 0):
// ◦
// Creates an empty slice of type []int
// ◦
// 0 means initial length is 0, but it can grow later with append

// SumRange returns the sum of nums using range.
func SumRange(nums []int) int {
    sum:= 0
    for _, v := range nums {
        sum += v
    }
    return sum
}
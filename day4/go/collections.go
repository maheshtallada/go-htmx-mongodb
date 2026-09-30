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
//    for , v := range nums:
//    This means: iterate through every value in the slice nums
//     ignores the index, v stores each number


// make([]int, 0):
// Creates an empty slice of type []int
// 0 means initial length is 0, but it can grow later with append

// SumRange returns the sum of nums using range.
func SumRange(nums []int) int {
    sum:= 0
    for _, v := range nums {
        sum += v
    }
    return sum
}


/*
Original Java-like attempt (kept for comparison; its syntax and count logic need correction):

Map<int, String> newMap = new HashMap<>();

for (String str : words) {
    if (!newMap.has(str)) {
        newMap.set(str, 1);
    }
    newMap[str] = newMap.get(str) + 1;
}

Corrected Java version:

import java.util.HashMap;
import java.util.Map;

public static Map<String, Integer> wordCount(String[] words) {
    Map<String, Integer> counts = new HashMap<>();

    for (String word : words) {
        counts.put(word, counts.getOrDefault(word, 0) + 1);
    }

    return counts;
}
*/

// WordCount returns how many times each word appears.
func WordCount(words []string) map[string]int {
	// make creates an initialized map that can store string keys and int counts.
	counts := make(map[string]int)

	for _, word := range words {
		// A missing map key reads as 0, so ++ sets it to 1 or adds 1 to its existing count.
		counts[word]++
	}

	return counts
}

// MapSlice applies fn to each element and returns a new slice.
// Signature breakdown:
//   - [T any, U any] declares two type parameters. `any` means either may be any type.
//   - in []T accepts a slice whose elements have input type T.
//   - fn func(T) U accepts a function that takes one T and returns one U.
//   - []U returns a slice of the output type U.
// The caller usually does not need to specify T and U; Go infers them from the arguments.
func MapSlice[T any, U any](in []T, fn func(T) U) []U {
	// Allocate a separate output slice with one slot per input element.
	// Keeping the same length also makes the result order match the input.
	out := make([]U, len(in))

	// range gives each input index and value; apply fn to transform the value.
	for i, value := range in {
		// Store the transformed U value in the corresponding output position.
		out[i] = fn(value)
	}

	// Return the new slice without changing the input slice.
	return out
}

// Keys returns the keys of a map in unspecified order.
func Keys[K comparable, V any](m map[K]V) []K {
	// K must be comparable because Go map keys must support equality checks.
	// V can be any type; this function only reads the keys.
	keys := make([]K, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
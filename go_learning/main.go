package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("-------------只出现一次的数字----------")
	nums := []int{4, 1, 2, 1, 2}
	result := singleNumber(nums)
	println(result) // Output: 4

	fmt.Println("-------------回文数----------")
	fmt.Println(isPalindrome(121))  // true
	fmt.Println(isPalindrome(-121)) // false
	fmt.Println(isPalindrome(10))   // false

	fmt.Println("-------------有效的括号----------")
	fmt.Println(isValid("()"))     // true
	fmt.Println(isValid("()[]{}")) // true
	fmt.Println(isValid("(]"))     // false
	fmt.Println(isValid("([)]"))   // false
	fmt.Println(isValid("{[]}"))   // true

	fmt.Println("-------------最长公共前缀----------")
	strs := []string{"flower", "flow", "flight"}
	fmt.Println(longestCommonPrefix(strs))                              // Output: "fl"
	fmt.Println(longestCommonPrefix([]string{"dog", "racecar", "car"})) // Output: ""

	fmt.Println("-------------加一----------")
	digits := []int{1, 2, 3}
	fmt.Println(plusOne(digits))         // Output: [1, 2, 4]
	fmt.Println(plusOne([]int{9, 9, 9})) // Output: [1, 0, 0, 0]
	fmt.Println("-------------删除排序数组中的重复项----------")

	nums2 := []int{1, 1, 2}
	length := removeDuplicates(nums2)
	fmt.Println(length)         // Output: 2
	fmt.Println(nums2[:length]) // Output: [1, 2]

	fmt.Println("-------------合并区间----------")
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	merged := merge(intervals)
	fmt.Println(merged) // Output: [[1, 6], [8, 10], [15, 18]]

	fmt.Println("-------------两数之和----------")
	nums3 := []int{2, 7, 11, 15}
	target := 9
	twoSumResult := twoSum(nums3, target)
	if twoSumResult != nil {
		fmt.Printf("Indices of the two numbers that add up to %d: %v\n", target, twoSumResult) // Output: [0, 1]
	} else {
		fmt.Println("No two numbers found that add up to the target.")
	}
}

// singleNumber finds the element that appears only once in an array where every other element appears twice.
func singleNumber(nums []int) int {
	result := 0
	for _, num := range nums {
		result ^= num
	}
	return result
}

// 回文数
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	if x < 10 {
		return true
	}
	if x%10 == 0 {
		return false
	}

	reversed := 0
	for x > reversed {
		reversed = reversed*10 + x%10
		x /= 10
	}

	return x == reversed || x == reversed/10
}
func reversed(x int) int {
	reversed := 0
	for x > 0 {
		reversed = reversed*10 + x%10
		x /= 10
	}
	return reversed
}

// 有效的括号 (使用栈)
func isValid(s string) bool {
	stack := []rune{}
	mapping := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		if open, exists := mapping[char]; exists {
			if len(stack) == 0 || stack[len(stack)-1] != open {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, char)
		}
	}

	return len(stack) == 0
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, str := range strs[1:] {
		for len(str) < len(prefix) || str[:len(prefix)] != prefix {
			prefix = prefix[:len(prefix)-1]
			if len(prefix) == 0 {
				return ""
			}
		}
	}
	return prefix
}

// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一，并返回一个新的数组表示结果。
func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	return append([]int{1}, digits...)
}

// 删除排序数组中的重复项
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	j := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[j] {
			j++
			nums[j] = nums[i]
		}
	}
	return j + 1
}

// 合并区间
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	// 按照区间的起始位置排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := [][]int{intervals[0]}
	for _, interval := range intervals[1:] {
		last := merged[len(merged)-1]
		if interval[0] <= last[1] { // 有重叠
			last[1] = max(last[1], interval[1])
		} else {
			merged = append(merged, interval)
		}
	}
	return merged
}

// 两数之和，给定一个整数数组 nums 和一个目标值 target，找出 nums 中的两个数，使它们的和等于 target，并返回它们的数组下标。
func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		if j, found := numMap[target-num]; found {
			return []int{j, i}
		}
		numMap[num] = i
	}
	return nil // 如果没有找到，返回 nil
}

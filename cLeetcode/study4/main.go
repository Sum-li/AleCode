package main

func main() {
	println(findMedianSortedArrays1([]int{1, 3, 4}, []int{2}))
	println(findMedianSortedArrays2([]int{1, 3, 4}, []int{2}))
}
func findMedianSortedArrays1(nums1 []int, nums2 []int) float64 {
	index1, index2, result := 0, 0, 0.0
	arr := []int{}
	point, remainder := (len(nums1)+len(nums2)+1)/2, (len(nums1)+len(nums2)+1)%2
	for i := 0; i < point+1; i++ {
		if index1 >= len(nums1) {
			arr = append(arr, nums2[index2:]...)
			break
		}
		if index2 >= len(nums2) {
			arr = append(arr, nums1[index1:]...)
			break
		}
		if nums1[index1] < nums2[index2] {
			arr = append(arr, nums1[index1])
			index1++
		} else {
			arr = append(arr, nums2[index2])
			index2++
		}
	}
	result = float64(arr[point-1]+arr[point-1+remainder]) / 2
	return result
}

func findMedianSortedArrays2(nums1 []int, nums2 []int) float64 {
	index1, index2, result := 0, 0, 0.0
	arr := make([]int, len(nums1)+len(nums2))
	point, remainder := (len(nums1)+len(nums2)+1)/2, (len(nums1)+len(nums2)+1)%2
	for i := 0; i < point+1; i++ {
		if index1 >= len(nums1) {
			arr = append(arr[:i], nums2[index2:index2+point-i+1]...)
			break
		}
		if index2 >= len(nums2) {
			arr = append(arr[:i], nums1[index1:index1+point-i+1]...)
			break
		}
		if nums1[index1] < nums2[index2] {
			arr[i] = nums1[index1]
			index1++
		} else {
			arr[i] = nums2[index2]
			index2++
		}
	}
	result = float64(arr[point-1]+arr[point-1+remainder]) / 2
	return result
}

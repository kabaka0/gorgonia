package tensor

import (
	"github.com/stretchr/testify/assert"
	//"runtime"
	//"fmt"
	"testing"
)

func TestMaskedInspection(t *testing.T) {
	assert := assert.New(t)

	var retT *Dense

	//vector case
	T := New(Of(Bool), WithShape(1, 12))
	T.ResetMask(false)
	assert.False(T.MaskedAny().(bool))
	for i := 0; i < 12; i += 2 {
		T.mask[i] = true
	}
	assert.True(T.MaskedAny().(bool))
	assert.True(T.MaskedAny(0).(bool))
	assert.False(T.MaskedAll().(bool))
	assert.False(T.MaskedAll(0).(bool))
	assert.Equal(6, T.MaskedCount())
	assert.Equal(6, T.MaskedCount(0))
	assert.Equal(6, T.NonMaskedCount())
	assert.Equal(6, T.NonMaskedCount(0))

	//contiguous mask case
	/*equivalent python code
	  ---------
	  import numpy.ma as ma
	  a = ma.arange(12).reshape((2, 3, 2))
	  a[0,0,0]=ma.masked
	  a[0,2,0]=ma.masked
	  print(ma.getmask(a).all())
	  print(ma.getmask(a).any())
	  print(ma.count_masked(a))
	  print(ma.count(a))
	  print(ma.getmask(a).all(0))
	  print(ma.getmask(a).any(0))
	  print(ma.count_masked(a,0))
	  print(ma.count(a,0))
	  print(ma.getmask(a).all(1))
	  print(ma.getmask(a).any(1))
	  print(ma.count_masked(a,1))
	  print(ma.count(a,1))
	  print(ma.getmask(a).all(2))
	  print(ma.getmask(a).any(2))
	  print(ma.count_masked(a,2))
	  print(ma.count(a,2))
	  -----------
	*/
	T = New(Of(Bool), WithShape(2, 3, 2))
	T.ResetMask(false)

	for i := 0; i < 2; i += 2 {
		for j := 0; j < 3; j += 2 {
			for k := 0; k < 2; k += 2 {
				a, b, c := T.maskStrides[0], T.maskStrides[1], T.maskStrides[2]
				T.mask[i*a+b*j+c*k] = true
			}
		}
	}

	assert.Equal([]bool{true, false, false, false, true, false,
		false, false, false, false, false, false}, T.mask)

	assert.Equal(false, T.MaskedAll())
	assert.Equal(true, T.MaskedAny())
	assert.Equal(2, T.MaskedCount())
	assert.Equal(10, T.NonMaskedCount())

	retT = T.MaskedAll(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false, false, false}, retT.bools())
	retT = T.MaskedAny(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, false, false, true, false}, retT.bools())
	retT = T.MaskedCount(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]int{1, 0, 0, 0, 1, 0}, retT.ints())
	retT = T.NonMaskedCount(0).(*Dense)
	assert.Equal([]int{1, 2, 2, 2, 1, 2}, retT.ints())

	retT = T.MaskedAll(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false}, retT.bools())
	retT = T.MaskedAny(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, false, false}, retT.bools())
	retT = T.MaskedCount(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]int{2, 0, 0, 0}, retT.ints())
	retT = T.NonMaskedCount(1).(*Dense)
	assert.Equal([]int{1, 3, 3, 3}, retT.ints())

	retT = T.MaskedAll(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false, false, false}, retT.bools())
	retT = T.MaskedAny(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false, false, false}, retT.bools())
	retT = T.MaskedCount(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]int{1, 0, 1, 0, 0, 0}, retT.ints())
	retT = T.NonMaskedCount(2).(*Dense)
	assert.Equal([]int{1, 2, 1, 2, 2, 2}, retT.ints())

	//non-contiguous mask case 1
	/*equivalent python code
	  ---------
	  import numpy.ma as ma
	  a = ma.arange(12).reshape((2, 3, 2))
	  a[0,0,0]=ma.masked
	  a[0,0,1]=ma.masked
	  a[0,2,0]=ma.masked
	  a[0,2,1]=ma.masked
	  #repeat print() statements
	  -----------
	*/
	T = New(Of(Bool), WithShape(2, 3, 2), WithMaskStrides([]bool{true, true, false}))
	T.ResetMask(false)

	for i := 0; i < 2; i += 2 {
		for j := 0; j < 3; j += 2 {
			for k := 0; k < 2; k += 2 {
				a, b, c := T.maskStrides[0], T.maskStrides[1], T.maskStrides[2]
				T.mask[i*a+b*j+c*k] = true
			}
		}
	}

	assert.Equal([]bool{true, false, true, false, false, false}, T.mask)

	assert.Equal(false, T.MaskedAll())
	assert.Equal(true, T.MaskedAny())
	assert.Equal(4, T.MaskedCount())
	assert.Equal(8, T.NonMaskedCount())

	retT = T.MaskedAll(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false, false, false}, retT.bools())
	retT = T.MaskedAny(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{true, true, false, false, true, true}, retT.bools())
	retT = T.MaskedCount(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]int{1, 1, 0, 0, 1, 1}, retT.ints())
	retT = T.NonMaskedCount(0).(*Dense)
	assert.Equal([]int{1, 1, 2, 2, 1, 1}, retT.ints())

	retT = T.MaskedAll(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false}, retT.bools())
	retT = T.MaskedAny(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{true, true, false, false}, retT.bools())
	retT = T.MaskedCount(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]int{2, 2, 0, 0}, retT.ints())
	retT = T.NonMaskedCount(1).(*Dense)
	assert.Equal([]int{1, 1, 3, 3}, retT.ints())

	retT = T.MaskedAll(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false, false, false}, retT.bools())
	retT = T.MaskedAny(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false, false, false}, retT.bools())
	retT = T.MaskedCount(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]int{2, 0, 2, 0, 0, 0}, retT.ints())
	retT = T.NonMaskedCount(2).(*Dense)
	assert.Equal([]int{0, 2, 0, 2, 2, 2}, retT.ints())

	//non-contiguous mask case 2
	/*equivalent python code
	 ---------
	 import numpy.ma as ma
	 a = ma.arange(12).reshape((2, 3, 2))
	 a[0,0,0]=ma.masked
	 a[0,1,0]=ma.masked
	 a[0,2,0]=ma.masked
	#repeat print() statements
	 -----------
	*/
	T = New(Of(Bool), WithShape(2, 3, 2), WithMaskStrides([]bool{true, false, true}))
	T.ResetMask(false)

	for i := 0; i < 2; i += 2 {
		for j := 0; j < 3; j += 2 {
			for k := 0; k < 2; k += 2 {
				a, b, c := T.maskStrides[0], T.maskStrides[1], T.maskStrides[2]
				T.mask[i*a+b*j+c*k] = true
			}
		}
	}

	assert.Equal([]bool{true, false, false, false}, T.mask)

	assert.Equal(false, T.MaskedAll())
	assert.Equal(true, T.MaskedAny())
	assert.Equal(3, T.MaskedCount())
	assert.Equal(9, T.NonMaskedCount())

	retT = T.MaskedAll(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false, false, false}, retT.bools())
	retT = T.MaskedAny(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false, true, false}, retT.bools())
	retT = T.MaskedCount(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]int{1, 0, 1, 0, 1, 0}, retT.ints())
	retT = T.NonMaskedCount(0).(*Dense)
	assert.Equal([]int{1, 2, 1, 2, 1, 2}, retT.ints())

	retT = T.MaskedAll(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, false, false}, retT.bools())
	retT = T.MaskedAny(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, false, false}, retT.bools())
	retT = T.MaskedCount(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]int{3, 0, 0, 0}, retT.ints())
	retT = T.NonMaskedCount(1).(*Dense)
	assert.Equal([]int{0, 3, 3, 3}, retT.ints())

	retT = T.MaskedAll(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false, false, false}, retT.bools())
	retT = T.MaskedAny(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{true, true, true, false, false, false}, retT.bools())
	retT = T.MaskedCount(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]int{1, 1, 1, 0, 0, 0}, retT.ints())
	retT = T.NonMaskedCount(2).(*Dense)
	assert.Equal([]int{1, 1, 1, 2, 2, 2}, retT.ints())

	//non-contiguous mask case 3
	/*equivalent python code
	  ---------
	  import numpy.ma as ma
	  a = ma.arange(12).reshape((2, 3, 2))
	  a[0,0,0]=ma.masked
	  a[1,0,0]=ma.masked
	  a[0,2,0]=ma.masked
	  a[1,2,0]=ma.masked
	  #repeat print() statements
	  -----------
	*/
	T = New(Of(Bool), WithShape(2, 3, 2), WithMaskStrides([]bool{false, true, true}))
	T.ResetMask(false)

	for i := 0; i < 2; i += 2 {
		for j := 0; j < 3; j += 2 {
			for k := 0; k < 2; k += 2 {
				a, b, c := T.maskStrides[0], T.maskStrides[1], T.maskStrides[2]
				T.mask[i*a+b*j+c*k] = true
			}
		}
	}

	assert.Equal([]bool{true, false, false, false, true, false}, T.mask)

	assert.Equal(false, T.MaskedAll())
	assert.Equal(true, T.MaskedAny())
	assert.Equal(4, T.MaskedCount())
	assert.Equal(8, T.NonMaskedCount())

	retT = T.MaskedAll(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, false, false, true, false}, retT.bools())
	retT = T.MaskedAny(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, false, false, true, false}, retT.bools())
	retT = T.MaskedCount(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]int{2, 0, 0, 0, 2, 0}, retT.ints())
	retT = T.NonMaskedCount(0).(*Dense)
	assert.Equal([]int{0, 2, 2, 2, 0, 2}, retT.ints())

	retT = T.MaskedAll(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false}, retT.bools())
	retT = T.MaskedAny(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false}, retT.bools())
	retT = T.MaskedCount(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]int{2, 0, 2, 0}, retT.ints())
	retT = T.NonMaskedCount(1).(*Dense)
	assert.Equal([]int{1, 3, 1, 3}, retT.ints())

	retT = T.MaskedAll(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false, false, false}, retT.bools())
	retT = T.MaskedAny(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, true, false, true}, retT.bools())
	retT = T.MaskedCount(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]int{1, 0, 1, 1, 0, 1}, retT.ints())
	retT = T.NonMaskedCount(2).(*Dense)
	assert.Equal([]int{1, 2, 1, 1, 2, 1}, retT.ints())

	//non-contiguous mask case 4
	/*equivalent python code
	  ---------
	  import numpy.ma as ma
	  a = ma.arange(12).reshape((2, 3, 2))
	  a[0,0,0]=ma.masked
	  a[0,1,0]=ma.masked
	  a[0,2,0]=ma.masked
	  a[1,0,0]=ma.masked
	  a[1,1,0]=ma.masked
	  a[1,2,0]=ma.masked
	  #repeat print() statements
	  -----------
	*/
	T = New(Of(Bool), WithShape(2, 3, 2), WithMaskStrides([]bool{false, false, true}))
	T.ResetMask(false)

	for i := 0; i < 2; i += 2 {
		for j := 0; j < 3; j += 2 {
			for k := 0; k < 2; k += 2 {
				a, b, c := T.maskStrides[0], T.maskStrides[1], T.maskStrides[2]
				T.mask[i*a+b*j+c*k] = true
			}
		}
	}

	assert.Equal([]bool{true, false}, T.mask)

	assert.Equal(false, T.MaskedAll())
	assert.Equal(true, T.MaskedAny())
	assert.Equal(6, T.MaskedCount())
	assert.Equal(6, T.NonMaskedCount())

	retT = T.MaskedAll(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false, true, false}, retT.bools())
	retT = T.MaskedAny(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false, true, false}, retT.bools())
	retT = T.MaskedCount(0).(*Dense)
	assert.Equal([]int{3, 2}, []int(retT.shape))
	assert.Equal([]int{2, 0, 2, 0, 2, 0}, retT.ints())
	retT = T.NonMaskedCount(0).(*Dense)
	assert.Equal([]int{0, 2, 0, 2, 0, 2}, retT.ints())

	retT = T.MaskedAll(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false}, retT.bools())
	retT = T.MaskedAny(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]bool{true, false, true, false}, retT.bools())
	retT = T.MaskedCount(1).(*Dense)
	assert.Equal([]int{2, 2}, []int(retT.shape))
	assert.Equal([]int{3, 0, 3, 0}, retT.ints())
	retT = T.NonMaskedCount(1).(*Dense)
	assert.Equal([]int{0, 3, 0, 3}, retT.ints())

	retT = T.MaskedAll(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{false, false, false, false, false, false}, retT.bools())
	retT = T.MaskedAny(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]bool{true, true, true, true, true, true}, retT.bools())
	retT = T.MaskedCount(2).(*Dense)
	assert.Equal([]int{2, 3}, []int(retT.shape))
	assert.Equal([]int{1, 1, 1, 1, 1, 1}, retT.ints())
	retT = T.NonMaskedCount(2).(*Dense)
	assert.Equal([]int{1, 1, 1, 1, 1, 1}, retT.ints())

	//non-contiguous mask case 5
	T = New(Of(Bool), WithShape(2, 3, 2), WithMaskStrides([]bool{false, false, false}))
	assert.Equal(0, len(T.mask))
	assert.Equal(0, T.MaskedCount())
	assert.Equal(12, T.NonMaskedCount())
	assert.Equal(false, T.MaskedAll())
	assert.Equal(false, T.MaskedAny())

}

func BenchmarkMaskedAll(b *testing.B) {
	T := New(Of(Bool), WithShape(1, 100000))
	T.ResetMask(true)
	for n := 0; n < b.N; n++ {
		T.MaskedAll()
	}
}

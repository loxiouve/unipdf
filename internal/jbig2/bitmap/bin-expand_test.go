/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loxiouve/unipdf/v3/common"
)

// TestExpandBinaryReplicate tests the expand binary functions.
func TestExpandBinaryReplicate(t *testing.T) {
	// Having some test bitmap with the data:
	//
	// 10001000 10001000 11000000
	// 10001000 10001000 10000000
	// 10001000 10001001 00000000
	// 10001000 10001000 00000000
	data := []byte{
		0x88, 0x88, 0xC0,
		0x88, 0x88, 0x80,
		0x88, 0x89, 0x00,
		0x88, 0x88, 0x00,
	}

	s, err := NewWithData(20, 4, data)
	require.NoError(t, err)
	common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))

	t.Run("EqualFactors", func(t *testing.T) {
		t.Run("One", func(t *testing.T) {
			// When the factors are equal to themself and to '1'
			// the bitmap should be copied.
			d, err := expandReplicate(s, 1)
			require.NoError(t, err)

			// the pointers are different
			assert.False(t, d == s)
			// but the values are the same
			assert.Equal(t, d, s)

			// When the factors are equal to themself and to '1'
			// the bitmap should be copied.
			d, err = expandBinaryReplicate(s, 1, 1)
			require.NoError(t, err)

			// the pointers are different
			assert.False(t, d == s)
			// but the values are the same
			assert.Equal(t, d, s)

			d, err = expandBinaryPower2(s, 1)
			require.NoError(t, err)

			// the pointers are different
			assert.False(t, d == s)
			// but the values are the same
			assert.Equal(t, d, s)
		})

		t.Run("Two", func(t *testing.T) {
			// When the factors are both equal and set to '2'
			// then the expandBinaryFactor2 function should be used
			// and the result would be twice bigger
			d, err := expandReplicate(s, 2)
			require.NoError(t, err)

			assert.Equal(t, 2*s.Width, d.Width)
			assert.Equal(t, 2*s.Height, d.Height)
			assert.Equal(t, 5, d.RowStride)

			// the data should be expanded by '2'
			// 11000000 11000000 11000000 11000000 11110000
			// 11000000 11000000 11000000 11000000 11110000
			// 11000000 11000000 11000000 11000000 11000000
			// 11000000 11000000 11000000 11000000 11000000
			// 11000000 11000000 11000000 11000011 00000000
			// 11000000 11000000 11000000 11000011 00000000
			// 11000000 11000000 11000000 11000000 00000000
			// 11000000 11000000 11000000 11000000 00000000

			expected := []byte{
				0xC0, 0xC0, 0xC0, 0xC0, 0xF0,
				0xC0, 0xC0, 0xC0, 0xC0, 0xF0,
				0xC0, 0xC0, 0xC0, 0xC0, 0xC0,
				0xC0, 0xC0, 0xC0, 0xC0, 0xC0,
				0xC0, 0xC0, 0xC0, 0xC3, 0x00,
				0xC0, 0xC0, 0xC0, 0xC3, 0x00,
				0xC0, 0xC0, 0xC0, 0xC0, 0x00,
				0xC0, 0xC0, 0xC0, 0xC0, 0x00,
			}
			assert.Equal(t, expected, d.Data)
		})

		t.Run("Four", func(t *testing.T) {
			// When the factor of expading is equal to '4' then
			// the result bitmap is a copy expanded by four in all directions.
			// While `diff := s.Rowstride * 4  - d.Rowstride`
			// the diff value can have three different values.
			// i.e.:
			// s.width: 13 -> rowstride = 2; d.Width = 52 -> 7 	 | 4 * 2 = 8  | 8 - 7 	= 1
			// s.width: 20 -> rowstride = 3; d.Width = 80 -> 10  | 4 * 3 = 12 | 12 - 10 = 2
			// s.width: 10 -> rowstride = 2; d.Width = 40 -> 5   | 4 * 2 = 8  | 8 - 5   = 3
			t.Run("Diff1", func(t *testing.T) {
				// Having a test bitmap of width 13 and height 4 with data:
				//
				// 10100100 11010000
				// 01011101 00101000
				// 11010110 10111000
				// 01010101 01010000
				//			    ^ padding starts here
				data := []byte{
					0xA4, 0xD0,
					0x5D, 0x28,
					0xD6, 0xB8,
					0x55, 0x50,
				}
				s, err := NewWithData(13, 4, data)
				require.NoError(t, err)

				d, err := expandReplicate(s, 4)
				require.NoError(t, err)

				assert.Equal(t, 4*s.Width, d.Width)
				assert.Equal(t, 4*s.Height, d.Height)
				assert.Equal(t, 4*s.RowStride-d.RowStride, 1)

				// the expanded data should be:
				// 11110000 11110000 00001111 00000000 11111111 00001111 00000000
				// 11110000 11110000 00001111 00000000 11111111 00001111 00000000
				// 11110000 11110000 00001111 00000000 11111111 00001111 00000000
				// 11110000 11110000 00001111 00000000 11111111 00001111 00000000

				// 00001111 00001111 11111111 00001111 00000000 11110000 11110000
				// 00001111 00001111 11111111 00001111 00000000 11110000 11110000
				// 00001111 00001111 11111111 00001111 00000000 11110000 11110000
				// 00001111 00001111 11111111 00001111 00000000 11110000 11110000

				// 11111111 00001111 00001111 11110000 11110000 11111111 11110000
				// 11111111 00001111 00001111 11110000 11110000 11111111 11110000
				// 11111111 00001111 00001111 11110000 11110000 11111111 11110000
				// 11111111 00001111 00001111 11110000 11110000 11111111 11110000

				// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
				// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
				// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
				// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
				expected := []byte{
					0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,
					0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,
					0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,
					0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,

					0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,
					0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,
					0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,
					0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,

					0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,
					0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,
					0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,
					0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,

					0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
					0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
					0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
					0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
				}
				assert.Equal(t, expected, d.Data)
			})

			t.Run("Diff2", func(t *testing.T) {
				d, err := expandReplicate(s, 4)
				require.NoError(t, err)

				assert.Equal(t, 4*s.Width, d.Width)
				assert.Equal(t, 4*s.Height, d.Height)
				assert.Equal(t, 10, d.RowStride)

				// Expected data:
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11111111 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11111111 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11111111 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11111111 00000000

				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000

				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00001111 00000000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00001111 00000000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00001111 00000000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00001111 00000000 00000000

				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 00000000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 00000000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 00000000 00000000
				// 11110000 00000000 11110000 00000000 11110000 00000000 11110000 00000000 00000000 00000000
				expected := []byte{
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xFF, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xFF, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xFF, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xFF, 0x00,

					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00,

					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x0F, 0x00, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x0F, 0x00, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x0F, 0x00, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x0F, 0x00, 0x00,

					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0x00, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0x00, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0x00, 0x00,
					0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0xF0, 0x00, 0x00, 0x00,
				}
				assert.Equal(t, expected, d.Data)
			})

			t.Run("Diff3", func(t *testing.T) {
				// Having a test bitmap of width 10 and height 4 with data:
				//
				// 10100100 11000000
				// 01011101 00000000
				// 11010110 10000000
				// 01010101 01000000
				//			  ^ padding starts here

				data := []byte{
					0xA4, 0xC0,
					0x5D, 0x00,
					0xD6, 0x80,
					0x55, 0x40,
				}
				s, err := NewWithData(10, 4, data)
				require.NoError(t, err)

				d, err := expandReplicate(s, 4)
				require.NoError(t, err)

				assert.Equal(t, s.Width*4, d.Width)
				assert.Equal(t, s.Height*4, d.Height)
				assert.Equal(t, s.RowStride*4-d.RowStride, 3)

				// the expanded data should be like:
				//
				// 11110000 11110000 00001111 00000000 11111111
				// 11110000 11110000 00001111 00000000 11111111
				// 11110000 11110000 00001111 00000000 11111111
				// 11110000 11110000 00001111 00000000 11111111

				// 00001111 00001111 11111111 00001111 00000000
				// 00001111 00001111 11111111 00001111 00000000
				// 00001111 00001111 11111111 00001111 00000000
				// 00001111 00001111 11111111 00001111 00000000

				// 11111111 00001111 00001111 11110000 11110000
				// 11111111 00001111 00001111 11110000 11110000
				// 11111111 00001111 00001111 11110000 11110000
				// 11111111 00001111 00001111 11110000 11110000

				// 00001111 00001111 00001111 00001111 00001111
				// 00001111 00001111 00001111 00001111 00001111
				// 00001111 00001111 00001111 00001111 00001111
				// 00001111 00001111 00001111 00001111 00001111
				expected := []byte{
					0xF0, 0xF0, 0x0F, 0x00, 0xFF,
					0xF0, 0xF0, 0x0F, 0x00, 0xFF,
					0xF0, 0xF0, 0x0F, 0x00, 0xFF,
					0xF0, 0xF0, 0x0F, 0x00, 0xFF,

					0x0F, 0x0F, 0xFF, 0x0F, 0x00,
					0x0F, 0x0F, 0xFF, 0x0F, 0x00,
					0x0F, 0x0F, 0xFF, 0x0F, 0x00,
					0x0F, 0x0F, 0xFF, 0x0F, 0x00,

					0xFF, 0x0F, 0x0F, 0xF0, 0xF0,
					0xFF, 0x0F, 0x0F, 0xF0, 0xF0,
					0xFF, 0x0F, 0x0F, 0xF0, 0xF0,
					0xFF, 0x0F, 0x0F, 0xF0, 0xF0,

					0x0F, 0x0F, 0x0F, 0x0F, 0x0F,
					0x0F, 0x0F, 0x0F, 0x0F, 0x0F,
					0x0F, 0x0F, 0x0F, 0x0F, 0x0F,
					0x0F, 0x0F, 0x0F, 0x0F, 0x0F,
				}
				assert.Equal(t, expected, d.Data)
			})
		})

		t.Run("Eight", func(t *testing.T) {
			t.Run("Diff4", func(t *testing.T) {
				d, err := expandReplicate(s, 8)
				require.NoError(t, err)

				assert.Equal(t, s.Width*8, d.Width)
				assert.Equal(t, s.Height*8, d.Height)
				assert.Equal(t, s.RowStride*8-d.RowStride, 4)
				// expected data:
				//
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000

				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000

				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000

				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
				// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000

				expected := []byte{
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,

					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,

					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,

					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				}
				assert.Equal(t, expected, d.Data)
			})
		})

		t.Run("Different", func(t *testing.T) {
			d, err := expandReplicate(s, 3)
			require.NoError(t, err)

			assert.Equal(t, s.Width*3, d.Width)
			assert.Equal(t, s.Height*3, d.Height)

			// 11100000 00001110 00000000 11100000 00001110 00000000 11111100 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000000 11111100 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000000 11111100 00000000

			// 11100000 00001110 00000000 11100000 00001110 00000000 11100000 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000000 11100000 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000000 11100000 00000000

			// 11100000 00001110 00000000 11100000 00001110 00000111 00000000 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000111 00000000 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000111 00000000 00000000

			// 11100000 00001110 00000000 11100000 00001110 00000000 00000000 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000000 00000000 00000000
			// 11100000 00001110 00000000 11100000 00001110 00000000 00000000 00000000
			expected := []byte{
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xFC, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xFC, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xFC, 0x00,

				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xE0, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xE0, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xE0, 0x00,

				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x07, 0x00, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x07, 0x00, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x07, 0x00, 0x00,

				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0x00, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0x00, 0x00,
				0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0x00, 0x00,
			}
			assert.Equal(t, expected, d.Data)
		})
	})

	t.Run("NonEqualFactors", func(t *testing.T) {
		d, err := expandBinaryReplicate(s, 3, 2)
		require.NoError(t, err)

		assert.Equal(t, s.Width*3, d.Width)
		assert.Equal(t, s.Height*2, d.Height)

		// 11100000 00001110 00000000 11100000 00001110 00000000 11111100 00000000
		// 11100000 00001110 00000000 11100000 00001110 00000000 11111100 00000000

		// 11100000 00001110 00000000 11100000 00001110 00000000 11100000 00000000
		// 11100000 00001110 00000000 11100000 00001110 00000000 11100000 00000000

		// 11100000 00001110 00000000 11100000 00001110 00000111 00000000 00000000
		// 11100000 00001110 00000000 11100000 00001110 00000111 00000000 00000000

		// 11100000 00001110 00000000 11100000 00001110 00000000 00000000 00000000
		// 11100000 00001110 00000000 11100000 00001110 00000000 00000000 00000000
		expected := []byte{
			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xFC, 0x00,
			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xFC, 0x00,

			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xE0, 0x00,
			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0xE0, 0x00,

			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x07, 0x00, 0x00,
			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x07, 0x00, 0x00,

			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0x00, 0x00,
			0xE0, 0x0E, 0x00, 0xE0, 0x0E, 0x00, 0x00, 0x00,
		}
		assert.Equal(t, expected, d.Data)
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Run("Factor", func(t *testing.T) {
			_, err := expandReplicate(s, -1)
			assert.Error(t, err)

			_, err = expandBinaryReplicate(s, 2, -1)
			assert.Error(t, err)

			_, err = expandBinaryPower2(s, -1)
			assert.Error(t, err)
		})

		t.Run("NilSource", func(t *testing.T) {
			_, err = expandReplicate(nil, 5)
			assert.Error(t, err)

			_, err = expandBinaryReplicate(nil, 2, 3)
			assert.Error(t, err)

			_, err = expandBinaryPower2(nil, 5)
			assert.Error(t, err)
		})
	})
}

func TestBinaryPower2Low(t *testing.T) {
	t.Run("InvalidFactor", func(t *testing.T) {
		// allowed factors are 2,4,8
		d := New(80, 80)
		s := New(80, 80)
		err := expandBinaryPower2Low(d, s, 10)
		require.Error(t, err)
	})

	t.Run("InvalidSize", func(t *testing.T) {
		// for factor two the 'd' image should be twice as big as the 's'.
		d := New(80, 80)
		s := New(80, 80)
		err := expandBinaryPower2Low(d, s, 2)
		require.Error(t, err)
	})

	t.Run("Two", func(t *testing.T) {
		// When the factors are both equal and set to '2'
		// then the expandBinaryFactor2 function should be used
		// and the result would be twice bigger
		// Having some test bitmap with the data:
		//
		// 10001000 10001000 11000000
		// 10001000 10001000 10000000
		// 10001000 10001001 00000000
		// 10001000 10001000 00000000
		data := []byte{
			0x88, 0x88, 0xC0,
			0x88, 0x88, 0x80,
			0x88, 0x89, 0x00,
			0x88, 0x88, 0x00,
		}

		s, err := NewWithData(20, 4, data)
		require.NoError(t, err)

		d := New(40, 8)
		err = expandBinaryPower2Low(d, s, 2)
		require.NoError(t, err)

		assert.Equal(t, 2*s.Width, d.Width)
		assert.Equal(t, 2*s.Height, d.Height)
		assert.Equal(t, 5, d.RowStride)

		// the data should be expanded by '2'
		// 11000000 11000000 11000000 11000000 11110000
		// 11000000 11000000 11000000 11000000 11110000
		// 11000000 11000000 11000000 11000000 11000000
		// 11000000 11000000 11000000 11000000 11000000
		// 11000000 11000000 11000000 11000011 00000000
		// 11000000 11000000 11000000 11000011 00000000
		// 11000000 11000000 11000000 11000000 00000000
		// 11000000 11000000 11000000 11000000 00000000

		expected := []byte{
			0xC0, 0xC0, 0xC0, 0xC0, 0xF0,
			0xC0, 0xC0, 0xC0, 0xC0, 0xF0,
			0xC0, 0xC0, 0xC0, 0xC0, 0xC0,
			0xC0, 0xC0, 0xC0, 0xC0, 0xC0,
			0xC0, 0xC0, 0xC0, 0xC3, 0x00,
			0xC0, 0xC0, 0xC0, 0xC3, 0x00,
			0xC0, 0xC0, 0xC0, 0xC0, 0x00,
			0xC0, 0xC0, 0xC0, 0xC0, 0x00,
		}
		assert.Equal(t, expected, d.Data)
	})

	t.Run("Four", func(t *testing.T) {
		// Having a test bitmap of width 13 and height 4 with data:
		//
		// 10100100 11010000
		// 01011101 00101000
		// 11010110 10111000
		// 01010101 01010000
		//			    ^ padding starts here
		data := []byte{
			0xA4, 0xD0,
			0x5D, 0x28,
			0xD6, 0xB8,
			0x55, 0x50,
		}
		s, err := NewWithData(13, 4, data)
		require.NoError(t, err)
		d := New(13*4, 16)

		err = expandBinaryPower2Low(d, s, 4)
		require.NoError(t, err)

		// the expanded data should be:
		// 11110000 11110000 00001111 00000000 11111111 00001111 00000000
		// 11110000 11110000 00001111 00000000 11111111 00001111 00000000
		// 11110000 11110000 00001111 00000000 11111111 00001111 00000000
		// 11110000 11110000 00001111 00000000 11111111 00001111 00000000

		// 00001111 00001111 11111111 00001111 00000000 11110000 11110000
		// 00001111 00001111 11111111 00001111 00000000 11110000 11110000
		// 00001111 00001111 11111111 00001111 00000000 11110000 11110000
		// 00001111 00001111 11111111 00001111 00000000 11110000 11110000

		// 11111111 00001111 00001111 11110000 11110000 11111111 11110000
		// 11111111 00001111 00001111 11110000 11110000 11111111 11110000
		// 11111111 00001111 00001111 11110000 11110000 11111111 11110000
		// 11111111 00001111 00001111 11110000 11110000 11111111 11110000

		// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
		// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
		// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
		// 00001111 00001111 00001111 00001111 00001111 00001111 00000000
		expected := []byte{
			0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,
			0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,
			0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,
			0xF0, 0xF0, 0x0F, 0x00, 0xFF, 0x0F, 0x00,

			0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,
			0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,
			0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,
			0x0F, 0x0F, 0xFF, 0x0F, 0x00, 0xF0, 0xF0,

			0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,
			0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,
			0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,
			0xFF, 0x0F, 0x0F, 0xF0, 0xF0, 0xFF, 0xF0,

			0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
			0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
			0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
			0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x00,
		}
		assert.Equal(t, expected, d.Data)
	})

	t.Run("Eight", func(t *testing.T) {
		// 10001000 10001000 11000000
		// 10001000 10001000 10000000
		// 10001000 10001001 00000000
		// 10001000 10001000 00000000
		data := []byte{
			0x88, 0x88, 0xC0,
			0x88, 0x88, 0x80,
			0x88, 0x89, 0x00,
			0x88, 0x88, 0x00,
		}

		s, err := NewWithData(20, 4, data)
		require.NoError(t, err)

		d := New(20*8, 4*8)
		err = expandBinaryPower2Low(d, s, 8)
		require.NoError(t, err)
		// expected data:
		//
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 11111111 00000000 00000000

		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000

		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 11111111 00000000 00000000 00000000 00000000

		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000
		// 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000

		expected := []byte{
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,

			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,

			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00,

			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}
		assert.Equal(t, expected, d.Data)
	})
}

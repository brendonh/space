package space

import (
	"testing"
)


func TestCubeFaces(t *testing.T) {
	faces := CubeFacesAll()
	
	checkFace := func(face uint8) {
		if !faces.Has(face) {
			t.Error("Missing face:", face)
		}
		faces.Unset(face)
		if faces.Has(face) {
			t.Error("Face not removed:", face)
		}
		faces.Set(face)
		if !faces.Has(face) {
			t.Error("Face not added:", face)
		}
	}

	checkFace(CUBE_TOP)
	checkFace(CUBE_LEFT)
	checkFace(CUBE_BOTTOM)
	checkFace(CUBE_RIGHT)
	checkFace(CUBE_FRONT)
	checkFace(CUBE_BACK)
}
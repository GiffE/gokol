package math

import "math"

type Mat4[T ~float32 | ~float64] [16]T
type Mat4f = Mat4[float32]

func Mat4FromQuaternion[T ~float32 | ~float64](quat Quaternion[T]) Mat4[T] {
	x2 := quat.V[0] + quat.V[0]
	y2 := quat.V[1] + quat.V[1]
	z2 := quat.V[2] + quat.V[2]

	xx2 := x2 * quat.V[0]
	xy2 := x2 * quat.V[1]
	xz2 := x2 * quat.V[2]

	yy2 := y2 * quat.V[1]
	yz2 := y2 * quat.V[2]
	zz2 := z2 * quat.V[2]

	sy2 := y2 * quat.S
	sz2 := z2 * quat.S
	sx2 := x2 * quat.S

	return Mat4[T]{
		1 - yy2 - zz2, xy2 + sz2, xz2 - sy2, 0,
		xy2 - sz2, 1 - xx2 - zz2, yz2 + sx2, 0,
		xz2 + sy2, yz2 - sx2, 1 - xx2 - yy2, 0,
		0, 0, 0, 1,
	}
}

func Mat4FromTranslation[T ~float32 | ~float64](v Vec3[T]) Mat4[T] {
	return Mat4[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v[0], v[1], v[2], 1,
	}
}

func Mat4FromAngleZ[T ~float32 | ~float64](thetaRad T) Mat4[T] {
	s, c := math.Sincos(float64(thetaRad))

	return Mat4[T]{
		T(c), T(s), 0, 0,
		-T(s), T(c), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (lhs Mat4[T]) Mul4(rhs Mat4[T]) Mat4[T] {
	return Mat4[T]{
		lhs[0]*rhs[0] + lhs[4]*rhs[1] + lhs[8]*rhs[2] + lhs[12]*rhs[3],
		lhs[1]*rhs[0] + lhs[5]*rhs[1] + lhs[9]*rhs[2] + lhs[13]*rhs[3],
		lhs[2]*rhs[0] + lhs[6]*rhs[1] + lhs[10]*rhs[2] + lhs[14]*rhs[3],
		lhs[3]*rhs[0] + lhs[7]*rhs[1] + lhs[11]*rhs[2] + lhs[15]*rhs[3],
		lhs[0]*rhs[4] + lhs[4]*rhs[5] + lhs[8]*rhs[6] + lhs[12]*rhs[7],
		lhs[1]*rhs[4] + lhs[5]*rhs[5] + lhs[9]*rhs[6] + lhs[13]*rhs[7],
		lhs[2]*rhs[4] + lhs[6]*rhs[5] + lhs[10]*rhs[6] + lhs[14]*rhs[7],
		lhs[3]*rhs[4] + lhs[7]*rhs[5] + lhs[11]*rhs[6] + lhs[15]*rhs[7],
		lhs[0]*rhs[8] + lhs[4]*rhs[9] + lhs[8]*rhs[10] + lhs[12]*rhs[11],
		lhs[1]*rhs[8] + lhs[5]*rhs[9] + lhs[9]*rhs[10] + lhs[13]*rhs[11],
		lhs[2]*rhs[8] + lhs[6]*rhs[9] + lhs[10]*rhs[10] + lhs[14]*rhs[11],
		lhs[3]*rhs[8] + lhs[7]*rhs[9] + lhs[11]*rhs[10] + lhs[15]*rhs[11],
		lhs[0]*rhs[12] + lhs[4]*rhs[13] + lhs[8]*rhs[14] + lhs[12]*rhs[15],
		lhs[1]*rhs[12] + lhs[5]*rhs[13] + lhs[9]*rhs[14] + lhs[13]*rhs[15],
		lhs[2]*rhs[12] + lhs[6]*rhs[13] + lhs[10]*rhs[14] + lhs[14]*rhs[15],
		lhs[3]*rhs[12] + lhs[7]*rhs[13] + lhs[11]*rhs[14] + lhs[15]*rhs[15],
	}
}

type Quaternion[T ~float32 | ~float64] struct {
	V Vec3[T]
	S T
}

func QuaternionFromAxisAngle[T ~float32 | ~float64](axis Vec3[T], angleRad T) Quaternion[T] {
	sin, cos := math.Sincos(float64(angleRad) * 0.5)
	return Quaternion[T]{
		S: T(cos),
		V: axis.Scale(T(sin)),
	}
}

func (lhs Quaternion[T]) Mul(rhs Quaternion[T]) Quaternion[T] {
	return Quaternion[T]{
		S: lhs.S*rhs.S - lhs.V[0]*rhs.V[0] - lhs.V[1]*rhs.V[1] - lhs.V[2]*rhs.V[2],
		V: Vec3[T]{
			lhs.S*rhs.V[0] + lhs.V[0]*rhs.S + lhs.V[1]*rhs.V[2] - lhs.V[2]*rhs.V[1],
			lhs.S*rhs.V[1] + lhs.V[1]*rhs.S + lhs.V[2]*rhs.V[0] - lhs.V[0]*rhs.V[2],
			lhs.S*rhs.V[2] + lhs.V[2]*rhs.S + lhs.V[0]*rhs.V[1] - lhs.V[1]*rhs.V[0],
		},
	}
}

func Ortho2D[T ~float32 | ~float64](left, right, bottom, top T) Mat4[T] {
	rl := right - left
	tb := top - bottom

	return Mat4[T]{
		2 / rl, 0, 0, 0,
		0, 2 / tb, 0, 0,
		0, 0, 1, 0, // Depth range [0, 1]
		-(right + left) / rl, -(top + bottom) / tb, 0, 1,
	}
}

func Perspective[T ~float32 | ~float64](fov, aspectRatio, near, far T) Mat4[T] {
	tanThetaOver2 := T(math.Tan(float64(fov * (math.Pi / 360))))
	return Mat4[T]{
		1 / tanThetaOver2, 0, 0, 0,
		0, aspectRatio / tanThetaOver2, 0, 0,
		0, 0, (near + far) / (near - far), -1,
		0, 0, (2 * near * far) / (near - far), 0,
	}
}

func LookAt[T ~float32 | ~float64](eye, center, up Vec3[T]) Mat4[T] {
	forward := center.Sub(eye).Normalize()
	right := forward.Cross(up).Normalize()
	up_ := right.Cross(forward)

	return Mat4[T]{
		right[0], up_[0], -forward[0], 0,
		right[1], up_[1], -forward[1], 0,
		right[2], up_[2], -forward[2], 0,
		-right.Dot(eye), -up_.Dot(eye), forward.Dot(eye), 1,
	}
}

func Rotate[T ~float32 | ~float64](angle T, axis Vec3[T]) Mat4[T] {
	a := axis.Normalize()
	c := T(math.Cos(float64(DegToRad(angle))))
	s := T(math.Sin(float64(DegToRad(angle))))
	C := T(1) - c

	x := a[0]
	y := a[1]
	z := a[2]

	return Mat4[T]{
		x*x*C + c, x*y*C - z*s, x*z*C + y*s, 0,
		y*x*C + z*s, y*y*C + c, y*z*C - x*s, 0,
		z*x*C - y*s, z*y*C + x*s, z*z*C + c, 0,
		0, 0, 0, 1,
	}
}

func DegToRad[T ~float32 | ~float64](deg T) (rad T) {
	return deg * (math.Pi / 180)
}

func RadToDeg[T ~float32 | ~float64](rad T) (deg T) {
	return rad * (180 / math.Pi)
}

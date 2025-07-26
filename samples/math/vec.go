package math

import "math"

type Vec2f = Vec2[float32]
type Vec3f = Vec3[float32]

type Vec2[T ~float32 | ~float64] [2]T

func (v Vec2[T]) X() T { return v[0] }
func (v Vec2[T]) Y() T { return v[1] }

func (lhs Vec2[T]) Dot(rhs Vec2[T]) T {
	return (lhs[0] * rhs[0]) + (lhs[1] * rhs[1])
}

func (lhs Vec2[T]) Magnitude() T {
	return T(math.Sqrt(float64(lhs.Dot(lhs))))
}

func (lhs Vec2[T]) Scale(s T) Vec2[T] {
	return Vec2[T]{
		lhs[0] * s,
		lhs[1] * s,
	}
}

func (lhs Vec2[T]) Normalize() Vec2[T] {
	return lhs.Scale(1 / lhs.Magnitude())
}

func (lhs Vec2[T]) Cross(rhs Vec2[T]) T {
	return lhs[0]*rhs[1] - lhs[1]*rhs[0]
}

func (lhs Vec2[T]) Add(rhs Vec2[T]) Vec2[T] {
	return Vec2[T]{
		lhs[0] + rhs[0],
		lhs[1] + rhs[1],
	}
}

func (lhs Vec2[T]) Sub(rhs Vec2[T]) Vec2[T] {
	return Vec2[T]{
		lhs[0] - rhs[0],
		lhs[1] - rhs[1],
	}
}

type Vec3[T ~float32 | ~float64] [3]T

func (lhs Vec3[T]) Dot(rhs Vec3[T]) T {
	return (lhs[0] * rhs[0]) + (lhs[1] * rhs[1]) + (lhs[2] * rhs[2])
}

func (lhs Vec3[T]) Magnitude() T {
	return T(math.Sqrt(float64(lhs.Dot(lhs))))
}

func (lhs Vec3[T]) Scale(s T) Vec3[T] {
	return Vec3[T]{
		lhs[0] * s,
		lhs[1] * s,
		lhs[2] * s,
	}
}

func (lhs Vec3[T]) Normalize() Vec3[T] {
	return lhs.Scale(1 / lhs.Magnitude())
}

func (lhs Vec3[T]) Cross(rhs Vec3[T]) Vec3[T] {
	return Vec3[T]{
		lhs[1]*rhs[2] - rhs[1]*lhs[2],
		lhs[2]*rhs[0] - rhs[2]*lhs[0],
		lhs[0]*rhs[1] - rhs[0]*lhs[1],
	}
}

func (lhs Vec3[T]) Add(rhs Vec3[T]) Vec3[T] {
	return Vec3[T]{
		lhs[0] + rhs[0],
		lhs[1] + rhs[1],
		lhs[2] + rhs[2],
	}
}

func (lhs Vec3[T]) Sub(rhs Vec3[T]) Vec3[T] {
	return Vec3[T]{
		lhs[0] - rhs[0],
		lhs[1] - rhs[1],
		lhs[2] - rhs[2],
	}
}

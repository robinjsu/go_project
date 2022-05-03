from typing import NamedTuple

class Point(NamedTuple):
    x: int
    y: int

class Box(NamedTuple):
    x0: int
    y0: int
    x1: int
    y1: int

    def contains(self, p: Point) -> bool:
        return p.x >= self.x0 and p.x < self.x1 and p.y >= self.y0 and p.y < self.y1

    
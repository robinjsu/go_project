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
        '''
        Returns whether the given point is within the bounds of the box. Exclusive of the Box.x1 and Box.y1 values
        :param p: Point object to test
        '''
        return p.x >= self.x0 and p.x < self.x1 and p.y >= self.y0 and p.y < self.y1

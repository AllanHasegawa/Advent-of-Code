We don't really wanna rotated 90 degress. Just the x/y flipped.

See how the coordinates gets complicated below.


# Default

012 (0,0)=0;(0,1)=1;(0,2)=2
345 (1,0)=3;(1,1)=4;(1,2)=5
678 (2,0)=6;(2,1)=7;(2,2)=8

# Rotated

630 (0,0)=6;(0,1)=3;(0,2)=0
741 (1,0)=7;(1,1)=4;(1,2)=1
852 (2,0)=8;(2,1)=5;(2,2)=2

---

default(x,y) = rotated(y,size(x)-x-1)
default(0,0) = rotated(0,2)
default(0,1) = rotated(1,2)
default(0,2) = rotated(2,2)
default(1,0) = rotated(0,1)
default(1,1) = rotated(1,1)
default(1,2) = rotated(2,1)
default(2,0) = rotated(0,0)
default(2,1) = rotated(1,0)
default(2,2) = rotated(2,0)
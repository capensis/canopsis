

def isThree(n):
    return True if n == 3 else False

v = [1,2,4,51,1,2, 3,5]
a = 2
b = 2

if isThree(a) or isThree(b) or any(isThree(x) for x in v):
    print "IS THREE"


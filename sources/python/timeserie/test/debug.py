from canopsis.timeserie import TimeSerie
from canopsis.timeserie.timewindow import TimeWindow, Period
from time import time

t = time()

duration = 20
gap = 5

stop = int(t)
start = stop - duration

points = [(i, i) for i in range(start, stop)]

p = Period(second=gap)
tw = TimeWindow(start=start, stop=stop)
ts = TimeSerie(period=p)

print(start, stop)
print(ts.calculate(timewindow=tw, points=points))

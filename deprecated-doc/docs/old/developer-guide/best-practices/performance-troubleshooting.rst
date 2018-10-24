.. _dev-practices-perf:

Performance troubleshooting
===========================

This guide aims to help developpers to troubleshot memory leaks in their implementations of widgets and algorithms against canopsis core.

As the core is still not really mature (as of february 2015), and there is no automatable or exact process to do it, it is more a style guide and a preview of few methods used to gather information.


Use Chrome's Timeline tool to find the location of the problem
--------------------------------------------------------------

Chrome's Timeline tool can be used to graph the page's memory consumption over time, and more handy informations, such as :

- the number of listeners currently registered (listeners as in "mouse click listener on a button")
- the number of DOM nodes instanciated (not only nodes that are attached to the page DOM)

The hardest part of the job when hunting memory leaks is the more often to detect them and to find their location. Thus, it is usually a good practice to start with a wide scope and to narrow the test scope progressively.

First of all, go to a view that you want to test. It can be a view with lots of widgets (with lots of mixins, managing lots of data), and start the timeline recording by hitting the recording button.

Then you'll have to provoke view or data refresh, to check if memory is correctly released.

As a first glance, the more refresh you cause, the better it is. It is also adviced to make these refresh be triggered with a regular interval (the "periodicrefresh" mixin is handy for that).


Checking if data is correctly released
--------------------------------------

Now you can see how the RAM usage and some other indicators evolved as you made the frontend content refresh.

First important tracked metric to check if everything is ok, the number of registrated listeners. This number should not evolve dramatically (and it should stay constant between refreshs is you used the "periodicrefresh" mixin).


.. note ::
   Don't worry if the number of registrated DOM nodes or your ram consumption is soaring on your recording, there are garbage collection processes which runs regularly (one every few minutes and one every tens of minutes). Being aware of that, you should consider benchmarking regularly on periods that lasts more than 10 minutes.


If you don't have a constant number at the beggining and at the end of your recording, and you can't explain it logically (i.e you have the same amount of data managed on the view), then it might be the sign of a memory leak problem.


Narrow the inspected scope
--------------------------

The second step to take when a memory leak is found on a broad scope is to narrow it to find the more precisely possible where the leak is.

To succeed easily on this part of the process, you'll have to remove possible causes of the problem and do a benchmark again to check if the leak is still here.

Here is some ways to find more precisely where the problem is :

- Remove some widgets on the view
- Remove some widget's mixins
- Reconfigure widgets and their mixins to manage less data (or no data at all)


Thereby, you'll might be able to logically locate the source of the problem within few tries.


Fix the leak
------------

Here are recurrent causes of memory leaks :

- Listeners that are set on didInsertElement but not cleared on willDestroyElement
- Objects not released from the memory

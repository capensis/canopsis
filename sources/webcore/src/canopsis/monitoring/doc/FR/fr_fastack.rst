.. _FR__UI_FastACK:

==============
Canopsis Event
==============

This document describes the Fast Ack UI Button

.. contents::
   :depth: 2

----------
References
----------


-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Vincent CANDEAU", "2016/04/04", "0.1", "Document creation", ""

--------
Contents
--------

Description
-----------

Fast ack will have the exact same features of Normal ACK :ref:`ACK <FR__Ack>`.

The purpose of this feature is designed to put an acknowlegdement on event without having to fill any form.* to ack an event without filling the popup information. 

The new button isdisplay on the left of the Slow ACK Button

The Fast ACK feature can be acces throught a button near the normal ACK button. 

Like Normal ACK, Fast ACK feature have a right on it. 
Normal ACK et Fast ACK right is both needed if you want to display the Fast ACK Button

An event in Canopsis is the representation of asynchronously incoming data, sent by
a :ref:`connector <FR__Connector>`.

Functionnal test
----------------

Case: Normal case
~~~~~~~~~~~~~~~~~
- Login on  canopsis with root/root
- Go to events page
- On each line you shouldn't see button with either check (ACK) or **underline check (Fast ACK)**
- **CTRL + E (Edit mode)**
- On the event list go to **sendevent mixin options**
- Inside the form change **fastackmsg** message
- Click on the **underline check (Fast ACK)**
- You should see a purple badge with
    - Actual date
    - root owner
    - Previously configured message    
- You shouldn't see a blue badge with ticket number

Case: On manager account assign right ACK without Fast ACK 
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
- Login on  canopsis with root/root
- Go to profiles = permissions
- Edit "Manager" role
- Add rights on view_event (read/write)
- Add "ack" right
- Remove "fastack" right 
- Save
- Login on canopsis with canopsis/canopsis
- Go to events page
- On each line you should see button with check (ACK) but you shouldn't see the **underline check (Fast ACK)**


Case: On manager account assign right Fast ACK without Fast ACK 
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
- Login on canopsis with root/root
- Go to profiles = permissions
- Edit "Manager" role
- Add rights on view_event (read/write)
- Remove "ack" right
- Add "fastack" right 
- Save
- Login on canopsis with canopsis/canopsis
- Go to events page
- On each line you shouldn't see button with either check (ACK) or **underline check (Fast ACK)**
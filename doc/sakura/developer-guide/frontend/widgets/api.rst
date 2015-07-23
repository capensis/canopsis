.. _dev-frontend-widgets-api:

API
===

1. Widget
---------

It's on this super-class that all widgets are based. Models and controllers of
child classes derive from Widget ( model ) and WidgetCOntroller.


2. Model
--------

We find following attributes in model :
* title :  ,
* widget_type : class of widget ,
* item : reference of associated item
* connectionParameters : attribute to develop when database connectors will exist.

We found as attributes :

* container : reference on container
  - belongsTo relationship
  - options : async

* widget : reference on widget
  - belongsTo relationship
  - options : embedded,async polymorphic

* row : position of left superior corner on X axis
* col : position of left superior corner on Y axis
* rowspan: height of item by report of layout
* colspan: width of item by report of layout


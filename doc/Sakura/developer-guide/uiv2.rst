Developments in UIV2
====================
I. Application
--------------
1.1. Adapter
------------
1.2. Serializer
---------------
We have defined and redefined any methods :

    * normalizeEmbeddedRelationships :
        - parameters : 
            + type: type of main model
            + hash collected by ajax query : { data : array( records ),
                                               success : true }
        - returns normalized hash : { data : array(normalized records),
                                      embedded : array( normalized embedded records ),
                                      success : true }
        - manages hasMany/belongsTo relationships and polymorphic option


    * Redefined methods :
        - normalizeId : 
          The redefined method takes as parameter the type of the hash and uses
          it to generate id.

        - normalize :
          The method has been redefined to take into account the modification on
          normalizeId

        - normalizePayload :
          By default, this method returns data in entry without traitment. 
          Redefined method calls normalizeId and normalizeEmbeddedRelationships
          to format payload.
        
        - extractFindAll :
          Redefined method calls normalizePayload to normalize data in entry.
          then it pushes in store embedded records. Finally, method returns 
          normalized main records to the store.

II. Userview
------------
2.1. Adapter
------------
We redefined several methods to recover data in application store because type
in Mongo database is 'view' and not 'userview'.

    * find : 
    * findAll :
    
2.2. Serializer
---------------
UserviewSerializer derives directly of ApplicationSerializer 

2.3. Routes
-----------
Currently we have two routes:
* userview : route for controller userview
* userviews : route for controller userviews
  - model : recovers all userview records in store
  - setupController : set content of userviews controller with data returned
    by model method.


2.4. Model
----------
Model has as attributes:
* name :
* container : main container of view
* internal :
* enable :

2.5. Controllers
----------------
We have two controllers for userview model:
* userview : controller for one model only.
* userviews : controller to manage data of several userview models ( arraycontroller ).

2.6. Templates
--------------
This Ember object has two templates:
* userview : displays container content
* userviews : displays id, name and of every record in a table with an button
to access of individual view.

III. Widgets
------------
3.1. Widget
-----------
It's on this super-class that all widgets are based. Models and controllers of
child classes derive from Widget ( model ) and WidgetCOntroller.

3.1.1. Model
~~~~~~~~~~~~
We find following attributes in model :
* title :  ,
* widget_type : class of widget ,
* item : reference of associated item
* connectionParameters : attribute to develop when database connectors will exist.

3.1.2. Controller
~~~~~~~~~~~~~~~~~

3.2. Container
--------------

This widget allows to structure one view or a set of widgets

3.2.1. Model
~~~~~~~~~~~~
Specific container attributes are :
* userview : reference on linked view
* items : list of associated items
* layout_cols : number of columns of container layout,
* layout_rows : number of rows.

3.2.2. Controller
~~~~~~~~~~~~~~~~~

3.2.3. Template
~~~~~~~~~~~~~~~
create 'gridster' list

3.3. Cell
---------
Widget of test

3.3.1. Model
~~~~~~~~~~~~
Specific attributes :
* color : color of background 

3.3.2. Controller
~~~~~~~~~~~~~~~~~

3.3.3. Template 
~~~~~~~~~~~~~~~

3.4. Grid
---------

3.3.1. Model
~~~~~~~~~~~~
Specific attributes :
* color : color of background 

3.3.2. Controller
~~~~~~~~~~~~~~~~~

3.3.3. Template 
~~~~~~~~~~~~~~~

IV. Item
--------
The goal of this structure is to manage sizing and positioning of widget 
object in interface.

4.1. Model
----------
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

4.2. Controller
----------------

4.3. Template
-------------
Create an elemnt of html list with gridster properties.

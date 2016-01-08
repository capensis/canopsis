.. _dev-frontend-widgets-mixins:

Managing plug-ins in mixins
===========================

Introduction
------------

Canopsis frontend allow to display widget with **add-in features**. These features can include additionnal html elements that is put directly inside widgets, and additionnal Javascript actions.

Requirements
------------

This chapter describes a more advanced topic than the previous ones. It will require a solid knowledge in *Javascript*, *OOP*, and concepts such as *Software Architecture* and *Modularity*.

Key concepts
------------

A **widget** add-in is mainly composed of a javascript file that registers an EmberJS Mixin. This mixin can be added or removed to a widget at the view loading, or when a widget is modified.


A **partial** is an Ember template that is included directly in a MVC context, without changing this context (ie. as if the template's code was manually pasted in the template where the partial is called).


A **slot** is a place where elements can be put dynamically without writing down statically what to display in the code.


Partial slots
-------------

It is sometimes useful to make widgets add-ins add bits of html in the widget layout, such as additional buttons, labels, and so on.

To make this possible, **partial slots** have to be put in the base widget templates. They will be used as anchors to make mixin append html bits into the widget.

Example
"""""""

Let's take back our label widget that was defined in the first example of this guide.

This widget basically insert a "<h1>" html tag (a title), with a text inside.

If we assume that we want to add a slot for action buttons that are providided by add-ins, the base template of the widget should be modified :


``core/widgets/label/template.html``

.. code-block:: html

   {{partialslot slot=controller._partials.actionbuttons}}
   <h1>
      {{texttodisplay}}
   </h1>


Now that a slot is added to the widget, open the mixin javascript file. It should look like this :

``core/mixins/infobutton.js``

.. code-block:: javascript

   define([
       'ember',
       'app/application'
   ], function(Ember, Application) {

       var mixin = Ember.Mixin.create({
         [some code here]
       });

       Application.ActionbuttonMixin = mixin;

       return mixin;
   });


The template of the button we want to add could be like this :

``core/templates/actionbutton-showinfo.html``


.. code-block:: html

     <button class="btn btn-default btn-sm" {{action 'showInfo'}} >
         {{glyphicon "eye"}}
     </button>


The mixin code now need to be modified to manage the button template and the corresponding action :

``core/mixins/infobutton.js``

.. code-block:: javascript

   define([
       'ember',
       'app/application'
   ], function(Ember, Application) {

       var mixin = Ember.Mixin.create({
          partials: {
              actionbuttons: ['actionbutton-showinfo']
          },

          actions: {
              showInfo: function() {
                  console.info('clicked on button');
              }
          },

          [some code here]
       });

       Application.ActionbuttonMixin = mixin;

       return mixin;
   });


When mixins are applied to widgets, ``partials`` dictionnaries are all merged in ``_partials`` (with the same logic as actions are merged in every Ember controller) and slots can provide all the templates that are designed to be displayed.

When this mixin is applied to any widget, if the widget implements the correct partial slot, a button will be available to the user.


Advantages
----------

This is an encouraged practice because it keeps algorithms and templates modular and make them appliable to any kind of widgets.

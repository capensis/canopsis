.. _dev-frontend-widgets-userpref:

User preferences
================


Overview
--------

User preferences in canopsis allow a user to store a widget statement into database making this widget having a specific state at application load.
Thus, a user preference is stored in the user preference mongo collection and **is identified with both a user id and a widget id**.

Implementation
--------------

User preferences in canopsis is a document that allow developers to save a widget user's statement into json format within the user preference database collection.
These user preference will be loaded at widget (re)load and will give it a statement that will render it different for a single widget to different users.

To declare a new user preference, open the schema of the model in which the user preference will be set. The property managed as an user preference is actually like any other schema property, but with a "isUserPreference" option set to true.

That's it. The Ember model class generated have a submodel dedicated to store user preferences, accessible via the 'userPreferencesModel' key, and the preferences values can be accessed from the model instance directly.

Persistance
-----------

User preferences are loaded when the model instance is loaded, and saved when the user model is saved. But on some cases, developers might want to save user preferences only, without saving the whole model.

This is done thanks to the **userconfiguration** `Ember js mixin <http://emberjs.com/api/classes/Ember.Mixin.html>`_ that add a method to the (widget) controller using it.

* ``saveUserConfiguration`` is called by the developper to trigger data save to the mongo collection for the widget it is called on.



Example
-------

Let say the widget text have a custom title for each user to save, then in the user interface the user will be able to fill a text field with a custom title. A user action will then tell the widget to set the title in ``userParams`` javascript object.


.. code-block:: javascript

   {
      "type": "object",

      '...': '...',

      "properties": {
         "customTitle": {
            "isUserPreference": true,
            "type": "string"
         },

         '...': '...'
      }
   }


.. code-block:: javascript

   set(widgetController, 'model.customTitle', 'user input title');
   widgetController.saveUserConfiguration();

Then run saveUserConfiguration will filldatabase with a document that looks like:

.. code-block:: javascript

   {
      widgetId: 'widgetGUID_aUserId',
      'userId': 'aUserId',
      'widget_preferences': {
         'customTitle': 'user input title',
         'anotherPreference': 'aPreferenceValue'
      }
      '...': '...'
   }

At widget initilization, this previous document is loaded into the widget and the model is filled with the appropriate fields.


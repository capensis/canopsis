User preferences
================


overview
--------

User preferences in canopsis allow a user to store a widget statement into database making this widget having a specific state at application load.
Thus, a user preference is stored in the user preference mongo collection and **is identified with both a user id and a widget id**.

Implementation
--------------

User preferences in canopsis is a document that allow developers to save a widget user's statement into json format within the user preference database collection.
These user preference will be loaded at widget (re)load and will give it a statement that will render it different for a single widget to different users.

This is done thanks to the **userconfiguration** `Ember js mixin <http://emberjs.com/api/classes/Ember.Mixin.html>`_ that add two methods to the widget controller using it.

* ``saveUserConfiguration`` is called by the developper to trigger data save to the mongo collection for the widget it is called on.
* ``loadUserConfiguration`` is called on widget load and fetch and apply user configuration to the widget at initialization time. It is called in the widget root class constructor and is applyed only if the widget implements the mixin.

The **saveUserConfiguration** method expects the widget's ``userParams`` javascript object is filled with keys that have a preference as value.
The **loadUserConfiguration** method will fetch database user configuration for the current widget and will fill widget's ``userParams`` javascript object making user preference information available at widget instanciation and usable in widget representation.

Example
-------

Let say the widget text have a custom title for each user to save, then in the user interface the user will be able to fill a text field with a custom title. A user action will then tell the widget to set the title in ``userParams`` javascript object.

.. code-block:: javascript

   set(widgetController, 'userParams.customTitle', 'user input title');
   widgetController.saveUserConfiguration();

Then run saveUserConfiguration will filldatabase with a document that looks like:

.. code-block:: javascript

   {
      widgetId: 'widgetGUID_123456',
      'userId': 'aUserId',
      'widget_preferences': {
         'customTitle': 'user input title',
         'anotherPreference': 'aPreferenceValue'
      }
      '...': '...'
   }

At widget initilization, this previous document is loaded into the widget and the **'widget_preferences'** field will fill the **userParams** widget field.


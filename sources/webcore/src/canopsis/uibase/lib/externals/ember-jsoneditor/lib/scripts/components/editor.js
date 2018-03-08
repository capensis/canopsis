'use strict';
/**
JSON Editor component.
*/
Ember.JSONEditor.EditorComponent = Ember.Component.extend({
    /**
    Element tag name.
    */
    tagName: 'div',
    /**
    Element classes.
    */
    classNames: ['jsoneditor-component'],

    /**
    Cached editor.
    */
    _editor: undefined,
    /**

    */
    editor: function() {
        var self = this;
        console.log('editor');
        if (Ember.isEqual(self.get('_editor'), undefined)) {
            // Empty, create it.
            var container = self.get('element');
            if (Ember.isEqual(container, undefined)) {
                return undefined;
            } else {
                var options = self.get('options');
                var json = self.get('json');
                var editor = new jsoneditor.JSONEditor(container, options, json);
                self.set('_editor', editor);
                return editor;
            }
        } else {
            // Editor is already created and cached.
            return self.get('_editor');
        }
    }.property('element', 'options', 'json'),

    /**
    JSON object.
    */
    json: {},

    /**
    Object with options.
    */
    options: function() {
        var props = this.getProperties([
            'mode',
            'modes',
            '_change',
            'search',
            'history',
            'name',
            'indentation',
            'error'
        ]);
        // Rename
        props.change = props._change;
        delete props._change;
        // Add reference to this component
        props.component = this;
        return props;
    }.property(
        'mode',
        'modes',
        '_change',
        'search',
        'history',
        'name',
        'indentation',
        'error'
    ),

    /**
    Editor mode. Available values:
    'tree' (default), 'view',
    'form', 'text', and 'code'.
    */
    mode: 'tree',

    /**
    Create a box in the editor menu where the user can switch between the specified modes.
    Available values: see option mode.
    */
    modes: ['tree','view','form','text','code'],

    /**
    Callback method, triggered
    on change of contents
    */
    change: function() {
        console.log('JSON Editor changed!');
    },

    /**
     Set a callback method triggered when an error occurs.
     Invoked with the error as first argument.
     The callback is only invoked for errors triggered by a users action.
    */
    error: function(error) {
        console.error('An error occured: ', error);
    },

    /**
    Editor updated JSON.
    */
    _updating: false,

    /**
    Change event handler.
    Triggers `change()` which is user defined.
    */
    _change: function() {
        var self = this.component;
        var editor = self.get('editor');
        var json = editor.get();
        //
        self.set('_updating', true);
        self.set('json', json);
        self.set('_updating', false);
        // Trigger Change event
        if (!!self.change) {
            self.change();
        }
    },

    /**
    Enable search box.
    True by default
    Only applicable for modes
    'tree', 'view', and 'form'
    */
    search: true,
    /**
    Enable history (undo/redo).
    True by default
    Only applicable for modes
    'tree', 'view', and 'form'
    */
    history: true,
    /**
    Field name for the root node.
    Only applicable for modes
    'tree', 'view', and 'form'
    */
    name: 'JSONEditor',
    /**
    Number of indentation
    spaces. 4 by default.
    Only applicable for
    modes 'text' and 'code'
    */
    indentation: 4,

    /**
    Editor observer.
    */
    editorDidChange: function() {
        var self = this;
        self.get('editor');
    }.observes('editor').on('didInsertElement'),

    /**
    JSON observer.
    */
    jsonDidChange: function() {
        var self = this;
        if (Ember.isEqual(self.get('_updating'), false)) {
          var editor = self.get('editor');
          var json = self.get('json');
          editor.set(json);
        }
    }.observes('json'),

    /**
    Mode observer.
    */
    modeDidChange: function() {
        var self = this;
        var editor = self.get('editor');
        var mode = self.get('mode');
        editor.setMode(mode);
    }.observes('mode'),

    /**
    Name observer.
    */
    nameDidChange: function() {
        var self = this;
        var editor = self.get('editor');
        var name = self.get('name');
        editor.setName(name);
    }.observes('name'),


});

Ember.Handlebars.helper('json-editor', Ember.JSONEditor.EditorComponent);

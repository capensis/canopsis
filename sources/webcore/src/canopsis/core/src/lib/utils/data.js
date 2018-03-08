/**
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'DataUtils',
    after: ['UtilityClass', 'HashUtils'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');
        var hashUtils = container.lookupFactory('utility:hash');

        var _loggedUserController,
            _applicationSingleton;

        var dataUtils = Utility.create({

            name: 'data',

            getLoggedUserController: function() {
                return _loggedUserController;
            },

            setLoggedUserController: function(loggedUserController) {
                _loggedUserController = loggedUserController;
            },

            getEmberApplicationSingleton: function() {
                return _applicationSingleton;
            },

            setEmberApplicationSingleton: function(applicationInstance) {
                _applicationSingleton = applicationInstance;
            },

            getStore: function() {
                console.warn("this should not be used as there is not only one store in Canopsis. This might lead to unexpected behaviour");
                return this.getEmberApplicationSingleton().__container__.lookup('store:main');
            },

            //TODO change parentElement term to something more descriptive
            addRecordToRelationship: function(record, parentElement, relationshipKey, cardinality) {
                console.log('addRecordToRelationship', arguments);
                if (cardinality === "hasMany") {
                    console.log("addRecordToRelationship hasMany", relationshipKey, arguments, parentElement);
                    parentElement.get(relationshipKey).pushObject(record);
                } else if (cardinality === "belongsTo") {
                    console.log("addRecordToRelationship belongsTo", relationshipKey, arguments, parentElement);
                    parentElement.set(relationshipKey, record);
                }
            },

            /**
             * @method download
             * @param {string} content the file content
             * @param {string} filename the file name
             * @param {string} contentType the file content type
             *
             * Automatically download content as a file
             */
            download: function (content, filename, contentType) {
                if(!contentType) {
                    contentType = 'application/octet-stream';
                }

                var a = document.createElement('a');
                var blob = new Blob([content], {'type': contentType});
                a.href = window.URL.createObjectURL(blob);
                a.download = filename;
                a.click();
            },

            /**
             * @method uploadFilePopup
             * @param {fn(fileInput)} callback to handle when the user select a file
             *
             * Shows a file selection popup window
             */
            uploadFilePopup: function(callback) {
                if (window.File && window.FileReader && window.FileList && window.Blob) {
                    //do your stuff!
                    var input = $(document.createElement('input'));
                    input.attr("type", "file");
                    input.trigger('click'); // opening dialog
                    input.change(function () {
                        if(typeof callback === 'function') {
                            callback(this);
                        }
                    });
                } else {
                    alert('The File APIs are not fully supported by your browser.');
                }
            },

            /**
             * @method uploadTextFilePopup
             * @param {fn(name:string, filetype:string, filesize:number, content:string)} callback to handle when the user select a file
             *
             * Shows a file selection popup window, and handle it with a callback dedicated to a single text file.
             */
            uploadTextFilePopup: function(callback) {
                this.uploadFilePopup(function(fileInput) {
                    var file = fileInput.files[0];

                    if (file) {
                        var r = new FileReader();
                        r.onload = function(e) {
                            var contents = e.target.result;
                            if(typeof callback === 'function') {
                                callback(file.name, file.type, file.size, contents);
                            }
                        }
                        r.readAsText(file);
                    } else {
                        alert("Failed to load file");
                    }
                })
            },

            /**
             * @function cleanJSONIds
             * @param {Object} recordJSON
             * @return {Object} the cleaned record
             */
            cleanJSONIds: function (recordJSON) {
                for (var key in recordJSON) {
                    var item = recordJSON[key];
                    //see if the key need to be cleaned
                    if(key === 'id' || key === '_id' || key === 'widgetId' || key === 'preference_id' || key === 'EmberClass') {
                        delete recordJSON[key];
                    }

                    //if this item is an object, then recurse into it
                    //to remove empty arrays in it too
                    if (typeof item === 'object') {
                        this.cleanJSONIds(item);
                    }
                }

                if(recordJSON !== null && recordJSON !== undefined && (recordJSON.crecord_type !== undefined || recordJSON.xtype !== undefined)) {
                    recordJSON['id'] = hashUtils.generateId(recordJSON.xtype || recordJSON.crecord_type || 'item');
                    recordJSON['_id'] = recordJSON['id'];
                }

                return recordJSON;
            }
        });
        application.register('utility:data', dataUtils);
    }
});

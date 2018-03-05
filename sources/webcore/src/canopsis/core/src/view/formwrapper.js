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
    name: 'FormwrapperView',
    after: 'DragUtils',
    initialize: function(container, application) {
        var drag = container.lookupFactory('utility:drag');

        var get = Ember.get,
            set = Ember.set,
            isNone =Ember.isNone;

        /**
         * @class FormwrapperView
         * @extends Ember.View
         * @constructor
         */
        var view = Ember.View.extend({
            /**
             * @method init
             */
            init: function() {
                this._super();
                console.log("formwrapper view init", this, get(this, 'controller'));

                set(this,'controller.widgetwrapperView', this);
            },

            /**
             * @method didInsertElement
             */
            didInsertElement: function () {
                //TODO watch out ! garbage collector might not work here! Possible memory leak.
                drag.setDraggable(this.$('#formwrapper .modal-title'), this.$('#formwrapper'));
            },

            /**
             * @method willDestroyElement
             */
            willDestroyElement: function () {
                this.$("#formwrapper").modal("hide");
                this._super();
            },


            /**
             * @method registerHooks
             * Controller -> View Hooks
             */
            registerHooks: function() {
                this.hooksRegistered = true;

                console.log("registerHooks", this);
                this.get("controller").on('validate', this, this.hidePopup);
                this.get("controller").on('hide', this, this.hidePopup);
                this.get("controller").on('rerender', this, this.rerender);

                var formwrapperView = this;

                //TODO "on" without "off"
                this.$('#formwrapper').on('hidden.bs.modal', function () {
                    formwrapperView.onPopupHidden.apply(formwrapperView, arguments);
                });
            },

            /**
             * @method unregisterHooks
             */
            unregisterHooks: function() {
                this.get("controller").off('validate', this, this.hidePopup);
                this.get("controller").off('hide', this, this.hidePopup);
                this.get("controller").off('rerender', this, this.rerender);
            },


            /**
             * @method showPopup
             */
            showPopup: function() {
                console.log("view showPopup");
                if(!this.hooksRegistered) {
                    this.registerHooks();
                }

                //show and display centered !
                this.$("#formwrapper").modal('show').css('top',0).css('left',0);

                if(get(this, 'controller.form')) {
                    get(this, 'controller.form').send('show');
                }
            },

            /**
             * @method hidePopup
             */
            hidePopup: function() {
                console.log("view hidePopup");
                this.$("#formwrapper").modal("hide");
            },

            /**
             * @method onPopupHidden
             */
            onPopupHidden: function() {
                console.log("onPopupHidden", arguments);
                var submit = get(this, 'controller.form.submit');
                if (!isNone(submit) && submit.state() === "pending") {
                    console.info("rejecting form submission");
                    get(this, 'controller.form').send("abort");
                }
            }
        });

        application.register('view:formwrapper', view);
    }
});

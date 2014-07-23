/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

define([
	'ember',
	'app/application'
], function(Ember, Application) {

	/**
	 * Implements configuration management for controllers
	 * @mixin
	 */
	Application.ConfigurableMixin = Ember.Mixin.create({
		actions: {
			//TODO as this shows the WIDGET edit form, it has to be refactored in the widget clas as soon as it will be created
			showEditForm: function() {
				//FIXME this does not work, because of the hack line 61
				crecord_type = "account";
				console.log("Form generation for", crecord_type);

				var crecordformController = Application.CrecordformController.create();
				crecordformController.set("crecord_type", crecord_type);
				crecordformController.set("editMode", "edit");
				crecordformController.set("editedRecordController", this);

				this.send('showEditFormWithController', crecordformController);
			}
		},

		init: function() {
			console.log("init");
			this.refreshConfiguration();
			this._super.apply(this, arguments);
		},

		refreshConfiguration: function() {
			console.log("refreshConfiguration");

			var me = this;
			try{
				this.store.findQuery("account", { filter: { "firstname" : "Cano" } }).then(function(queryResult) {
					console.log("results found", queryResult);
					me.set("configuration", queryResult.get("content")[0]);
				});
			} catch (e) {
				console.error(e.message, e.stack);
			}
		}
	});

	return Application.ConfigurableMixin;
});

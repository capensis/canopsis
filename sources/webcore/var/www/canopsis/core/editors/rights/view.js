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
	'app/application',
	'app/view/arraytocollectioncontrol'
], function(Application) {
		Application.RightView = Application.ArrayToCollectionControlView.extend({

			init: function() {
				var value = this.getValue();
				var contentREF = this.getContent();

				var readTemplate = {name : "r", icon : "eye-open" , label : "Read" };
				var writeTemplate = {name : "w", icon : "pencil", label : "Write" };
				this.addTemplate(readTemplate, value, contentREF);
				this.addTemplate(writeTemplate, value, contentREF);

				this._super(true);
			}

	});
	return Application.RightView;
});
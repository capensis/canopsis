//need:app/lib/view/cfile_window.js
/*
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
*/
//Common class provied usefull functions for controllers from canopsis lib. It allows code factorisation and easy fonction access
Ext.define('canopsis.controller.common', {
	requires: [
		'canopsis.lib.view.cfile_window'
	],
	
	//This method displays an import pop up that allows user to import an object from both a json file or json character data
	//parent item refers to the window the pop up is linked to, object type defines the model to use in order to validate object import
	filepopup: function(parentItem, objectType) {
		
		//Capitalize object type definition for class matching
		objectType = objectType.toLowerCase();
		var objectTypeCapitalized = objectType.charAt(0).toUpperCase() + objectType.slice(1);
		
		log.debug('Open file popup', this.logAuthor);
		
		var importView = function(objs) {
			console.log(objs);

			if(!Ext.isArray(objs)) {
				objs = [objs];
			}

			for(var i = 0; i < objs.length; i++) {
				var obj = objs[i];

				var record = Ext.create('canopsis.model.' + objectTypeCapitalized, obj);
				var recordId = objectType + '.' + global.account.user + '.' + global.gen_id();

				record.set('_id', recordId);
				record.set('id', recordId);
				record.set('leaf', true);
				parentItem.add_to_home(record, false);
			}
		};


		var config = {
			_fieldLabel: _(objectTypeCapitalized + ' dump'),
			copyPasteZone: true,
			constrainTo: parentItem
		};

		var popup = Ext.create('canopsis.lib.view.cfile_window', config);
		popup.show();

		popup.on('save', function(info) {
			if(info.file) {
				var file = info.value[0];

				if(file.type === '' || file.type === 'application/json') {
					log.debug('Import '+ objectType +' file', this.logAuthor);

					var reader = new FileReader();
					reader.onload = (function(e) {
						importView(Ext.decode(e.target.result));
					}).bind(this);

					reader.readAsText(file);
					popup.close();
				}
				else {
					log.debug('Wrong file type: ' + file.type, this.logAuthor);
					global.notify.notify(_('Wrong file type'), _('Please choose a correct json file'), 'info');
				}
			}
			else  {
				importView(info.value);
				popup.close();
			}
		}, this);
	}

});

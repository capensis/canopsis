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

Ext.define('canopsis.lib.form.field.ccolorfield', {
	extend: 'Ext.ux.ColorField',

	alias: 'widget.ccolorfield',
	colors: global.default_colors,

	
    menuListeners : {
        select: function(m, d){
        	this.setValue(d);
        	//dirty hack to make ccolorfield work in cellediting
			if(this.ownerCt.editingPlugin && this.ownerCt.editingPlugin.record)
			{
				var rec =  this.ownerCt.editingPlugin.record;
				rec.set('color', this.getValue());
				this.ownerCt.editingPlugin.ccomponentlist.addRecord(rec);
			}
       },
        show : function(){
            this.onFocus();
        },
        hide : function(){
            this.focus();
            var ml = this.menuListeners;
            this.menu.un("select", ml.select,  this);
            this.menu.un("show", ml.show,  this);
            this.menu.un("hide", ml.hide,  this);
        }
    }
});

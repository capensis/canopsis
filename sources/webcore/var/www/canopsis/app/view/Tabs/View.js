/*
#--------------------------------
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
# ---------------------------------
*/
Ext.define('canopsis.view.Tabs.View' , {
	extend: 'Ext.tab.Panel',
	alias: 'widget.TabsView',

	activeTab: 0, // index or id
	bodyBorder: false,
	componentCls: 'cps-headertabs',
	plain: true,
	tabBar:{
		//plain:true,
		items:[{
		    xtype: 'tbfill'
		},{
		    iconCls: 'icon-control-play',
			tooltip: _('Rotate view'),
			xtype:'button',
			border: 0,
			style:'background-color:#e0e0e0;background-image:none;',
			enableToggle: true,
			scope:this,
			toggleHandler: function(button, state) {
				if (state) {
					Ext.Msg.prompt(
						_('Question'),
						_('Enter the delay to stay on each view, in minutes'),
						function(button,text,obj){
							var number = parseInt(text)
							if(isNaN(number)){
								global.notify.notify(_('Warning'),_('You must enter only number'))
								return
							}
							if(button == 'ok'){
								this.up('tabpanel').fireEvent('AutoRotateView',true, number)
								this.setIconCls('icon-control-pause');
							}
						},
						button
					)
				}else {
					button.setIconCls('icon-control-play');
					button.up('tabpanel').fireEvent('AutoRotateView',false)
				}
			}
		},{
		    iconCls: 'icon-control-repeat',
			tooltip: _('Refresh view'),
		    xtype:'button',
		    style:'background-color:#e0e0e0;background-image:none;',
		    border: 0,
		    handler:function(btn,state){
				this.up('tabpanel').fireEvent('reload_active_view')
			}
		}]
	},

});
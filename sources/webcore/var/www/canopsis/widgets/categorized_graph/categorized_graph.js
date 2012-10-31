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

// initComponent -> afterContainerRender 	-> setchartTitle -> ready -> doRefresh -> onRefresh -> addDataOnChart
//                                			-> setOptions                             			-> getSerie
//											-> createChart


Ext.define('widgets.categorized_graph.categorized_graph' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.categorized_graph',

	logAuthor: '[categorized_graph]',

	defaultHtml: '<center><span class="icon icon-loading" />' + _('This wiget is no longer available, please look at "diagram" widget') +'</center>'
	

});
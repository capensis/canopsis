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
    'utils',
], function(Ember, utils) {

    Ember.Handlebars.helper('ack', function(value) {

	    var tooltipHtml = ['<i>' + _('Date') + '</i> : <br/>',
	    	utils.dates.timestamp2String(value.timestamp) +' <br/> ',
	    	value.author +' <br/><br/> ',
	    	'<i>'+_('Commentaire') +' :</i> : <br/>' + value.comment].join('');

		var guid = utils.hash.generate_GUID();
		var ack  = '<span id="'+ guid +'" class="badge bg-maroon" data-html="true" title="" data-original-title="' + tooltipHtml + '"><i class="fa fa-check"></i></span>';

		//Triggers tooltip display once loaded /!\ hack
		setTimeout(function () {
			$('#' + guid).tooltip();
		}, 1000);

		return new Ember.Handlebars.SafeString(ack);
    });

});
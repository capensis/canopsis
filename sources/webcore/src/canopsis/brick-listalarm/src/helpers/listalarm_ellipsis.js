/*
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

 (function () {

    var set = Ember.set,
        get = Ember.get,
        isNone = Ember.isNone,
        isArray = Ember.isArray;


    var helper = function (txt, nbChar, name) {
        var style = "";

        if (typeof (nbChar) !== 'number')
			nbChar = 0

        var html = '';
		var output = "";
        if (txt === undefined) {
            txt = ""
        }
        if (txt.length > nbChar) {
			output = txt.substring(0, nbChar)
			output = output + "…"
		} else {
			output = txt
		}
        html += '<p onclick="showOutput(\'';
        html += txt.replace(/([^>\r\n]?)(\r\n|\n\r|\r|\n)/g, "<br />").replace(/\"/g, "&quot;").replace(/'/g, "&rsquo;");
        html += '\',\'' + name + ' \')" ' + style + '>' + output + '</p>';

        return new Ember.String.htmlSafe(html);
    };
    //declaring helper this way allow it to be used as simple function somewhere else.
    Handlebars.registerHelper('listalarm_ellipsis', helper);
    Ember.Handlebars.helper('listalarm_ellipsis', helper);
    window.ellipsis = helper;
    window.showOutput = function (output, name) {
        if ($("#modal-default-output").length) {
            hideOutput();
        }
        var modal = '';
        modal += '<div class="modal fade in" id="modal-default-output" style="display: block; margin-left: 40%; width: 20%; padding-right: 15px;">';
        modal += '  <div class="modal-dialog" style="width: auto;">';
        modal += '    <div class="modal-content">';
        modal += '      <div class="modal-header">';
        modal += '        <button type="button" class="close" data-dismiss="modal" aria-label="Close" onclick=hideOutput()>';
        modal += '        <span aria-hidden="true">×</span></button>';
        modal += '        <h4 class="modal-title">'+ name +'</h4>';
        modal += '      </div>';
        modal += '      <div class="modal-body">';
        modal += '        <p style="margin: auto; word-wrap: break-word;">' + output + '</p>';
        modal += '      </div>';
        modal += '    </div>';
        modal += '  </div>';
        modal += '</div>';
        $('body').append(modal);
    }

    window.hideOutput = function () {
        $('#modal-default-output').remove();
    }
})();

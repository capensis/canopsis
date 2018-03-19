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
        
        
    var helper = function (txt, nbChar) {
        var style = "";
        
        if (typeof (nbChar) !== 'number')
        nbChar = 0
        
        if (nbChar > 0 && txt.length > nbChar) {
            style += 'style="';
            style += 'text-overflow:ellipsis;';
            style += 'width: ' + nbChar + 'ch;';
            style += 'white-space:nowrap;';
            style += 'overflow:hidden;"';
        }
        var coucou = function () {
            alert('YEAHHH');
        }

        var html = '';
        html += '<p onclick=showOutput("';
        html += txt;
        html += '") ';
        html += style;
        html += '>';
        html += txt;
        html += '</p>';
        return new Ember.String.htmlSafe(html);
    };

    //declaring helper this way allow it to be used as simple function somewhere else.
    Handlebars.registerHelper('listalarm_ellipsis', helper);
    Ember.Handlebars.helper('listalarm_ellipsis', helper);
    window.ellipsis = helper;
    window.showOutput = function (output) {
        if ($("#modal-default-output").length) {
            hideOutput();
        }
        var modal = '';
        modal += '<div class="modal fade in" id="modal-default-output" style="display: block; height: 30%; margin-left: 25%; padding-right: 15px;">';
        modal += '  <div class="modal-dialog">';
        modal += '    <div class="modal-content">';
        modal += '      <div class="modal-header">';
        modal += '        <button type="button" class="close" data-dismiss="modal" aria-label="Close" onclick=hideOutput()>';
        modal += '        <span aria-hidden="true">×</span></button>';
        modal += '        <h4 class="modal-title">Output</h4>';
        modal += '      </div>';
        modal += '      <div class="modal-body">';
        modal += '        <p>' + output + '</p>';
        modal += '      </div>';
        modal += '    </div>';
        modal += '  </div>';
        modal += '</div>';
        $('body').append(modal);
        //putTrigger()
        /*$(window).click(function () {
            //Hide the menus if visible
            if ($("#modal-default-output").length) {
                hideOutput();
            }
        });

        $('#modal-default-output').click(function (event) {
            event.stopPropagation();
        });*/
    }
    window.putTrigger= function (){
         $(window).click(function () {
             //Hide the menus if visible
             if ($("#modal-default-output").length) {
                 hideOutput();
             }
         });

         $('#modal-default-output').click(function (event) {
             event.stopPropagation();
         });
    }
    window.hideOutput = function () {
        $('#modal-default-output').remove();
    }
})();

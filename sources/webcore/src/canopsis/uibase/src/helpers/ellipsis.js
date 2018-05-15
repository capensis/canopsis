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

(function() {

    var set = Ember.set,
        isNone = Ember.isNone,
        isArray = Ember.isArray;

    var helper = function(txt, nbChar) {
        var style = "";

        if(typeof(nbChar) !== 'number')
            nbChar = 0

        if(nbChar > 0 && txt.length > nbChar ){
            style += 'style="';
            style += 'text-overflow:ellipsis;';
            style += 'width: ' + nbChar + 'ch;';
            style += 'white-space:nowrap;';
            style += 'overflow:hidden;"';
        }

        var html = '<p ' + style + ' onclick="selfEllipsis(' + nbChar + ',this)">' + txt + '</p>';

        return new Ember.String.htmlSafe(html);
    };

    //declaring helper this way allow it to be used as simple function somewhere else.
    Handlebars.registerHelper('ellipsis', helper);
    Ember.Handlebars.helper('ellipsis', helper);
    window.ellipsis = helper;
    window.selfEllipsis = function(nbChar,obj){

        if(typeof(nbChar) !== 'number' || nbChar <= 0)
            return

        if(obj.style.textOverflow){
            obj.removeAttribute("style");
        }else{
            if($(obj).text().length > nbChar){
                obj.style.textOverflow = "ellipsis";
                obj.style.width = nbChar + "ch";
                obj.style.whiteSpace= "nowrap";
                obj.style.overflow= "hidden";
            }
        }
        
    }
})();
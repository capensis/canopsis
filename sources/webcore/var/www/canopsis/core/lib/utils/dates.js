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

define([], function() {

	var dates = {
		timestamp2String: function (value, format) {
			function addZero(i) {
			    return (i < 10 ? '0'+ i +'' : i +'');
			}

			var a = new Date(value*1000);
			var months = ["January", "February", "March", "April", "May",
				      "June", "July", "August", "September", "October",
				      "November", "December"];
			var year = a.getFullYear();
			var month = months[a.getMonth()];
			var date = addZero(a.getDate());
			var hour = addZero(a.getHours());
			var min = addZero(a.getMinutes());
			var sec = addZero(a.getSeconds());
			var time = "";

			switch(format) {
				case 'f':
				    time = [date, month, year].join(' ') + ' ' + [hour, min, sec].join(':') ;
				    break;
				case 'r':
				    time = [date, addZero(a.getMonth()), year].join('/') + ' <br>' + [hour, min, sec].join(':') ;
				    break;

				default:
				    time = [date, month, year].join(' ') + ' ' + [hour, min, sec].join(':') ;
				    break;
			}
			return time;
		},

		//bottom is not used yet
		locale: 'fr',
		months: ['January','February','March','April','May','June','July','August','September','October','November','December'],
		days: ['Monday', 'Thuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],

		setLang: function (lang) {
			dates.locale = lang;
			if (lang === 'fr') {
				dates.month = ['Janvier', 'Février', 'Mars', 'Avril', 'Mai', 'Juin', 'Juillet', 'Août', 'Octobre', 'Novembre', 'Décembre'];
				dates.days = ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche'];
			}
		},
		dateFormat:'YYYY/MM/DD'
	};

	return dates;
});
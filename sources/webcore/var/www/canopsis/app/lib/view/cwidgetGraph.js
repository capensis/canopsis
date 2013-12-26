//need:app/lib/view/cperfstoreValueConsumerWidget.js
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
Ext.define('canopsis.lib.view.cwidgetGraph', {
    extend: 'canopsis.lib.view.cperfstoreValueConsumerWidget',

    logAuthor: '[widgets][graph]',

    time_window: global.commonTs.day,

    initComponent: function() {
        this.callParent(arguments);

        this.setChartOptions();
        this.series = [];

        this.chart = undefined;

        this.on('boxready', this.createChart, this);
        this.on('resize', this.renderChart, this);
    },

    afterContainerRender: function() {
        this.callParent(arguments);

        if(this.chart !== undefined) {
            this.renderChart();
        }
    },

    setChartOptions: function() {
        var me = this;

        this.options = {
            cwidget: function() {
                return me;
            },

            crosshair: {
                mode: 'x'
            },

            grid: {
                borderWidth: {
                    top: 0,
                    bottom: 0,
                    right: 0,
                    left: 0
                },
                hoverable: true,
                clickable: true
            },

            xaxis: {
                position: 'bottom',
                mode: 'time',
                timeformat: '%H:%M:%S'
            }
        };
    },

    createChart: function() {
        this.plotcontainer = $('#' + this.wcontainerId);
        this.chart = $.plot(this.plotcontainer, this.series, this.options);
    },

    renderChart: function() {
        log.debug('Rendering chart');

        this.chart.destroy();
        this.createChart();

        this.chart.draw();
    },

    addSerie: function(config) {
        log.debug('Add serie:');
        log.dump(config);

        this.series.push(config);
    },

    getSerie: function(serieId) {
        for(var i = 0; i < this.series.length; i++) {
            var serie = this.series[i];

            if(serieId === serie.serieId) {
                return serie;
            }
        }

        return null;
    },

    doRefresh: function(from, to) {
        this.refreshNodes(from, to);
    },

    onRefresh: function(data, from, to) {
        this.callParent(arguments);

        console.log(data);

        if(data.length > 0) {
            for(var i = 0; i < data.length; i++) {
                var info = data[i];
                var node = this.nodesByID[info.node];

                var serie = this.getSerie(info.node);

                if(serie === null) {
                    serie = {
                        serieId: info.node,
                        label: node.label,
                        data: []
                    };

                    this.addSerie(serie);
                }

                for(var j = 0; j < info.values.length; j++) {
                    var value = info.values[j];

                    serie.data.push([value[0] * 1000, value[1]]);
                }
            }
        }

        this.renderChart();
    }
});
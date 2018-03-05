/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @author Capensis
 */

(function($) {
    var options = {
        values: {
        }
    };

    function init(plot) {
        function draw(plot, ctx) {
            $.each(plot.getData(), function(idx, series) {
                var values = series.values || {};

                if (!values.xalign) {
                    values.xalign = function(x) {
                        return x;
                    };
                }

                if (!values.yalign) {
                    values.yalign = function(y) {
                        var shift = (series.yaxis.max * 15.0) / 100.0;

                        return y + shift;
                    };
                }
                if (!values.labelFormatter) {
                    values.labelFormatter = function(series, text) {
                        void(series);

                        return text;
                    };
                }

                if (values.show) {
                    var ps = series.datapoints.pointsize;
                    var points = series.datapoints.points;

                    var ctx = plot.getCanvas().getContext('2d');
                    var offset = plot.getPlotOffset();

                    ctx.textBaseline = 'top';
                    ctx.textAlign = 'center';

                    var shiftx = values.xalign;
                    var shifty = values.yalign;

                    for(var i = 0; i < points.length; i += ps) {
                        var point = {
                            'x': shiftx(points[i]),
                            'y': shifty(points[i + 1])
                        };

                        var text = points[i + 1];
                        var c = plot.p2c(point);

                        if(text === null || text == undefined) {
                            text = 'no data';
                        }

                        if(typeof text !== 'string') {
                            text = text.toString();
                        }

                        text = values.labelFormatter(series, text);

                        var oldStyle = ctx.fillStyle;
                        ctx.fillStyle = series.color || oldStyle;

                        ctx.fillText(text,
                            c.left + offset.left,
                            c.top + offset.top
                        );

                        ctx.fillStyle = oldStyle;
                    }
                }
            });
        }

        plot.hooks.draw.push(draw);
    }

    $.plot.plugins.push({
        init: init,
        options: options,
        name: 'chartvalues',
        version: '0.1'
    });
})(jQuery);

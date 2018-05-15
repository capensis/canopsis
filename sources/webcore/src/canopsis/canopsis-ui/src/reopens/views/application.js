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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'CanopsisUiApplicationViewReopen',
    after: 'ApplicationView',
    initialize: function(container, application) {
        var ApplicationView = container.lookupFactory('view:application');

        var get = Ember.get,
            set = Ember.set;


        ApplicationView.reopen({
            didInsertElement: function() {
                console.log('main template rendered, trigger $ documentready');

                $("[data-toggle='offcanvas']").click(function(e) {
                    e.preventDefault();

                    //If window is small enough, enable sidebar push menu
                    if ($(window).width() <= 992) {
                        $('.row-offcanvas').toggleClass('active');
                        $('.left-side').removeClass("collapse-left");
                        $(".right-side").removeClass("strech");
                        $('.row-offcanvas').toggleClass("relative");
                    } else {
                        //Else, enable content streching
                        $('.left-side').toggleClass("collapse-left");
                        $(".right-side").toggleClass("strech");
                    }
                });

                //Add hover support for touch devices
                $('.btn').bind('touchstart', function() {
                    $(this).addClass('hover');
                }).bind('touchend', function() {
                    $(this).removeClass('hover');
                });

                $("[data-widget='collapse']").click(function() {
                    //Find the box parent
                    var box = $(this).parents(".box").first();
                    //Find the body and the footer
                    var bf = box.find(".box-body, .box-footer");
                    if (!box.hasClass("collapsed-box")) {
                        box.addClass("collapsed-box");
                        bf.slideUp();
                    } else {
                        box.removeClass("collapsed-box");
                        bf.slideDown();
                    }
                });

                /*
                 * ADD SLIMSCROLL TO THE TOP NAV DROPDOWNS
                 * ---------------------------------------
                 */
                //TODO uncomment while ready


                // $(".navbar .menu").slimscroll({
                //     height: "200px",
                //     alwaysVisible: false,
                //     size: "3px"
                // }).css("width", "100%");

                /*
                 * INITIALIZE BUTTON TOGGLE
                 * ------------------------
                 */
                $('.btn-group[data-toggle="btn-toggle"]').each(function() {
                    var group = $(this);
                    $(this).find(".btn").click(function(e) {
                        group.find(".btn.active").removeClass("active");
                        $(this).addClass("active");
                        e.preventDefault();
                    });

                });

                $("[data-widget='remove']").click(function() {
                    //Find the box parent
                    var box = $(this).parents(".box").first();
                    box.slideUp();
                });

                /* Sidebar tree view */
                $(".sidebar .treeview").tree();

                /*
                 * Make sure that the sidebar is streched full height
                 * ---------------------------------------------
                 * We are gonna assign a min-height value every time the
                 * wrapper gets resized and upon page load. We will use
                 * Ben Alman's method for detecting the resize event.
                 *
                 **/
                function _fix() {
                    //Get window height and the wrapper height
                    var height = $(window).height() - $("body > .header").height();
                    $(".wrapper").css("min-height", height + "px");
                    var content = $(".wrapper").height();
                    //If the wrapper height is greater than the window
                    if (content > height)
                        //then set sidebar height to the wrapper
                        $(".left-side, html, body").css("min-height", content + "px");
                    else {
                        //Otherwise, set the sidebar to the height of the window
                        $(".left-side, html, body").css("min-height", height + "px");
                    }
                     $(".right-side").css("min-height", height + "px");
                }
                //Fire upon load
                _fix();
                //Fire when wrapper is resized
                $(".wrapper").resize(function() {
                    _fix();
                    fix_sidebar();
                });

                //Fix the fixed layout sidebar scroll bug
                fix_sidebar();
            }
        });
    }
});

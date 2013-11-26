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

widget_stream_event_template = Ext.create('Ext.XTemplate',
	'<tpl if="referer">',
		"<div class='separator-top'></div>",
	'</tpl>',
	'<tpl if="referer == undefined">',
		"<div class='separator'></div>",
	'</tpl>',
	"<table class='table'>",
	'<tr>',
		"<td class='state' style='background-color: ", '<tpl if="state==0">green</tpl>', '<tpl if="state==1">orange</tpl>', '<tpl if="state==2">red</tpl>', '<tpl if="state==3">grey</tpl>', ";'/>",
		"<td class='picture ", '<tpl if="referer">comment</tpl>', "' >",
			'<tpl if="authorId">',
				"<img src='/account/getAvatar/{authorId}'/></img>",
			'</tpl>',
			'<tpl if="authorId == undefined">',
				"<img src='widgets/stream/logo/{connector}.png'></img>",
			'</tpl>',
		'</td>',
		"<td class='", '<tpl if="referer">comment</tpl>', "'>",
			"<div class='content'>",
				'<header>',
					'<tpl if="author">',
						"<h1 class='title'>{author}</h1>",
					'</tpl>',
					'<tpl if="author == undefined">',
						'<tpl if="display_name">',
							"<h1 class='title'>{display_name}</h1>",
						'</tpl>',
						'<tpl if="display_name == undefined">',
							'<tpl if="source_type == \'resource\'">',
								"<h1 class='title'>{component} - {resource}</h1>",
							'</tpl>',
							'<tpl if="source_type == \'component\'">',
								"<h1 class='title'>{component}</h1>",
							'</tpl>',
						'</tpl>',
					'</tpl>',
					"<span id='{event_Component_id}-time' class='timestamp'>{event_date}</span></br>",
				'</header>',
				"<span class='output'>{output}</span></br> <span class='long_output'>{long_output}</span>",
				'<tpl if="referer == undefined">',
					"<div class='afooter'>",
							'<tpl if="derogation_id != undefined">',
								"<span class='icon icon-crecord_type-derogation' id='{event_Component_id}-derogation' style='float: left;'></span>",
							'</tpl>',
							"<span class='icon icon-comment'></span><div class='comment-counter' id='{event_Component_id}-nbcomment'></div>",
							"<span id='{event_Component_id}-expend-comments' class='icon icon-plus'></span>",
					'</div>',
				'</tpl>',
			'</div>',
		'</td>',
	'</tr>',
	"<tr class='tr-comments'>",
		'<td/><td/>',
		"<td id='{event_Component_id}-comments-td' class='comments' style='display: none;' colspan=2>",
			"<div id='{event_Component_id}-comments'></div>",
		'</td>',
	'</tr>',
	'</table>',
	{ compiled: true }
);

widget_stream_event_template_compact = Ext.create('Ext.XTemplate',
	"<table class='table'>",
	'<tr>',
		"<td class='state' style='background-color: ", '<tpl if="state==0">green</tpl>', '<tpl if="state==1">orange</tpl>', '<tpl if="state==2">red</tpl>', '<tpl if="state==3">grey</tpl>', ";'/>",
		"<td class='picture_compact' ", '<tpl if="referer">comment</tpl>', "' >",
			'<tpl if="authorId">',
				"<img src='/account/getAvatar/{authorId}'/></img>",
			'</tpl>',
			'<tpl if="authorId == undefined">',
				"<img src='widgets/stream/logo/{connector}.png'></img>",
			'</tpl>',
		'</td>',
		"<td class='", '<tpl if="referer">comment</tpl>', "'>",
			"<div class='content_compact'>",

				'<header>',
					"<span id='{event_Component_id}-time' class='timestamp'>{event_date}</span>",
					'<tpl if="author">',
						"<h1 class='title'>{author}</h1>",
					'</tpl>',
					'<tpl if="author == undefined">',
						'<tpl if="display_name">',
							"<h1 class='title'>{display_name}</h1>",
						'</tpl>',
						'<tpl if="display_name == undefined">',
							'<tpl if="source_type == \'resource\'">',
								"<h1 class='title'>{component} - {resource}</h1>",
							'</tpl>',
							'<tpl if="source_type == \'component\'">',
								"<h1 class='title'>{component}</h1>",
							'</tpl>',
						'</tpl>',
					'</tpl>',
					"<span class='output'>{output}</span>",
				'</header>',
				'<tpl if="referer == undefined">',
					"<div class='afooter_compact'>",
							'<tpl if="derogation_id != undefined">',
								"<span class='icon icon-crecord_type-derogation' id='{event_Component_id}-derogation' style='float: left;'></span>",
							'</tpl>',
							"<span class='icon icon-comment'></span><div class='comment-counter' id='{event_Component_id}-nbcomment'></div>",
							"<span id='{event_Component_id}-expend-comments' class='icon icon-plus'></span>",
					'</div>',
				'</tpl>',
			'</div>',
		'</td>',
	'</tr>',
	"<tr class='tr-comments'>",
		'<td/><td/>',
		"<td id='{event_Component_id}-comments-td' class='comments' style='display: none;' colspan=2>",
			"<div id='{event_Component_id}-comments'></div>",
		'</td>',
	'</tr>',
	'</table>',
	{ compiled: true }
);

widget_stream_derogation_tooltip_template = Ext.create('Ext.XTemplate',
	'<p>',
		'<b>{name}</b>', '<br/>',
		'{description}', '<br/>',
	'</p>',
	{ compiled: true }
);

Ext.define('widgets.stream.event', {
	extend: 'Ext.Component',
	alias: 'widget.stream.event',

	logAuthor: '[widgets][stream][event]',

	cls: 'event',

	raw: {},
	event_id: undefined,
	timestamp: 0,

	stream: undefined,

	comments: [],

	el_comments: undefined,
	el_btn_exp_comments: undefined,
	el_time: undefined,

	el_nbcomment: undefined,
	nbcomment: 0,

	initComponent: function() {
		if(this.id) {
			this.event_id = this.id;
			this.id = this.stream.id + '.' + this.id;
		}

		log.debug('Create event: ' + this.id, this.logAuthor);

		this.timestamp = parseInt(this.raw.timestamp);

		this.html = this.build();

		if(this.stream.compact) {
			this.style = {
				padding: '0px'
			};
		}

		this.callParent(arguments);
	},

	build: function(raw) {
		if(!raw) {
			raw = this.raw;
		}

		if(!raw['referer']) {
			raw['referer'] = undefined;
		}

		if(!raw['authorId']) {
			raw['authorId'] = undefined;
		}

		if(!raw['author']) {
			raw['author'] = undefined;
		}

		if(!raw['display_name']) {
			raw['display_name'] = undefined;
		}

		if(!raw['derogation_id']) {
			raw['derogation_id'] = undefined;
		}

		if(!raw['output']) {
			raw['output'] = '';
		}

		if(!raw['long_output']) {
			raw['long_output'] = '';
		}

		if(raw['output']) {
			raw['output'] = raw['output'].replace('\n', '</br>');
		}

		if(raw['long_output']) {
			raw['output'] = raw['output'].replace('\n', '</br>');
		}

		if(typeof(raw['connector']) === 'string') {
			raw['connector'] = raw['connector'].toLowerCase();
		}

		raw['event_date'] = rdr_elapsed_time(raw['timestamp']);
		raw['event_Component_id'] = this.id;

		if(this.stream.compact) {
			return widget_stream_event_template_compact.applyTemplate(raw);
		}
		else {
			return widget_stream_event_template.applyTemplate(raw);
		}
	},

	afterRender: function() {
		this.callParent(arguments);

		var el = this.getEl();

		//Get elements
		this.el_comments = el.getById(this.id + '-comments-td');
		this.el_btn_exp_comments = el.getById(this.id + '-expend-comments');
		this.el_nbcomment = el.getById(this.id + '-nbcomment');

		this.el_time = el.getById(this.id + '-time');

		this.el_derogation_icon = el.getById(this.id + '-derogation');

		this.init_comment_counter();

		this.bindEvents();
	},

	bindEvents: function() {
		if(this.el_btn_exp_comments) {
			this.el_btn_exp_comments.on('click', this.toggle_comments, this);
			this.el_nbcomment.on('click', this.toggle_comments, this);
		}

		this.create_tooltip();
	},

	create_tooltip: function() {
		if(this.el_derogation_icon) {
			if(this.raw.derogation_name && this.raw.derogation_description) {
				var config = {
					name: this.raw.derogation_name,
					description: this.raw.derogation_description
				};

				Ext.create('Ext.tip.ToolTip', {
					target: this.el_derogation_icon,
					html: widget_stream_derogation_tooltip_template.apply(config)
				});
			}
			else {
				log.debug('No required information, no tooltip', this.logAuthor);
			}
		}
	},

	create_comments_container: function() {
		if(!this.comments_container) {
			var items = [];

			if(this.stream.enable_comments) {
				log.debug("Create comment's container", this.logAuthor);
				this.comment_form = Ext.create('Ext.form.Panel', {
					layout: 'fit',
					border: false,
					margin: 3,
					items: [{
						xtype: 'textfield',
						emptyText: _('Leave a') + ' ' + _('comment') + ' ?',
						name: 'message',
						listeners: {
							specialkey: {
								fn: function(field, e) {
									void(field);

									if(e.getKey() === e.ENTER) {
										this.submit_comment();
									}
								},
								scope: this
							}

						}
					}]
				});

				items.push(this.comment_form);
			}

			this.comments_container = Ext.create('Ext.container.Container', {
				layout: 'anchor',
				items: items
			});

			this.comments_container.on('afterRender', function() {
				var me = this;

				now.stream_getComments(this.event_id, this.stream.max_comment, function(records) {
					log.debug(records.length + " comments for '" + me.event_id + "'", me.logAuthor);

					if(records.length > 0) {
						me.init_comment_counter();
						records.reverse();

						for(var i = 0; i < records.length; i++) {
							records[i] = Ext.create('widgets.stream.event', {
								raw: records[i],
								stream: me.stream
							});
						}

						me.comments_container.insert(0, records);
					}
				});
			}, this);

			this.comments_container.render(this.id + '-comments');
		}
	},

	submit_comment: function() {
		log.debug('Submit comment', this.logAuthor);

		var message = this.comment_form.getForm().getValues().message;
		this.comment_form.getForm().findField('message').reset();

		if(this.event_id) {
			this.stream.publish_comment(this.event_id, this.raw, message, this);
		}
		else {
			log.error("Invalid event_id: '" + this.event_id + "'", this.logAuthor);
		}
	},

	comment: function(event) {
		if(!this.comments_container && !this.stream.active) {
			return;
		}
		else if (!this.comments_container) {
			this.show_comments();
		}
		else {
			if(this.stream.active) {
				this.show_comments();
			}
			else {
				this.hide_comments();
			}

			log.debug(' + Insert comment', this.logAuthor);
			var nb = this.comments_container.items.length;
			log.debug(' + ' + nb + ' comments', this.logAuthor);

			this.comments_container.insert(nb - 1, event);

			//Clean before
			if(this.comments_container.items.length > (this.stream.max_comment + 1)) {
				log.debug(' + Remove first comment', this.logAuthor);
				var item = this.comments_container.getComponent(0);
				this.comments_container.remove(item.id, true);
			}

			this.nbcomment += 1;
			this.update_comment_counter();
		}

	},

	update_time: function() {
		if(this.el_time) {
			this.el_time.update(rdr_elapsed_time(this.timestamp));
		}

		if(this.comments_container) {
			for (var i = 0; i < this.comments_container.items.length; i++) {
				var event = this.comments_container.getComponent(i);

				//check if not a form
				if(event.event_id) {
					event.update_time();
				}
			}
		}
	},

	init_comment_counter: function() {
		var me = this;

		now.stream_countComments(this.event_id, function(count) {
			me.nbcomment = count;
			me.update_comment_counter();
		});
	},

	update_comment_counter: function() {
		if(this.el_nbcomment) {
			this.el_nbcomment.update(this.nbcomment + ' ' + _('comment(s)'));
		}
	},

	toggle_comments: function() {
		if(this.el_comments.isVisible()) {
			this.hide_comments();
		}
		else {
			this.show_comments();
		}

	},

	show_comments: function() {
		log.debug('Show comments', this.logAuthor);
		this.el_comments.show();

		this.el_btn_exp_comments.removeCls('icon-plus');
		this.el_btn_exp_comments.addCls('icon-minus');

		this.create_comments_container();

		if(this.comment_form) {
			this.comment_form.doLayout();
			this.comment_form.getForm().findField('message').focus();
		}
	},

	hide_comments: function() {
		this.el_comments.setVisibilityMode(Ext.Element.DISPLAY);

		log.debug('Hide comments', this.logAuthor);
		this.el_comments.hide();

		this.el_btn_exp_comments.addCls('icon-plus');
		this.el_btn_exp_comments.removeCls('icon-minus');
	},

 	beforeDestroy: function() {
		if(this.comments_container) {
			this.comments_container.destroy();
		}

		this.callParent(arguments);
 	}
});

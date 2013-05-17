
Ext.define('widgets.progressbar.bar', 
{
	extend: 'Ext.container.Container',
	border: 0,
	_res: undefined,
	_comp: undefined,
	_met: undefined,
	_label: undefined,
	_value : undefined,
	_gHeight: undefined,
	_colorStart:undefined,
	_colorMid: undefined,
	_colorEnd:undefined,
	_colorBg: undefined,
	_grad1: undefined,
	_grad2: undefined,
	_dispGrad: undefined,
	_fontSize: undefined,
	_boldText: undefined,
	
	initComponent: function()
	{
		this.callParent(arguments);

		this.idpb = this.id+'pb';

		if (!this._label)
			if (this._res)
				this.names = this._comp+', '+this._res+', '+this._met+': ';
			else
				this.names = this._comp+', '+this._met+': ';
		else
			this.names = this._label;

		var tpl = new Ext.Template
		(
			'<table border=0 align=center style="width:100%;margin-bottom:5px;">',
			'<tr><td style="min-width:150px;width:40%;">',
			'<div class="label">{names}</div></td>',
			'<td style="min-width:50px;">',
			'<div style="width:100%;line-height:'+this._gHeight+'px;" id={idpb}>',
			'<div class="progress-label"></div></div>',
			'</td></tr></table>'
		);

		var html = tpl.apply({
			names: this.names, 
			idpb: this.idpb,
			gHeight: this._gHeight
		});

		this.update(html);
	},

	afterRender: function()
	{
		this.jqpb = $("#"+this.idpb);
		this.displayBars();
	},

	displayBars: function()
	{
		if (this._value<40)
			var pcolor = this.shadeColor(this._colorStart,-20);
		else if (this._value<55)
			var pcolor = this.shadeColor(this._colorStart,-40); 
		else if (this._value<65)
			var pcolor = this.shadeColor(this._colorMid,-20); 
		else if (this._value<75)
			var pcolor = this.shadeColor(this._colorMid,-40); 
		else
			var pcolor = this.shadeColor(this._colorEnd,-20); 

		$("#"+this.idpb).progressbar({"value":this._value})
			.height(this._gHeight)
			.css({ 'background': this._colorBg })
			.css('border');

		var idpb_div = $('#'+this.idpb+' > div');

		if (this._dispGrad)
		{
			var grad_color = this.shadeColor(pcolor,100);
			idpb_div.css('background', '-webkit-gradient(linear, left top, left bottom, from('	+ grad_color +'), to('+pcolor+'))');
			idpb_div.css('background', '-webkit-linear-gradient(' 	+ grad_color +', '+pcolor+')');
			idpb_div.css('background', '-moz-linear-gradient('			+ grad_color +', '+pcolor+')');
			idpb_div.css('background', '-ms-linear-gradient('			+ grad_color +', '+pcolor+')');
			idpb_div.css('background', '-o-linear-gradient('			+ grad_color +', '+pcolor+')');
			idpb_div.css('background', 'linear-gradient(to bottom, '	+ grad_color +', '+pcolor+')');
		}
		else
			idpb_div.css('background',pcolor);

		var idpb_progresslabel = $('#'+this.idpb+" > .progress-label");

		idpb_progresslabel.css('background','transparent');

		if (this._boldText){
			idpb_progresslabel.css('font-weight','bold');
			$('#'+this.id+' .label').css('font-weight','bold');
		}

		if (this._fontSize)
			idpb_progresslabel.css('font-size',this._fontSize+'%');

		idpb_progresslabel.text(Math.floor(this._value)+"%");
	},
	shadeColor: function (color, percent) 
	{
		if (! color)
			return '#000000'

		var R = parseInt(color.substring(1,3),16);
		var G = parseInt(color.substring(3,5),16);
		var B = parseInt(color.substring(5,7),16);
		R = parseInt(R * (100 + percent) / 100);
		G = parseInt(G * (100 + percent) / 100);
		B = parseInt(B * (100 + percent) / 100);
		R = (R<255)?R:255;
		G = (G<255)?G:255;
		B = (B<255)?B:255;
		var RR = ((R.toString(16).length==1)?"0"+R.toString(16):R.toString(16));
		var GG = ((G.toString(16).length==1)?"0"+G.toString(16):G.toString(16));
		var BB = ((B.toString(16).length==1)?"0"+B.toString(16):B.toString(16));
		return "#"+RR+GG+BB;
	}
});

Ext.define('widgets.progressbar.progressbar' , 
{
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.progressbar',
	logAuthor: '[progressBarWidget]',
	innerText: undefined,
	refresh_number: 0,
	wcontainer_autoScroll: true,
	wcontainer_layout: {type:'anchor'},
	bodyPadding: '5 5 5 5',

	grad1: '#FFFFFF',
	colorBg: '#EEEEEE',
	colorStart: '#1BE01B',
	colorMid: '#FFCD43',
	colorEnd: '#E0251B',
	dispGrad: true,
	boldText: true,
	fontSize: 100,


	initComponent: function()
	{
		log.debug("initComponent", this.logAuthor)
		this.callParent(arguments);
		this.progressBarArray=[];
	},

	getNodeInfo: function(from,to) 
	{
		this.processNodes();
		if (this.nodeId) 
		{
			Ext.Ajax.request(
			{
				url: '/perfstore/values' 
					+ '/' + parseInt(to/1000)+ '/' + parseInt(to/1000),
				scope: this,
				params: this.post_params,
				method: 'POST',
				success: function(response) 
				{
					var data = Ext.JSON.decode(response.responseText);
					this.data = data.data;
					if (this.progressBarArray==0)
						this.createBars();
					else
						this.refreshBars();
				}, 
				failure: function(result, request) 
				{
					log.error
					(
						'get Node info, Ajax req failed ... (' 
						+ request.url + ')', this.logAuthor
					);
				}
			});
		}
	},

	createBars: function()
	{
		for (var i = 0; i < this.nodes.length; i++)
		{
			if (! this.data[i])
				continue;

			var _met = this.data[i].metric;
			var _comp = this.nodes[i].component;
			var _res = this.nodes[i].resource;
			var _val = this.data[i].values[0][1];
			var _max = this.nodes[i].extra_field.ma;
			if (this.nodes[i].extra_field.label != "")
				var _label = this.nodes[i].extra_field.label
			else
				var _label = null;
			var percent = (100*_val) / _max;
			if (_max)
			{
				this.createBarsObj
				(
					_res, _comp, _met, percent,
					this.gHeight, this.colorStart, this.colorMid, this.colorEnd, 
					this.colorBg, this.grad1, this.grad2, _label, this.dispGrad,
					this.fontSize, this.boldText
				);
			}
		}
	},
	refreshBars: function()
	{
		for (var i = 0; i < this.progressBarArray.length; i++)
		{
			_val = this.data[i].values[0][1];
			_max = this.nodes[i].extra_field.ma;
			percent = (100*_val)/_max;
			oldPercent = this.progressBarArray[i]._value;
			if (percent != oldPercent)
			{
				this.progressBarArray[i]._value = percent;
				this.progressBarArray[i].displayBars();
			}
		}
	},

	createBarsObj: function
		( res, comp, met, val, gHeight, colorStart, colorMid, colorEnd,
			colorBg, grad1, grad2, label, dispGrad, fontSize, boldText ) 
	{
		this.obj = Ext.create("widgets.progressbar.bar",
			{
				_res:res, _comp:comp, _met:met, _value:val,
				_gHeight:gHeight, _colorStart:colorStart, _colorMid:colorMid, 
				_colorEnd:colorEnd, _colorBg:colorBg, _grad1:grad1, _grad2:grad2,
				_label:label, _dispGrad:dispGrad, _fontSize:fontSize, _boldText:boldText
			}
		);
		this.down('#'+this.wcontainerId).add(this.obj);
		this.progressBarArray.push(this.obj);
	},

	processNodes: function() 
	{
		var post_params = [];
		for (var i = 0; i < this.nodes.length; i++)
			post_params.push(
			{
				id: this.nodes[i].id,
				metrics: this.nodes[i].metrics
			});
		this.post_params = 
		{
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_interval': this.aggregate_interval,
			'aggregate_max_points': this.aggregate_max_points
		};
	}
});

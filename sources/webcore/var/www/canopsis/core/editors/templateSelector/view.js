 define([
    'app/application',
    'app/view/arraytocollectioncontrol'
], function(Application) {
	Application.TemplateView = Application.ArrayToCollectionControlView.extend({
		//cssClass: "btn-items-",
		init: function() {
				var value = this.getValue();
				var contentREF = this.getContent();

				// Check if filter template byclass (add another "for" in order to filter with several class)
				//if (this.templateData.keywords.attr.model.options.templateClass !== undefined) {
				var classToGet = this.templateData.keywords.controller.content.model.options.templateClass;
				if (classToGet !== undefined) {
					for (var i = 0 ; i < Canopsis.templates.byClass[classToGet].length ; i++) {
						this.addTemplate(Canopsis.templates.byClass[classToGet][i], value, contentREF);
					}
				 //Else Add all template
				 } else {
				 	for (i = 0 ; i < Canopsis.templates.all.length ; i++) {
				 		this.addTemplate(Canopsis.templates.all[i], value, contentREF);
				 	}
				 }
			//Have to be done after
			this._super(true);
		}
	});
    return Application.TemplateView;
});

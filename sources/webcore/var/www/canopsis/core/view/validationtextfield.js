  define([
  'jquery',
  'ember',
  'app/application',
  'app/view/crecords'
], function($, Ember, Application) {

    /**
     * Generic textField for validation
     * Use Component-> validators -> validate (Ember.validators["validate"]) for validation
     */
    Application.ValidationTextField = Ember.TextField.extend({
        attr : "",
        formController : null,


        registerFieldWithController: function() {
            var formController  =  Canopsis.formwrapperController.form;
            var validationFields ;

            if ( formController ){
              var validationFields = formController.get('validationFields');
              if (validationFields){
                  validationFields.pushObject(this);
              }
            }
        }.on('didInsertElement'),

        focusOut: function() {
            this.validate();
        },

        validate : function() {
          var formController  = Canopsis.formwrapperController.form;
          var FCValidation    = formController.get('validation');
          if ( FCValidation  !== undefined ) {
              var attr = this.get('attr') ;
              var valideStruct =  Ember.validators.validate(attr);
              console.log("valideStruct",valideStruct);

              this.$().closest('div').next(".help-block").remove();

              if (!valideStruct.valid) {

                this.$().closest('div').addClass('has-error').after("<span class='help-block'>"+ valideStruct.error + "</span>");
              } else {

                this.$().closest('div').removeClass("has-error");
              }

              return valideStruct.valid;
          }
        }
    });

    void ($);
});

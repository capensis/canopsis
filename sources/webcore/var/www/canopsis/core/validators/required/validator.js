define([
], function() {

  function requiredValidator(attr, valideStruct) {

    console.log("requiredValidator :attr = ",attr) ;
    if (attr.model.options.required !== undefined && attr.model.options.required === true  && (attr.value === undefined || attr.value === ""  || attr.value === null)) {
        valideStruct.valid = false;
        valideStruct.error = " Field can't be empty";
    } else {

        valideStruct.valid = true;
        valideStruct.error = "";
    }

    return valideStruct;

  };

  return requiredValidator;
});

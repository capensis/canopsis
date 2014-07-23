define([
], function() {

  function mailValidator(attr, valideStruct)
  {
    var regex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    if (regex.test(attr.value))
      {
          valideStruct.valid = true ;
      }
    else
      {
          valideStruct.valid = false ;
          valideStruct.error = "Mail's format should be: X@Y.Z";
      }

    return valideStruct;
  };

  return mailValidator;
});
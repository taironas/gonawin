'use strict'

angular.module('directive.formValidation', []).
  directive('emails', function () {
    return {
      require: 'ngModel',
      link: function(scope, elem, attr, ctrl) {
	elem.on('blur', function (evt) {
          scope.$apply(function () {
	    var valid = true;
	    var re = /^[\S]+@[\S]+\.[\S]+$/;
	    emails = elem.val().split(',')
	    for(i in emails) {
	      if (!re.test(emails[i].trim())) {
		valid = false;
		break;
	      }
	    }
	    ctrl.$setValidity('emails', valid);
	  });
	});
      }
    }
  });

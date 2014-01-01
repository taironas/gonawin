angular.module('directive.formValidation', []).
directive('emails', function () {
	return {
		require: 'ngModel',
		link: function(scope, elem, attr, ngModel) {
				console.log('attr = ', attr);
				console.log('ngModel = ', ngModel);
				var emails = attr.emails.split(',');
	
				function validate(value) {
					console.log('emails = ', emails);
					var valid = false;
					ngModel.$setValidity('emails', valid);
					return valid ? value : undefined;
				}
	
				//For view -> model validation
				ngModel.$parsers.unshift(function(value) {
					var valid = false;
					ngModel.$setValidity('emails', valid);
					return valid ? value : undefined;
				});
	
				//For model -> view validation
				ngModel.$formatters.unshift(function(value) {
					ngModel.$setValidity('emails', false);
					return value;
				});
		}
	};
});
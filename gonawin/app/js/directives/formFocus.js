'use strict'

angular.module('directive.formFocus', []).
  directive('focus', function () {
    return function (scope, element, attrs) {
      attrs.$observe('focus', function (newValue) {
        newValue === 'true' && element[0].focus();
      });
    }
  });

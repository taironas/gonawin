'use strict'

angular.module('directive.addButton', []).directive('addbutton', function() {
  return {
    restrict: 'E',
    template: '<button class="btn btn-default" ng-click="action()">{{name}}</button>',
    replace: true,
    scope: {
      action: '&',
      name: '=',
      id: '='
    }
  };
});

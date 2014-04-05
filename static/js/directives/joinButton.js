'use strict'

angular.module('directive.joinButton', []).directive('joinbutton', function() {
  return {
    restrict: 'E',
    template: '<button class="btn btn-default" ng-click="action()">{{name}}</button>',
    replace: true,
    scope: {
      action: '&',
      name: '=',
      teamid: '='
    }
  };
});

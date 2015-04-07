'use strict';

angular.module('directive.joinButton', []).directive('joinbutton', function() {
  return {
    restrict: 'E',
    template: '<button class="btn btn-default btn-sm" ng-click="action()"><strong>{{name}}</strong></button>',
    replace: true,
    scope: {
      action: '&',
      name: '=',
      teamid: '='
    }
  };
});

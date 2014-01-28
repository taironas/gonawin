'use strict'

angular.module('directive.joinButton', []).directive('joinbutton', function() {
    return {
      restrict: 'E',
      template: '<button class="btn btn-primary" ng-click="action()">{{name}}</button>',
      scope: {
        action: '&',
        name: "="
      },
      replace: true
    };
  }
);
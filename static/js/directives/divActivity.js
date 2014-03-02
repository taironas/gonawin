'use strict'

angular.module('directive.divActivity', []).directive('gwActivity', function() {
  return {
    restrict: 'E',
    scope: {
      value: '=',
    },
    template: '<td><a href="/ng#/users/show/{{value.Actor.ID}}">{{value.Actor.DisplayName}}</a> {{value.Verb}} <a href="/ng#/teams/show/{{value.Object.ID}}">{{value.Object.DisplayName}}</a></td>',
    link: function(scope, element, attr) {
      console.log('link, scope = ', scope);
    }
  };
});
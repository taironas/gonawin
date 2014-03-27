'use strict'

angular.module('directive.divActivity', []).directive('gwActivity', function() {
  return {
    restrict: 'E',
    scope: {
      value: '=',
    },
    template: '<td><a href="/ng#/users/{{value.Actor.Id}}">{{value.Actor.DisplayName}}</a> {{value.Verb}} <a href="/ng#/teams/{{value.Object.Id}}">{{value.Object.DisplayName}}</a></td>',
    link: function(scope, element, attr) {
      console.log('link, scope = ', scope);
    }
  };
});

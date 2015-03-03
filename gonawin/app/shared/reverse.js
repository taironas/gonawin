'use strict'

// Use this filter to reverse an array. Can be use for ng-repeat as follows:
// <div ng-repeat="friend in friends | reverse">{{friend.name}}</div>
angular.module('filter.reverse', []).filter('reverse', function() {
  return function(items) {
    return items.slice().reverse();
  };
});

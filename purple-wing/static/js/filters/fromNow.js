'use strict'

angular.module('filter.fromNow', []).filter('fromNow', function() {
  return function(dateString) {
    return moment(dateString).fromNow();
  };
});
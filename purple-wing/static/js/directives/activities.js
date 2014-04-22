'use strict'

angular.module('directive.activities', []).directive('gwActivities', function() {
  return {
    controller: 'ActivitiesCtrl',
    templateUrl: 'templates/directives/activities.html'
  };
});

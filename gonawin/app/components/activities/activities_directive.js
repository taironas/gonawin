'use strict'

angular.module('directive.activities', []).directive('gwActivities', function() {
  return {
    controller: 'ActivitiesCtrl',
    templateUrl: 'app/components/activities/activities.html'
  };
});

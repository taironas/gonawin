'use strict';

var activitiesControllers = angular.module('activitiesControllers', []);

activitiesControllers.controller('ActivitiesCtrl', ['$scope', '$location', 'Activity', function($scope, $location, Activity) {
  console.log("activities controller");
  $scope.activities = Activity.query();
}]);

activitiesControllers.directive('gwActivities', function() {
  return {
    templateUrl: 'templates/directives/activities.html'
  };
});

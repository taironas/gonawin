'use strict';

var activitiesControllers = angular.module('activitiesControllers', []);

activitiesControllers.controller('ActivitiesCtrl', ['$scope', '$location', 'Activity', function($scope, $location, Activity) {
  console.log("activities controller");
  $scope.activities = Activity.query();
  
  $scope.activities.$promise.then(function(result){
    if(!$scope.activities || ($scope.activities && !$scope.activities.length))
      $scope.noActivitiesMessage = 'You have no activity';
  });
}]);

activitiesControllers.directive('gwActivities', function() {
  return {
    templateUrl: 'templates/directives/activities.html'
  };
});

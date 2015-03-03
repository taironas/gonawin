'use strict';

var activitiesControllers = angular.module('activitiesControllers', []);

activitiesControllers.controller('ActivitiesCtrl', ['$scope', '$location', 'Activity', function($scope, $location, Activity) {
  console.log("activities controller");
  // Fetch activities based on the count and page variables.
  // Concatenate new activities when more button is clicked.
  $scope.loadActivities = function()
  {
    Activity.query({ count:$scope.count, page:$scope.page}).$promise.then(function(response){
      $scope.activities = $scope.activities.concat(response);
      console.log('loadActivities: length = ', response.length);
      $scope.more = response.length === $scope.count;
    });
  }
  // Indicates if there more activities that could be loaded
  $scope.hasMore = function() {
    return $scope.more;
  }
  // Triggers the loading of new activities
  $scope.showMore = function() {
    $scope.page += 1;
    $scope.loadActivities();
  }
  
  $scope.count = 20;  // number of activities per page
  $scope.page = 1;    // current page
  $scope.activities = [];
  $scope.more = true;
  $scope.loadActivities();
}]);

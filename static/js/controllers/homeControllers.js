'use strict';

var homeControllers = angular.module('homeControllers', []);

homeControllers.controller('HomeCtrl', ['$scope', '$location', 'Activity', function($scope, $location, Activity) {
  console.log("home controller");
  $scope.activities = Activity.query();
  
  $scope.activities.$promise.then(function(result){
    if(!$scope.activities || ($scope.activities && !$scope.activities.length))
      $scope.noActivitiesMessage = 'You have no activity';
  });
}]);

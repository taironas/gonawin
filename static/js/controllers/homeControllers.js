'use strict';

var homeControllers = angular.module('homeControllers', []);

homeControllers.controller('HomeCtrl', ['$scope', function($scope, $location){
  console.log("home controller");
}]);

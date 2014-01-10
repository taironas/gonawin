'use strict';

var mainControllers = angular.module('mainControllers', []);

mainControllers.controller('MainCtrl', ['$scope', function($scope, $location){
	console.log("main controller");
}]);

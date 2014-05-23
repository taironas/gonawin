'use strict';

// Root controllers manage data at root '/'
// For now only messageInfo notification is handled here.
var rootControllers = angular.module('rootControllers', []);

rootControllers.controller('RootCtrl', ['$rootScope', '$scope', '$location', 'Tournament', function($rootScope, $scope, $location, Tournament) {
  console.log("root controller");
  // get message info from redirects.
  $scope.messageInfo = $rootScope.messageInfo;
  // reset to nil var message info in root scope.
  $rootScope.messageInfo = undefined;
  
  // fetch World Cup tournament
  $scope.tournamentData =  Tournament.getWorldCup();
}]);

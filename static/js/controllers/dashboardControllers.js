'use strict';

// Dashboard controllers update the dashboard information depending on the user location.
var dashboardControllers = angular.module('dashboardControllers', []);

dashboardControllers.controller('DashboardCtrl', ['$scope', '$rootScope', '$location', '$cookieStore', 'Session', 'sAuth', 
  function ($scope, $rootScope, $location, $cookieStore, Session, sAuth) {
    console.log('Dashboard module');
    console.log('DashboardCtrl, current user = ', $rootScope.currentUser);
    console.log('DashboardCtrl, isLoggedIn = ', $rootScope.isLoggedIn);
    $scope.DashboardUserScore = 'Lima la Asno';
}]);

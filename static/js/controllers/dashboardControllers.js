'use strict';

// Dashboard controllers update the dashboard information depending on the user location.
var dashboardControllers = angular.module('dashboardControllers', []);

dashboardControllers.controller('DashboardCtrl', ['$scope', '$rootScope', '$location', '$cookieStore', 'Session', 'sAuth', function ($scope, $rootScope, $location, $cookieStore, Session, sAuth) {
  console.log('Dashboard module');
  console.log('DashboardCtrl, current user = ', $rootScope.currentUser);
  console.log('DashboardCtrl, isLoggedIn = ', $rootScope.isLoggedIn);
  $scope.dashboard = {};

  $scope.dashboard.User = 'Lima la Asno';
  console.log('location - abs url: ', $location.absUrl());
  console.log('location - url: ', $location.url());
  console.log('location - path: ', $location.path());

  // set dashboard loaction:
  var url = $location.url();
  if(url.match('^/$') != null){
    $scope.dashboard.location = 'root';
  } else if(url.match('^/tournaments/*') != null){
    $scope.dashboard.location = 'tournament';
  } else if(url.match('^/teams/*') != null){
    $scope.dashboard.location = 'team';
  }

  // event triggered when url changes. 
  // note: this event is not triggered when page is refreshed.
  $scope.$on('$locationChangeSuccess', function(event) {
    console.log('location changed:');
    console.log('location - abs url: ', $location.absUrl());
    console.log('location - url: ', $location.url());
    console.log('location - path: ', $location.path());
    var url = $location.url();
    console.log('match root? ', url.match('^/$') != null);
    console.log('match tournaments? ', url.match('^/tournaments/*') != null);
    console.log('match teams? ', url.match('^/teams/*') != null);
    if(url.match('^/$') != null){
      $scope.dashboard.location = 'root';
    } else if(url.match('^/tournaments/*') != null){
      $scope.dashboard.location = 'tournament';
    } else if(url.match('^/teams/*') != null){
      $scope.dashboard.location = 'team';
    }
  });

}]);

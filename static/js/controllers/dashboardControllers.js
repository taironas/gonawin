'use strict';

// Dashboard controllers update the dashboard information depending on the user location.
var dashboardControllers = angular.module('dashboardControllers', []);

dashboardControllers.controller('DashboardCtrl', ['$scope', '$rootScope', '$location', '$cookieStore', 'Session', 'sAuth', 'User', function ($scope, $rootScope, $location, $cookieStore, Session, sAuth, User) {
  console.log('Dashboard module');
  console.log('DashboardCtrl, current user = ', $rootScope.currentUser);
  console.log('DashboardCtrl, isLoggedIn = ', $rootScope.isLoggedIn);
  $scope.dashboard = {};
  
  // Set the dashboard data with respect of the url:
  // Current urls watches:
  // - root: /
  // - team: /teams/:id/*
  // - tournament: /tournaments/:id/*
  $scope.setDashboard = function(){
    var ctx = 'set dashboard: ';
    console.log(ctx, 'location - abs url: ', $location.absUrl());
    console.log(ctx, 'location - url: ', $location.url());
    console.log(ctx, 'location - path: ', $location.path());

    var url = $location.url();
    console.log(ctx, 'match root? ', url.match('^/$') != null);
    console.log(ctx, 'match tournaments? ', url.match('^/tournaments/[0-9]+.*') != null);
    console.log(ctx, 'match teams? ', url.match('^/teams/[0-9]+.*') != null);
    
    if(url.match('^/$') != null){
      $scope.dashboard.location = 'root';
      $rootScope.currentUser.$promise.then(function(currentUser){
	$scope.dashboard.user = currentUser.User.Name;
	$scope.dashboard.ntournaments = currentUser.User.TournamentIds.length;
	$scope.dashboard.nteams = currentUser.User.TeamIds.length;
	// get user score information:
	$scope.dashboard.score = currentUser.User.Score;
      });
    } else if(url.match('^/tournaments/[0-9]+.*') != null){
      $scope.dashboard.location = 'tournament';
    } else if(url.match('^/teams/[0-9]+.*') != null){
      $scope.dashboard.location = 'team';
    } else{
      $scope.dashboard.location = 'default';
    }
  };
  
  // set dashboard with respect to url.
  $scope.setDashboard();

  // $locationChangeSuccess event is triggered when url changes.
  // note: this event is not triggered when page is refreshed.
  $scope.$on('$locationChangeSuccess', function(event) {
    console.log('location changed:');
    $scope.setDashboard();
  });

}]);

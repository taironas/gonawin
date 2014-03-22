'use strict';

// Dashboard controllers update the dashboard information depending on the user location.
var dashboardControllers = angular.module('dashboardControllers', []);

dashboardControllers.controller('DashboardCtrl', ['$scope', '$rootScope', '$routeParams', '$location', '$cookieStore', 'Session', 'sAuth', 'User', 'Team', '$route', function ($scope, $rootScope, $routeParams, $location, $cookieStore, Session, sAuth, User, Team, $route) {
  console.log('Dashboard module');
  console.log('DashboardCtrl, current user = ', $rootScope.currentUser);
  console.log('DashboardCtrl, isLoggedIn = ', $rootScope.isLoggedIn);
  $scope.dashboard = {};
  
  // Set the dashboard data with respect of the url:
  // Current urls watches:
  // - root: /
  // - team: /teams/:id/*
  // - tournament: /tournaments/:id/*
  //
  // the context variable has the current state of the dashboard.
  // context can take values 'user', 'team', 'tournament', 'default'.
  $scope.setDashboard = function($routeParams){
    var ctx = 'set dashboard: ';
    console.log(ctx, 'location - abs url: ', $location.absUrl());
    console.log(ctx, 'location - url: ', $location.url());
    console.log(ctx, 'location - path: ', $location.path());

    var url = $location.url();
    console.log(ctx, 'match root? ', url.match('^/$') != null);
    console.log(ctx, 'match tournaments? ', url.match('^/tournaments/[0-9]+.*') != null);
    console.log(ctx, 'match teams? ', url.match('^/teams/[0-9]+.*') != null);

    // reset dashboard before any set.
      $scope.dashboard = {};
    if(url.match('^/$') != null){
      $scope.dashboard.location = 'root';
      $scope.dashboard.context = 'user';
      $rootScope.currentUser.$promise.then(function(currentUser){
	$scope.dashboard.user = currentUser.User.Name;
	$scope.dashboard.ntournaments = currentUser.User.TournamentIds.length;
	$scope.dashboard.nteams = currentUser.User.TeamIds.length;
	// get user score information:
	$scope.dashboard.score = currentUser.User.Score;
      });
    } else if(url.match('^/tournaments/[0-9]+.*') != null){
      $scope.dashboard.location = 'tournament';
      $scope.dashboard.context = 'tournament';

    } else if(url.match('^/teams/[0-9]+.*') != null){
      $scope.dashboard.location = 'team';
      $scope.dashboard.context = 'team';
      // We are not able to get the current team id with $routeParams, so we use $route instead.
      // from angular documentation: 
      // http://docs.angularjs.org/api/ngRoute/service/$routeParams
      //   Note that the $routeParams are only updated after a route change completes successfully. 
      //   This means that you cannot rely on $routeParams being correct in route resolve functions. 
      //   Instead you can use $route.current.params to access the new route's parameters.
      console.log(ctx, 'route', $route);
      Team.get({ id:$route.current.params.id }).$promise.then(function(teamResult){
      	console.log(ctx, 'get team ', teamResult);
	$scope.dashboard.ntournaments = teamResult.Team.TournamentIds.length;
	$scope.dashboard.nmembers = teamResult.Players.length;
	$scope.dashboard.accuracy = teamResult.Team.Accuracy;
      });
      
    } else{
      $scope.dashboard.location = 'default';
      $scope.dashboard.context = 'default';
    }
  };
  
  // set dashboard with respect to url in the global controller.
  // We do this because the $locationChangeSuccess event is not triggered by a refresh.
  // As the controller is called when there is a refresh we are able to set the Dashboard with the proper information.
  $scope.setDashboard($routeParams);

  // $locationChangeSuccess event is triggered when url changes.
  // note: this event is not triggered when page is refreshed.
  $scope.$on('$locationChangeSuccess', function(event) {
    console.log('location changed:');
    $scope.setDashboard($routeParams);
  });

}]);

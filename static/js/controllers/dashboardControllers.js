'use strict';

// Dashboard controllers update the dashboard information depending on the user location.
var dashboardControllers = angular.module('dashboardControllers', []);

dashboardControllers.controller('DashboardCtrl', ['$scope', '$rootScope', '$routeParams', '$location', '$cookieStore', 'Session', 'sAuth', 'User', 'Team', 'Tournament', '$route', function ($scope, $rootScope, $routeParams, $location, $cookieStore, Session, sAuth, User, Team, Tournament, $route) {
  console.log('Dashboard module');
  console.log('DashboardCtrl, current user = ', $rootScope.currentUser);
  console.log('DashboardCtrl, isLoggedIn = ', $rootScope.isLoggedIn);
  $scope.dashboard = {};
  $scope.predicate = 'Score';
  
  // Set the dashboard data with respect of the url:
  // Current urls watches:
  // - root: /
  // - team: /teams/:id/*
  // - tournament: /tournaments/:id/*
  //
  // the context variable has the current state of the dashboard.
  // context can take values 'user', 'team', 'tournament', 'default'.
  //
  // On RouteParams:
  //
  //   We are not able to get the current team id with $routeParams, so we use $route instead.
  //   from angular documentation: 
  //   http://docs.angularjs.org/api/ngRoute/service/$routeParams
  //     Note that the $routeParams are only updated after a route change completes successfully. 
  //     This means that you cannot rely on $routeParams being correct in route resolve functions. 
  //     Instead you can use $route.current.params to access the new route's parameters.
  //
  $scope.setDashboard = function(){
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
	$scope.dashboard.id = currentUser.User.Name;
	if (currentUser.User.TournamentIds){
	  $scope.dashboard.ntournaments = currentUser.User.TournamentIds.length;
	} else {
	  $scope.dashboard.ntournaments = 0;
	}
	
	if (currentUser.User.TeamIds){
	  $scope.dashboard.nteams = currentUser.User.TeamIds.length;
	}else{
	   $scope.dashboard.nteams = 0;
	}
	
	// get user score information:
	$scope.dashboard.score = currentUser.User.Score;
      });
    } else if(url.match('^/tournaments/[0-9]+.*') != null){

      $scope.dashboard.location = 'tournament';
      $scope.dashboard.context = 'tournament';

      console.log(ctx, 'route', $route);
      Tournament.get({ id:$route.current.params.id }).$promise.then(function(tournamentResult){
      	console.log(ctx, 'get tournament ', tournamentResult);
	$scope.dashboard.tournament = tournamentResult.Tournament.Name;
	if(tournamentResult.Participants){
	  $scope.dashboard.nparticipants = tournamentResult.Participants.length;
	} else{
	  $scope.dashboard.nparticipants = 0;
	}
	if(tournamentResult.Teams){
	  $scope.dashboard.nteams = tournamentResult.Teams.length;
	} else{
	  $scope.dashboard.nteams = 0;
	}

      });
      $scope.dashboard.rank = {};

      Tournament.ranking({id:$route.current.params.id, rankby:'users', limit:'10'}).$promise.then(function(rankResult){
	console.log(ctx, 'get users ranking ', rankResult);
	$scope.dashboard.rank.users = rankResult.Users;
      });

      Tournament.ranking({id:$route.current.params.id, rankby:'teams', limit:'10'}).$promise.then(function(rankResult){
	console.log(ctx, 'get teams ranking', rankResult);
	$scope.dashboard.rank.teams = rankResult.Teams;
      });

      $rootScope.currentUser.$promise.then(function(currentUser){
	$scope.dashboard.user = currentUser.User.Name;
	$scope.dashboard.id = currentUser.User.Id;
      });

    } else if(url.match('^/teams/[0-9]+.*') != null){

      $scope.dashboard.location = 'team';
      $scope.dashboard.context = 'team';

      $rootScope.currentUser.$promise.then(function(currentUser){
	$scope.dashboard.user = currentUser.User.Name;
	$scope.dashboard.id = currentUser.User.Id;
      });

      console.log(ctx, 'route', $route);
      Team.get({ id:$route.current.params.id }).$promise.then(function(teamResult){
      	console.log(ctx, 'get team ', teamResult);
	$scope.dashboard.team = teamResult.Team.Name;

	if(teamResult.Team.TournamentIds){
	  $scope.dashboard.ntournaments = teamResult.Team.TournamentIds.length;
	}else{
	  $scope.dashboard.ntournaments = 0;
	}

	if(teamResult.Players){
	  $scope.dashboard.nmembers = teamResult.Players.length;
	}else{
	  $scope.dashboard.nmembers = 0;
	}

	$scope.dashboard.accuracy = teamResult.Team.Accuracy;
      });

      Team.ranking({id:$route.current.params.id, limit:'10'}).$promise.then(function(rankResult){
	console.log(ctx, 'get team ranking', rankResult);
	$scope.dashboard.members = rankResult.Users;
      });

      
    } else {
      $scope.dashboard.location = 'default';
      $scope.dashboard.context = 'default';
    }
  };
  
  // set dashboard with respect to url in the global controller.
  // We do this because the $locationChangeSuccess event is not triggered by a refresh.
  // As the controller is called when there is a refresh we are able to set the Dashboard with the proper information.
  $scope.setDashboard();

  // $locationChangeSuccess event is triggered when url changes.
  // note: this event is not triggered when page is refreshed.
  $scope.$on('$locationChangeSuccess', function(event) {
    console.log('location changed:');
    $scope.setDashboard();
  });

}]);

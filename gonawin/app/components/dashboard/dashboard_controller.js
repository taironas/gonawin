'use strict';

// Dashboard controllers update the dashboard information depending on the user location.
var dashboardControllers = angular.module('dashboardControllers', []);

dashboardControllers.controller('DashboardCtrl', ['$scope', '$rootScope', '$routeParams', '$location', '$cookieStore', 'Session', 'sAuth', 'User', 'Team', 'Tournament', '$route',
  function ($scope, $rootScope, $routeParams, $location, $cookieStore, Session, sAuth, User, Team, Tournament, $route) {
    console.log('Dashboard module');
    console.log('DashboardCtrl, current user = ', $rootScope.currentUser);
    console.log('DashboardCtrl, isLoggedIn = ', $rootScope.isLoggedIn);
    $scope.dashboard = {};
    $scope.predicate = 'Score';
    $scope.state = '';

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
    //   We are not able to get the current team id with $routeParams all the time, so we use $route when $routeParams does not work.
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
      console.log(ctx, 'match tournaments? ', url.match('^/tournaments/[0-9]+.*') != null);
      console.log(ctx, 'match teams? ', url.match('^/teams/[0-9]+.*') != null);
      console.log(ctx, 'route:--->', $route);

      if(url.match('^/tournaments/[0-9]+.*') != null){
        // if same state as before just exit.
        if($scope.state == 'tournaments' && $scope.needToRefresh == false){
          console.log('same!');
          return;
        }
        $scope.state = 'tournaments';

        // reset dashboard before getting data
        $scope.dashboard = {};
        $scope.dashboard.context = 'tournament';

        // depending on how we get here, either by refresh, redirect or following link,
        // sometimes routeparams will have the tournament id and other times route.current will have it.
        var tournamentId;
        if($route.current.params.id != undefined){
          tournamentId = $route.current.params.id;
        } else if($routeParams.id != undefined){
          tournamentId = $routeParams.id;
        } else {
          console.log(ctx, 'unable to get tournament id.');
          return;
        }

        Tournament.get({ id:tournamentId }).$promise.then(function(tournamentResult){
          console.log(ctx, 'get tournament ', tournamentResult);
          $scope.dashboard.tournament = tournamentResult.Tournament.Name;
          $scope.dashboard.tournamentid = tournamentResult.Tournament.Id;
          $scope.dashboard.id = tournamentResult.Tournament.Id;
          $scope.dashboard.name = tournamentResult.Tournament.Name;

          $scope.imageURL = tournamentResult.ImageURL;

          if(tournamentResult.Participants){
            $scope.dashboard.nparticipants = tournamentResult.Participants.length;
          } else {
            $scope.dashboard.nparticipants = 0;
          }
          if(tournamentResult.Teams){
            $scope.dashboard.nteams = tournamentResult.Teams.length;
          } else {
            $scope.dashboard.nteams = 0;
          }
        });

        $scope.dashboard.rank = {};
        $scope.rankBy = 'users';
        Tournament.ranking({id:tournamentId, rankby:'users', limit:'10'}).$promise.then(function(rankResult){
          console.log(ctx, 'get users ranking ', rankResult);
          $scope.dashboard.rank.users = rankResult.Users;
        });
        Tournament.ranking({id:tournamentId, rankby:'teams', limit:'10'}).$promise.then(function(rankResult){
          console.log(ctx, 'get teams ranking', rankResult);
          $scope.dashboard.rank.teams = rankResult.Teams;
        });

        $rootScope.currentUser.$promise.then(function(currentUser){
          $scope.dashboard.user = currentUser.displayName();
          $scope.dashboard.userid = currentUser.User.Id;
          $scope.dashboard.id = currentUser.User.Id;
        });
      } else if(url.match('^/tournaments/?$') != null) {
        // if same state as before just exit.
        if($scope.state == 'tournamentsindex' && $scope.needToRefresh == false) {
          console.log('same!');
          return;
        }
        $scope.state = 'tournamentsindex';

        // reset dashboard before getting data
        $scope.dashboard = {};
        $scope.dashboard.context = 'tournaments index';

        $rootScope.currentUser.$promise.then(function(currentUser){
          $scope.dashboard.user = currentUser.displayName();
	  $scope.dashboard.name = currentUser.displayName();
          $scope.dashboard.userid = currentUser.User.Id;
          $scope.dashboard.id = currentUser.User.Id;

          $scope.imageURL = currentUser.ImageURL;
	  

          if (currentUser.User.TournamentIds){
            $scope.dashboard.ntournaments = currentUser.User.TournamentIds.length;
          } else {
            $scope.dashboard.ntournaments = 0;
          }

          if (currentUser.User.TeamIds){
            $scope.dashboard.nteams = currentUser.User.TeamIds.length;
          } else {
            $scope.dashboard.nteams = 0;
          }
          // get user score information:
          $scope.dashboard.score = currentUser.User.Score;
          // get user tournaments.
          User.tournaments({id:$rootScope.currentUser.User.Id}).$promise.then(function(response){
            $scope.dashboard.tournaments = response.Tournaments;
          });
        });
      } else if(url.match('^/teams/[0-9]+.*') != null) {
        // if same state as before just exit.
        if($scope.state == 'teams' && $scope.needToRefresh == false){
          console.log('same!');
          return;
        }
        $scope.state = 'teams';

        // reset dashboard before getting data
        $scope.dashboard = {};
        $scope.dashboard.context = 'team';

        var teamId;
        if($route.current.params.id != undefined){
          teamId = $route.current.params.id;
        } else if($routeParams.id != undefined) {
          teamId = $routeParams.id;
        } else {
          console.log(ctx, 'unable to get tournament id.');
          return;
        }

        $rootScope.currentUser.$promise.then(function(currentUser) {
          $scope.dashboard.user = currentUser.displayName();
          $scope.dashboard.userid = currentUser.User.Id;
          $scope.dashboard.id = currentUser.User.Id;
        });

        console.log(ctx, 'route', $route);
        Team.get({ id:teamId }).$promise.then(function(teamResult) {
          console.log(ctx, 'get team ', teamResult);
          $scope.dashboard.team = teamResult.Team.Name;
          $scope.dashboard.teamid = teamResult.Team.Id;
          $scope.dashboard.id = teamResult.Team.Id;
          $scope.dashboard.name = teamResult.Team.Name;
          
	  $scope.imageURL = teamResult.ImageURL;

          if(teamResult.Team.TournamentIds){
            $scope.dashboard.ntournaments = teamResult.Team.TournamentIds.length;
          } else {
            $scope.dashboard.ntournaments = 0;
          }

          if(teamResult.Players){
            $scope.dashboard.nmembers = teamResult.Players.length;
          } else {
            $scope.dashboard.nmembers = 0;
          }

          $scope.dashboard.accuracy = teamResult.Team.Accuracy;
        });

        Team.ranking({id:teamId, limit:'10'}).$promise.then(function(rankResult){
          console.log(ctx, 'get team ranking', rankResult);
          $scope.dashboard.members = rankResult.Users;
        });
      } else if(url.match('^/teams/?$') != null) {
        // if same state as before just exit.
        if($scope.state == 'teamsindex' && $scope.needToRefresh == false) {
          console.log('same!');
          return;
        }
        $scope.state = 'teamsindex';

        // reset dashboard before getting data
        $scope.dashboard = {};
        $scope.dashboard.context = 'teams index';

        $rootScope.currentUser.$promise.then(function(currentUser) {
          $scope.dashboard.user = currentUser.displayName();
          $scope.dashboard.name = currentUser.displayName();
          $scope.dashboard.userid = currentUser.User.Id;
          $scope.dashboard.id = currentUser.User.Id;

	  $scope.imageURL = currentUser.ImageURL;

          if (currentUser.User.TournamentIds){
            $scope.dashboard.ntournaments = currentUser.User.TournamentIds.length;
          } else {
            $scope.dashboard.ntournaments = 0;
          }

          if (currentUser.User.TeamIds){
            $scope.dashboard.nteams = currentUser.User.TeamIds.length;
          } else {
            $scope.dashboard.nteams = 0;
          }
          // get user score information:
          $scope.dashboard.score = currentUser.User.Score;

          // get user teams.
          User.teams({id:$rootScope.currentUser.User.Id}).$promise.then(function(response){
            $scope.dashboard.teams = response.Teams;
          });
        });
      // Default dashboard: user
      } else {
        // if same state as before just exit.
        if($scope.state == 'user' && $scope.needToRefresh == false) {
          console.log('same!');
          return;
        }

        $scope.state = 'user';
        // reset dashboard before getting data
        $scope.dashboard = {};
        $scope.dashboard.context = 'user';

        $rootScope.currentUser.$promise.then(function(currentUser){
          $scope.dashboard.user = currentUser.User.Name;
          $scope.dashboard.name = currentUser.displayName();
          $scope.dashboard.userid = currentUser.User.Id;
          $scope.dashboard.id = currentUser.User.Id;

	  $scope.imageURL = currentUser.ImageURL;

          if (currentUser.User.TournamentIds){
            $scope.dashboard.ntournaments = currentUser.User.TournamentIds.length;
          } else {
            $scope.dashboard.ntournaments = 0;
          }

          if (currentUser.User.TeamIds){
            $scope.dashboard.nteams = currentUser.User.TeamIds.length;
          } else{
            $scope.dashboard.nteams = 0;
          }
          // get user score information:
          $scope.dashboard.score = currentUser.User.Score;
        });
      }
    };

    $scope.byUsersOnClick = function(){
      $scope.rankBy = 'users';
      return;
    };

    $scope.byTeamsOnClick = function(){
      $scope.rankBy = 'teams';
      return;
    };

    // set dashboard with respect to url in the global controller.
    // We do this because the $locationChangeSuccess event is not triggered by a refresh.
    // As the controller is called when there is a refresh we are able to set the Dashboard with the proper information.
    $scope.needToRefresh = false;
    $scope.setDashboard();

    // $locationChangeSuccess event is triggered when url changes.
    // note: this event is not triggered when page is refreshed.
    $scope.$on('$locationChangeSuccess', function(event) {
      console.log('location changed:');
      $scope.needToRefresh = false;
      $scope.setDashboard();
    });

    $scope.$on('setUpdatedDashboard', function(event, refresh) {
      $scope.needToRefresh = true;
      $scope.setDashboard();
    });
}]);

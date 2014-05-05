'use strict';

// Dashboard controllers update the dashboard information depending on the user location.
var dashboardControllers = angular.module('dashboardControllers', []);

dashboardControllers.controller('DashboardCtrl', ['$scope', '$rootScope', '$routeParams', '$location', '$cookieStore', 'Session', 'sAuth', 'User', 'Team', 'Tournament', '$route', function ($scope, $rootScope, $routeParams, $location, $cookieStore, Session, sAuth, User, Team, Tournament, $route) {
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
	console.log(ctx, 'match root? ', url.match('^/$') != null);
	console.log(ctx, 'match tournaments? ', url.match('^/tournaments/[0-9]+.*') != null);
	console.log(ctx, 'match teams? ', url.match('^/teams/[0-9]+.*') != null);
	console.log(ctx, 'route:--->', $route);

	if(url.match('^/$') != null || url.match('^/users/[0-9]+.*') != null){
	    $scope.state = 'root';

	    // reset dashboard before getting data
	    $scope.dashboard = {};

	    $scope.dashboard.location = 'root';
	    $scope.dashboard.context = 'user';
	    $rootScope.currentUser.$promise.then(function(currentUser){
		$scope.dashboard.user = currentUser.User.Name;
		$scope.dashboard.name = currentUser.User.Username;
		$scope.dashboard.userid = currentUser.User.Id;
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
	    // if same state as before just exit.
	    if($scope.state == 'tournaments'){
		console.log('same!');
		return;
	    }
	    // reset dashboard before getting data
	    $scope.dashboard = {};

	    $scope.state = 'tournaments';
	    $scope.dashboard.location = 'tournament';
	    $scope.dashboard.context = 'tournament';

	    // depending on how we get here, either by refresh, redirect or following link,
	    // sometimes routeparams will have the tournament id and other times route.current will have it.
	    var tournamentId;
	    if($route.current.params.id != undefined){
		tournamentId = $route.current.params.id;
	    } else if($routeParams.id != undefined){
		tournamentId = $routeParams.id;
	    } else{
		console.log(ctx, 'unable to get tournament id.');
		return;
	    }

	    Tournament.get({ id:tournamentId }).$promise.then(function(tournamentResult){
      		console.log(ctx, 'get tournament ', tournamentResult);
		$scope.dashboard.tournament = tournamentResult.Tournament.Name;
		$scope.dashboard.tournamentid = tournamentResult.Tournament.Id;
    $scope.dashboard.name = tournamentResult.Tournament.Name;

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
		$scope.dashboard.user = currentUser.User.Name;
		$scope.dashboard.userid = currentUser.User.Id;
	    });
	} else if(url.match('^/tournaments/?$') != null){

	    // if same state as before just exit.
	    if($scope.state == 'tournamentsindex'){
		console.log('same!');
		return;
	    }
	    // reset dashboard before getting data
	    $scope.dashboard = {};

	    $scope.dashboard.location = 'tournaments index';
	    $scope.dashboard.context = 'tournaments index';
	    $scope.state = 'tournamentsindex';

	    $rootScope.currentUser.$promise.then(function(currentUser){
		$scope.dashboard.user = currentUser.User.Name;
		$scope.dashboard.name = currentUser.User.Name;
		$scope.dashboard.userid = currentUser.User.Id;
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
		// get user tournaments.
		User.tournaments({id:$rootScope.currentUser.User.Id}).$promise.then(function(response){
		    $scope.dashboard.tournaments = response.Tournaments;
		});

	    });

	} else if(url.match('^/teams/[0-9]+.*') != null){
	    $scope.state = 'teams';

	    // reset dashboard before getting data
	    $scope.dashboard = {};

	    $scope.dashboard.location = 'team';
	    $scope.dashboard.context = 'team';

	    var teamId;
	    if($route.current.params.id != undefined){
		teamId = $route.current.params.id;
	    } else if($routeParams.id != undefined){
		teamId = $routeParams.id;
	    } else{
		console.log(ctx, 'unable to get tournament id.');
		return;
	    }


	    $rootScope.currentUser.$promise.then(function(currentUser){
		$scope.dashboard.user = currentUser.User.Name;
		$scope.dashboard.userid = currentUser.User.Id;
	    });

	    console.log(ctx, 'route', $route);
	    Team.get({ id:teamId }).$promise.then(function(teamResult){
      		console.log(ctx, 'get team ', teamResult);
		$scope.dashboard.team = teamResult.Team.Name;
		$scope.dashboard.teamid = teamResult.Team.Id;
    $scope.dashboard.name = teamResult.Team.Name;

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

	    Team.ranking({id:teamId, limit:'10'}).$promise.then(function(rankResult){
		console.log(ctx, 'get team ranking', rankResult);
		$scope.dashboard.members = rankResult.Users;
	    });
	} else if(url.match('^/teams/?$') != null){
	    // if same state as before just exit.
	    if($scope.state == 'teamsindex'){
		console.log('same!');
		return;
	    }
	    // reset dashboard before getting data
	    $scope.dashboard = {};

	    $scope.state = 'teamsindex';
	    $scope.dashboard.location = 'teams index';
	    $scope.dashboard.context = 'teams index';

	    $rootScope.currentUser.$promise.then(function(currentUser){
		console.log('current user!!!!!!', currentUser);
		$scope.dashboard.user = currentUser.User.Name;
		$scope.dashboard.name = currentUser.User.Name;
		$scope.dashboard.userid = currentUser.User.Id;
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

		// get user teams.
		User.teams({id:$rootScope.currentUser.User.Id}).$promise.then(function(response){
		    $scope.dashboard.teams = response.Teams;
		});

	    });

	} else {
	    $scope.dashboard.location = 'default';
	    $scope.dashboard.context = 'default';
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
    $scope.setDashboard();

    // $locationChangeSuccess event is triggered when url changes.
    // note: this event is not triggered when page is refreshed.
    $scope.$on('$locationChangeSuccess', function(event) {
	console.log('location changed:');
	$scope.setDashboard();
    });

}]);

'use strict';

var tournamentControllers = angular.module('tournamentControllers', []);

tournamentControllers.controller('TournamentListCtrl', ['$scope', 'Tournament', '$location', function($scope, Tournament, $location) {
    console.log('Tournament list controller');
    $scope.tournaments = Tournament.query();
    $scope.searchTournament = function(){
      console.log('TournamentListCtrl: searchTournament');
      console.log('keywords: ', $scope.keywords)
      $location.search('q', $scope.keywords).path('/tournaments/search');
    };
}]);

tournamentControllers.controller('TournamentSearchCtrl', ['$scope', '$routeParams', 'Tournament', '$location', function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament search controller');
    console.log('routeParams: ', $routeParams);
    $scope.tournaments = Tournament.search( {q:$routeParams.q});

    $scope.query = $routeParams.q;
    $scope.isSearching = true;
    $scope.searchTournament = function(){
      console.log('TournamentSearchCtrl: searchTournament');
      console.log('keywords: ', $scope.keywords)
      $location.search('q', $scope.keywords).path('/tournaments/search');
    };
}]);

tournamentControllers.controller('TournamentNewCtrl', ['$scope', 'Tournament', '$location', function($scope, Tournament, $location) {
    console.log('Tournament New controller');

    $scope.addTournament = function() {
	Tournament.save($scope.tournament,
			function(tournament) {
        $location.path('/tournaments/show/' + tournament.Id);
			},
			function(err) {
        console.log('save failed: ', err.data);
			});
    };
}]);

tournamentControllers.controller('TournamentShowCtrl', ['$scope', '$routeParams', 'Tournament', '$location', '$q', 
  function($scope, $routeParams, Tournament, $location, $q) {
    console.log('Tournament Show controller');
    
    $scope.tournamentData =  Tournament.get({ id:$routeParams.id });
    console.log('tournamentData', $scope.tournamentData);
    
    $scope.deleteTournament = function() {
      Tournament.delete({ id:$routeParams.id },
			  function(){
          $location.path('/');
			  },
			  function(err) {
          console.log('delete failed: ', err.data);
			  });
    };

    $scope.leaveTournamentAsTeam = function(teamId){
      console.log('team leave tournament ');
      console.log(teamId);
      Tournament.leaveAsTeam({id:$routeParams.id, teamId:teamId},
        function(tournament) {
          console.log('success leave as a team');
          $location.path('/tournaments/show/' + tournament.Id);
        },
        function(err) {
          console.log('leave as a team failed: ', err.data);
        });
    };

    $scope.joinTournamentAsTeam = function(teamId){
      console.log('team join tournament ');
      console.log(teamId);
      Tournament.joinAsTeam({id:$routeParams.id, teamId:teamId},
        function(tournament) {
          console.log('success join as a team');
          $location.path('/tournaments/show/' + tournament.Id);
        },
        function(err) {
          console.log('join as a team failed: ', err.data);
      });
    };

    $scope.joinOrLeaveTournamentAsTeam = function(team){
    console.log('join or leave tournament as team');
    console.log('id', team.Id);
      if(team.Joined){
        $scope.leaveTournamentAsTeam(team.Id);
      }else{
        $scope.joinTournamentAsTeam(team.Id);
      }
    };

    $scope.isTournamentAdmin = $scope.tournamentData.$promise.then(function(result){
      console.log('tournament is admin ready!');
      if(result.Tournament.AdminId == $scope.currentUser.Id){
        return true;
      }else{
        return false;
      }
    });

    // checks if user has joined a tournament
    $scope.joined = $scope.tournamentData.$promise.then(function(result){
      console.log('tournament joined ready!');
      return result.Joined;
    });

    $scope.candidates = Tournament.candidates( {id:$routeParams.id});
    $scope.candidates.$promise.then(function(result){
      console.log('candidates ready!', result);
    });
    
    $scope.joinTournament = function(){
      Tournament.join({ id:$routeParams.id }).$promise.then(function(result){
        $scope.tournamentData = Tournament.get({ id:$routeParams.id });
        $scope.joinButtonName = 'Leave';
        $scope.joinButtonMethod = $scope.leaveTournament;
      });
    };

    $scope.leaveTournament = function(){
      Tournament.leave({ id:$routeParams.id }).$promise.then(function(result){
        $scope.tournamentData = Tournament.get({ id:$routeParams.id });
        $scope.joinButtonName = 'Join';
        $scope.joinButtonMethod = $scope.joinTournament;
      });
    };
    
    $scope.tournamentData.$promise.then(function(tournamentResult){
      var deferred = $q.defer();
      if (tournamentResult.Joined) {
        deferred.resolve('Leave');
      }
      else {
        deferred.resolve('Join');
      }
      return deferred.promise;
    }).then(function(result){
      $scope.joinButtonName = result;
    });
    
    $scope.tournamentData.$promise.then(function(tournamentResult){
      var deferred = $q.defer();
      if (tournamentResult.Joined) {
        deferred.resolve($scope.leaveTournament);
      }
      else {
        deferred.resolve($scope.joinTournament);
      }
      return deferred.promise;
    }).then(function(result){
      $scope.joinButtonMethod = result;
    });
}]);

tournamentControllers.controller('TournamentEditCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
    $scope.tournament = Tournament.get({ id:$routeParams.id });
    
    $scope.updateTournament = function() {
      var tournament = Tournament.get({ id:$routeParams.id });
      Tournament.update({ id:$routeParams.id }, $scope.tournament,
			  function(){
          $location.path('/tournaments/show/' + $routeParams.id);
			  },
			  function(err) {
          console.log('update failed: ', err.data);
			  });
    }
}]);

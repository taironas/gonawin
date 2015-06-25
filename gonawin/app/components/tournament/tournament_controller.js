'use strict';

// Tournament controllers manage tournament entities (creation, update, deletion) by getting
// data from REST service (resource).
// Handle also user subscription to a tournament (join/leave and join as team/leave as team).
var tournamentControllers = angular.module('tournamentControllers', []);
// TournamentListCtrl: fetch all tournaments data
tournamentControllers.controller('TournamentListCtrl', ['$scope', '$rootScope', 'Tournament', '$location',
  function($scope, $rootScope, Tournament, $location) {
    console.log('Tournament list controller:');

    $rootScope.title = 'gonawin - Tournaments';

    $scope.countTournaments = 25;  // counter for the number of tournaments to display in view.
    $scope.pageTournaments = 1;    // page counter for tournaments, to know which page to display next.

    // main query to /j/tournaments to get all tournaments.
    $scope.tournaments = Tournament.query({count:$scope.countTournaments, page:$scope.pageTournaments});

    $scope.tournaments.$promise.then(function(response){
	if(!$scope.tournaments || ($scope.tournaments && !$scope.tournaments.length)){
	    $scope.noTournamentsMessage = 'There are no tournaments yet';
	}else if($scope.tournaments !== undefined){
	    $scope.showMoreTournaments = (response.length == $scope.countTournaments);
	}
    });

    // show more tournaments function:
    // retreive tournaments by page and increment page.
    $scope.moreTournaments = function(){
	console.log('more tournaments');
	$scope.pageTournaments = $scope.pageTournaments + 1;
	Tournament.query({count:$scope.countTournaments, page:$scope.pageTournaments}).$promise.then(function(response){
	    console.log('response: ', response);
	    $scope.tournaments = $scope.tournaments.concat(response);
	    $scope.showMoreTournaments = (response.length == $scope.countTournaments);
	});
    };

    $scope.searchTournament = function(){
	console.log('TournamentListCtrl: searchTournament');
	console.log('keywords: ', $scope.keywords)
	$location.search('q', $scope.keywords).path('/search');
    };

    // start world cup create action
    $scope.createWorldCup = function() {
      console.log('Creating world cup');
      Tournament.saveWorldCup($scope.tournament, function(tournament) {
        console.log('World Cup Tournament: ', tournament);
        $location.path('/tournaments/' + tournament.Id);
      }, function(err) {
        console.log('save failed: ', err.data);
        $scope.messageDanger = err.data;
      });
    };
    // end world cup create action

    // start champions league create action
    $scope.createChampionsLeague = function() {
      console.log('Creating champions league');
      Tournament.saveChampionsLeague($scope.tournament, function(tournament) {
        console.log('Champions League Tournament: ', tournament);
        $location.path('/tournaments/' + tournament.Id);
      }, function(err) {
        console.log('save failed: ', err.data);
        $scope.messageDanger = err.data;
      });
    };

    // start copa america create action
    $scope.createCopaAmerica = function() {
      console.log('Creating copa america');
      Tournament.saveCopaAmerica($scope.tournament, function(tournament) {
        console.log('Copa America Tournament: ', tournament);
        $location.path('/tournaments/' + tournament.Id);
      }, function(err) {
        console.log('save failed: ', err.data);
        $scope.messageDanger = err.data;
      });
    };


}]);

// TournamentCardCtrl: handles team card
tournamentControllers.controller('TournamentCardCtrl', ['$rootScope', '$scope', '$q', 'Tournament',
  function($rootScope, $scope, $q, Tournament) {
    console.log('Tournament Card controller: tournament = ', $scope.tournament);

    $scope.tournamentData = Tournament.get({ id:$scope.tournament.Id });
}]);

// TournamentNewCtrl: use this controller to create a new tournament.
tournamentControllers.controller('TournamentNewCtrl', ['$rootScope', '$scope', 'Tournament', '$location', function($rootScope, $scope, Tournament, $location) {
  console.log('Tournament New controller');

  $rootScope.title = 'gonawin - New Tournament';

  $scope.addTournament = function() {
    Tournament.save($scope.tournament,
		    function(response) {
		      $rootScope.messageInfo = response.MessageInfo;
		      $location.path('/tournaments/' + response.Tournament.Id);
		    },
		    function(err) {
		      console.log('save failed: ', err.data);
		      $scope.messageDanger = err.data;
		    });
  };
}]);

// TournamentShowCtrl: fetch data of specific tournament.
// Handle also deletion of this same tournament and join/leave and join/leave as team.
tournamentControllers.controller('TournamentShowCtrl', ['$rootScope', '$scope', '$routeParams', 'Tournament', '$location', '$q', '$route', function($rootScope, $scope, $routeParams, Tournament, $location, $q, $route) {
    console.log('Tournament Show controller: Start');

    $scope.tournamentData =  Tournament.get({ id:$routeParams.id });
    console.log('tournamentData', $scope.tournamentData);

    $scope.tournamentData.$promise.then(function(response){
      $rootScope.title = 'gonawin - ' + response.Tournament.Name;
    });

    // get message info from redirects.
    $scope.messageInfo = $rootScope.messageInfo;
    // reset to nil var message info in root scope.
    $rootScope.messageInfo = undefined;

    // get candidates data from tournament id
    $scope.candidateTeamsData = Tournament.candidates({id:$routeParams.id});

    // do we really need theses lines?
    $scope.candidateTeamsData.$promise.then(function(result) {
      $scope.candidates = result.Candidates;
    });

    // list of tournament groups
    $scope.groupsData = Tournament.groups({id:$routeParams.id});
    console.log('groupsData', $scope.groupsData);
    // admin function: reset tournament
    $scope.resetTournament = function() {
      Tournament.reset({id:$routeParams.id},
        function(result) {
			     console.log('reset succeed.');
			     $scope.messageInfo = result.MessageInfo;
			     $scope.groupsData.Groups = result.Groups;
        }, function(err) {
			     console.log('reset failed: ', err.data);
			     $scope.messageDanger = err.data;
        });
    };

    $scope.deleteTournament = function() {
      if(confirm('Are you sure?')) {
        Tournament.delete({ id:$routeParams.id },
          function(response) {
            $rootScope.messageInfo = response.MessageInfo;
            $location.path('/');
          }, function(err) {
            console.log('delete failed: ', err.data);
            $scope.messageDanger = err.data;
          });
      }
    };

    $scope.joinTournament = function(){
      Tournament.join({ id:$routeParams.id }).$promise.then(function(response){
          $scope.joinButtonName = 'Leave Tournament';
          $scope.joinButtonMethod = $scope.leaveTournament;
          $scope.messageInfo = response.MessageInfo;
          $scope.$broadcast('setUpdatedTournamentData');
          $rootScope.$broadcast('setUpdatedDashboard');
      });
    };

    $scope.leaveTournament = function(){
      if(confirm('Are you sure?\nYou will not be able to play on this tournament.')){
        Tournament.leave({ id:$routeParams.id }).$promise.then(function(response){
            $scope.joinButtonName = 'Join';
            $scope.joinButtonMethod = $scope.joinTournament;
            $scope.messageInfo = response.MessageInfo;
            $scope.$broadcast('setUpdatedTournamentData');
            $rootScope.$broadcast('setUpdatedDashboard');
        });
      }
    };

    $scope.joinTournamentAsTeam = function(teamId){
    	Tournament.joinAsTeam({id:$routeParams.id, teamId:teamId}).$promise.then(function(response){
    	    $scope.joinAsTeamButtonName[teamId] = 'Leave Tournament';
    	    $scope.joinAsTeamButtonMethod[teamId] = $scope.leaveTournamentAsTeam;
    	    $scope.messageInfo = response.MessageInfo;
    	    $scope.$broadcast('setUpdatedTournamentData');
          $rootScope.$broadcast('setUpdatedDashboard');
    	});
    };

    $scope.leaveTournamentAsTeam = function(teamId){
      if(confirm('Are you sure?\nYou will not be able to play on this tournament as team.')){
        Tournament.leaveAsTeam({id:$routeParams.id, teamId:teamId}).$promise.then(function(response){
            $scope.joinAsTeamButtonName[teamId] = 'Join';
            $scope.joinAsTeamButtonMethod[teamId] = $scope.joinTournamentAsTeam;
            $scope.messageInfo = response.MessageInfo;
            $scope.$broadcast('setUpdatedTournamentData');
            $rootScope.$broadcast('setUpdatedDashboard');
        });
      }
    };

    // tab at undefined means you are in tournament/:id url.
    // So matches should be active by default
    if($routeParams.tab === undefined) {
      $scope.tab = 'matches';
    } else {
      // Initialize tab variable to handle views:
      $scope.tab = $routeParams.tab;
    }

    // set isTournamentAdmin boolean:
    // This variable defines if the current user is admin of the current tournament.
    $scope.tournamentData.$promise.then(function(tournamentResult){
      console.log('tournament is admin ready');
      // as it depends of currentUser, make a promise
      var deferred = $q.defer();
      deferred.resolve((tournamentResult.Tournament.AdminIds.indexOf($scope.currentUser.User.Id)>=0));
      return deferred.promise;
    }).then(function(result){
      console.log('is tournament admin:', result);
      $scope.isTournamentAdmin = result;
    });

    // Checks if user has joined a tournament
    $scope.joined = $scope.tournamentData.$promise.then(function(result){
	     console.log('tournament joined ready!');
       return result.Joined;
    });

    $scope.tournamentData.$promise.then(function(tournamentResult){
    	var deferred = $q.defer();
    	if (tournamentResult.Joined) {
        deferred.resolve('Leave Tournament');
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

    // set up join as a team buttons.
    $scope.candidateTeamsData.$promise.then(function(candidatesResult) {
      var candidatesLength = 0;
      if(candidatesResult.Candidates){
          candidatesLength = candidatesResult.Candidates.length;
      }
    	$scope.joinAsTeamButtonName = new Array(candidatesLength);
    	$scope.joinAsTeamButtonMethod = new Array(candidatesLength);

        $scope.tournamentData.$promise.then(function(tournamentResult){
    	    for (var i=0 ; i<candidatesLength; i++)
    	    {
        		if(IsTeamJoined(candidatesResult.Candidates[i].Team.Id, tournamentResult.Teams))
        		{
      		    $scope.joinAsTeamButtonName[candidatesResult.Candidates[i].Team.Id] = 'Leave';
      		    $scope.joinAsTeamButtonMethod[candidatesResult.Candidates[i].Team.Id] = $scope.leaveTournamentAsTeam;
        		} else {
      		    $scope.joinAsTeamButtonName[candidatesResult.Candidates[i].Team.Id] = 'Join';
      		    $scope.joinAsTeamButtonMethod[candidatesResult.Candidates[i].Team.Id] = $scope.joinTournamentAsTeam;
        		}
    	    }
    	});
    });

    // set admin candidates and array of functions.
    Tournament.participants({ id:$routeParams.id }).$promise.then(function(participantsResult) {
      console.log('set admin candidates.');
    	$scope.adminCandidates = participantsResult.Participants;
    	var len = 0;
    	if(participantsResult.Participants){
  	    len = participantsResult.Participants.length;
    	}

    	$scope.addAdminButtonName = new Array(len);
    	$scope.addAdminButtonMethod = new Array(len);
	$scope.tournamentData.$promise.then(function(result){
	    for (var i=0 ; i<len; i++){
  		// check if user is admin already here.
		if(result.Tournament.AdminIds.indexOf(participantsResult.Participants[i].Id)>=0){
  	    	    $scope.addAdminButtonName[participantsResult.Participants[i].Id] = 'Remove Admin';
  	    	    $scope.addAdminButtonMethod[participantsResult.Participants[i].Id] = $scope.removeAdmin;
		}else{
		    $scope.addAdminButtonName[participantsResult.Participants[i].Id] = 'Add Admin';
  		    $scope.addAdminButtonMethod[participantsResult.Participants[i].Id] = $scope.addAdmin;
		}
    	    }
	});
    });

    // admin modal add buttons.
    // add admin state.
    $scope.addAdmin = function(userId){
	console.log('add admin');
    	Tournament.addAdmin({id:$routeParams.id, userId:userId}).$promise.then(function(response){
    	    $scope.addAdminButtonName[userId] = 'Remove admin';
    	    $scope.addAdminButtonMethod[userId] = $scope.removeAdmin;
    	    $scope.messageInfo = response.MessageInfo;
    	}, function(err) {
	    $scope.messageDanger = err.data;
	    console.log('add admin failed: ', err.data);
	});
    };
    // remove admin state.
    $scope.removeAdmin = function(userId){
	console.log('remove admin');
    	Tournament.removeAdmin({id:$routeParams.id, userId:userId}).$promise.then(function(response){
    	    $scope.addAdminButtonName[userId] = 'Add admin';
    	    $scope.addAdminButtonMethod[userId] = $scope.addAdmin;
    	    $scope.messageInfo = response.MessageInfo;
    	},function(err){
	    console.log('remove admin failed: ', err.data);
	    $scope.messageDanger = err.data;
	});
    };

    var IsTeamJoined = function(teamId, teams) {
  	if(!teams) {
  	    return false;
  	}

  	 for (var i=0 ; i<teams.length; i++){
  	    if(teams[i].Id == teamId){
  		      return true;
  	    }
  	   }
    };

    $scope.redirectToNewTeam = false;
    // Action triggered when 'Create new button' is clicked, modal window will be hidden.
    // We also set flag 'redirectToNewTeam' to true for listener to know if redirection is needed.
    $scope.newTeam = function() {
      console.log('.addTeamModal newTeam');
      $('.addTeamModal').modal('hide');
      $scope.redirectToNewTeam = true;
    };

    // listen 'hidden.bs.modal' event to redirect to new team page
    // Only redirect if flag 'redirectToNewTeam' is set.
    $('body').on('hidden.bs.modal', '.addTeamModal', function (e) {
      console.log('.addTeamModal hidden.bs.modal');
      // need to have scope for $location to work. So add 'apply' function
      // inside js listener
      if($scope.redirectToNewTeam === true) {
        $scope.$apply(function(){
          $location.path('/teams/new/');
        });
      }
    });

  $scope.tabs = {
      "matches":         { title: 'Matches',                url: 'components/tournament/tab_matches.html' },
      "firststage":       { title: 'First Stage',             url: 'components/tournament/tab_firststage.html' },
      "secondstage":      { title: 'Second Stage',            url: 'components/tournament/tab_secondstage.html' },
      "ranking":          { title: 'Ranking',                 url: 'components/tournament/tab_ranking.html' },
      "predictions":      { title: 'Predictions',             url: 'components/tournament/tab_predictions.html' },
      "admin.setresults": { title: 'Set Results',             url: 'components/tournament/tab_setresults.html' },
      "admin.setteams":   { title: 'Set Teams',               url: 'components/tournament/tab_setteams.html' }
  };

  // set the current tab based on the 'tab' parameter
  if($scope.tab === undefined) {
    $scope.currentTab = $scope.tabs["matches"].url;
  } else {
    $scope.currentTab = $scope.tabs[$scope.tab].url;
  }

  $scope.onClickTab = function (tab) {
    $scope.currentTab = tab.url;
  };

  $scope.isCopaAmerica = false;
  $scope.tournamentData.$promise.then(function(response){
      $scope.isCopaAmerica = response.Tournament.Name == '2015 Copa America';
  });

  $scope.tournamentHasGroups = false;
  Tournament.groups({id:$routeParams.id}).$promise.then(function(response){
      $scope.tournamentHasGroups = (response.Groups.length > 0)
  });

  $scope.showStageTab = function(){
      if($scope.isCopaAmerica){
          return false;
      }
      return $scope.tournamentHasGroups;
  };

  // $locationChangeSuccess event is triggered when url changes.
  // note: this event is not triggered when page is refreshed.
  $scope.$on('$locationChangeSuccess', function(event) {
    console.log('tournament show: location changed:');
    console.log('tournament show: routeparams', $routeParams);
    // set tab parameter to new tab parameter.
    // use this to have to last state in order to display the proper view.
    if($routeParams.tab !== undefined){
      $scope.tab = $routeParams.tab;
    }
  });
}]);

// TournamentEditCtrl: collects data to update an existing tournament.
tournamentControllers.controller('TournamentEditCtrl', ['$rootScope', '$scope', '$routeParams', 'Tournament', '$location',function($rootScope, $scope, $routeParams, Tournament, $location) {
  $scope.tournamentData = Tournament.get({ id:$routeParams.id });

  $scope.tournamentData.$promise.then(function(response){
    $rootScope.title = 'gonawin - ' + response.Tournament.Name;
  });

  $scope.updateTournament = function() {
    var tournamentData = Tournament.get({ id:$routeParams.id });
    Tournament.update({ id:$routeParams.id }, $scope.tournamentData.Tournament,
		      function(response){
			$rootScope.messageInfo = response.MessageInfo;
			$location.path('/tournaments/' + $routeParams.id);
		      },
		      function(err) {
			console.log('update failed: ', err.data);
			$scope.messageDanger = err.data;
		      });
  }
}]);

// TournamentCalendarCtrl: collects complete data of specific tournament (matches, predict)
tournamentControllers.controller('TournamentCalendarCtrl', ['$scope', '$routeParams', 'Tournament', '$location', function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament calendar controller');
    console.log('route params', $routeParams);
    $scope.tournamentData = Tournament.get({ id:$routeParams.id });

    $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:$routeParams.groupby});
    console.log('matchesData = ', $scope.matchesData);

    $scope.groupby = $routeParams.groupby;

    $scope.updateMatchesView = function() {
    	if($scope.groupby !== undefined) {
	       if($scope.groupby == 'phase') {
           $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'phase'});
           $scope.matchesData.$promise.then(function(result) {
             if(result.Phases !== undefined) {
               for(var i = 0; i < result.Phases.length; i++) {
                 if($scope.matchesData.Phases[i].Completed) {
                   $scope.matchesData.Phases[i].showPhase = false;
                 } else {
                   $scope.matchesData.Phases[i].showPhase = true;
                 }
               }
             }
           });
         } else if($scope.groupby == 'date') {
           $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'date'});
         }
      } else {
        $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'phase'});
         $scope.matchesData.$promise.then(function(result) {
           if(result.Phases !== undefined) {
             for(var i = 0; i < result.Phases.length; i++) {
               if($scope.matchesData.Phases[i].Completed) {
                 $scope.matchesData.Phases[i].showPhase = false;
               } else {
                 $scope.matchesData.Phases[i].showPhase = true;
               }
             }
           }
         });
       }
     };
     $scope.updateMatchesView();

     $scope.$on('setUpdatedTournamentData', function(event) {
       $scope.tournamentData = Tournament.get({ id:$routeParams.id });
     });

    $scope.activatePredict = function(matchIdNumber, index, parentIndex) {
      console.log('Tournament calendar controller: activate predict: matchid number, index, parent index', matchIdNumber, index, parentIndex);
      console.log('Tournament calendar controller: activate predict: matchesdata', $scope.matchesData);

      $scope.matchesData.Days[parentIndex].Matches[index].wantToPredict = true;
    };

    $scope.activatePredictPhase = function(matchIdNumber, index, parentIndex, parentParentIndex) {
      console.log('Tournament calendar controller: activate predict: matchid number, index, parent index, parent parent index', matchIdNumber, index, parentIndex, parentParentIndex);
      $scope.matchesData.Phases[parentParentIndex].Days[parentIndex].Matches[index].wantToPredict = true;
    };

    $scope.predict = function(matchIdNumber, index, parentIndex, result1, result2) {
      console.log('Tournament calendar controller: predict:', matchIdNumber);
      console.log('Tournament calendar controller: predict: index:', index);
      console.log('Tournament calendar controller: predict: parent parent index:', parentIndex);

      $scope.matchesData.Days[parentIndex].Matches[index].wantToPredict = false;

      Tournament.predict({id:$routeParams.id, matchId:matchIdNumber, result1:result1, result2:result2},
        function(result) {
          console.log('success in setting prediction!');
          $scope.matchesData.Days[parentIndex].Matches[index].Predict = result.Predict.Result1 + ' - ' + result.Predict.Result2;
          $scope.messageInfo = result.MessageInfo;
          console.log('match result: ', result.Predict.Result1 + ' - ' + result.Predict.Result2);
          $scope.matchesData.Days[parentIndex].Matches[index].HasPredict = true;
        }, function(err) {
          console.log('failure setting prediction! ', err.data);
          $scope.messageDanger = err.data;
        });
      console.log('match result: ', result1, ' ', result2);
    };

    $scope.predictPhase = function(matchIdNumber, index, parentIndex, parentParentIndex, result1, result2) {
      console.log('Tournament calendar controller: predict:', matchIdNumber);
      console.log('Tournament calendar controller: predict: index:', index);
      console.log('Tournament calendar controller: predict: parent  index:', parentIndex);
      console.log('Tournament calendar controller: predict: parent parent  index:', parentParentIndex);

      $scope.matchesData.Phases[parentParentIndex].Days[parentIndex].Matches[index].wantToPredict = false;

      Tournament.predict({id:$routeParams.id, matchId:matchIdNumber, result1:result1, result2:result2},
        function(result) {
          console.log('success in setting prediction!');
          $scope.matchesData.Phases[parentParentIndex].Days[parentIndex].Matches[index].Predict = result.Predict.Result1 + ' - ' + result.Predict.Result2;
          $scope.messageInfo = result.MessageInfo;
          console.log('match result: ', result.Predict.Result1 + ' - ' + result.Predict.Result2);
          $scope.matchesData.Phases[parentParentIndex].Days[parentIndex].Matches[index].HasPredict = true;
        }, function(err) {
          console.log('failure setting prediction! ', err.data);
          $scope.messageDanger = err.data;
        });
        console.log('match result: ', result1, ' ', result2);
      };

    $scope.cancel = function(matchIdNumber, index, parentIndex, result1, result2) {
      $scope.matchesData.Days[parentIndex].Matches[index].wantToPredict = false;
    };

    $scope.cancelPhase = function(matchIdNumber, index, parentIndex, parentParentIndex, result1, result2) {
      console.log($scope.matchesData);
      console.log('matchIdNumber, index, parentIndex, parentParentIndex, result1, result2)',matchIdNumber, index, parentIndex, parentParentIndex, result1, result2);
      $scope.matchesData.Phases[parentParentIndex].Days[parentIndex].Matches[index].wantToPredict = false;
    };

    $scope.showPhase = function(index) {
      console.log('showPhase on index ', index);
      if($scope.matchesData.Phases[index].showPhase)  {
        $scope.matchesData.Phases[index].showPhase = false;
      } else {
        $scope.matchesData.Phases[index].showPhase = true;
      }
    };

    // $locationChangeSuccess event is triggered when url changes.
    // note: this event is not triggered when page is refreshed.
    $scope.$on('$locationChangeSuccess', function(event) {
      console.log('tournament calendar!!!: location changed:');
      console.log('tournament calendar!!!!!!: routeparams', $routeParams);
      // set tab parameter to new tab parameter.
      // use this to have to last state in order to display the proper view.
      if($routeParams.tab !== undefined) {
        $scope.groupby = $routeParams.groupby;
        $scope.updateMatchesView();
      }
    });
}]);

// TournamentPredictionsCtrl: collects complete data of specific tournament (matches, predictions)
tournamentControllers.controller('TournamentPredictionsCtrl', ['$scope', '$routeParams', 'Tournament', 'User', '$location',
  function($scope, $routeParams, Tournament, User, $location) {
    console.log('Tournament predictions controller');

    Tournament.get({ id:$routeParams.id }).$promise.then(function(response) {
      if($scope.currentUser.Teams !== undefined && $scope.currentUser.Teams.length > 0) {
        $scope.teams = filteredTeams($scope.currentUser.Teams, response.Teams);
        $scope.selectedTeamId = $scope.teams[0].Id;

        $scope.matchesData = Tournament.calendarWithPrediction({id:$routeParams.id, teamId:$scope.selectedTeamId, groupby:$routeParams.groupby});
      }
    });

    $scope.update = function() {
      $scope.matchesData = Tournament.calendarWithPrediction({id:$routeParams.id, teamId:$scope.selectedTeamId, groupby:$routeParams.groupby});
    };

    $scope.groupby = $routeParams.groupby;
    // $locationChangeSuccess event is triggered when url changes.
    // note: this event is not triggered when page is refreshed.
    $scope.$on('$locationChangeSuccess', function(event) {
      console.log('tournament predictions: location changed:');
      console.log('tournament predictions: routeparams', $routeParams);
      // set tab parameter to new tab parameter.
      // use this to have to last state in order to display the proper view.
      if($routeParams.tab !== undefined){
        $scope.groupby = $routeParams.groupby;
      }
    });

    // keep user's teams added to the given tournament
    function filteredTeams(userTeams, tournamentTeams) {
      var teams = [];
      for(var i = 0; i < userTeams.length; i++) {
        for(var j = 0; j < tournamentTeams.length; j++) {
          if(userTeams[i].Id == tournamentTeams[j].Id) {
            teams.push(userTeams[i]);
          }
        }
      }

      return teams;
    }
  }
]);

// TournamentSetResultsCtrl (admin): update results.
// ToDo: Should only be available if you are admin
tournamentControllers.controller('TournamentSetResultsCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
  console.log('Tournament set results controller:');
  console.log('route params', $routeParams)
  $scope.tournamentData = Tournament.get({ id:$routeParams.id });

  $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:"phase"});

  // update result of a match.
  $scope.updateResult = function(match, matchindex, dayindex, phaseindex){
      console.log('TournamentSetResultsCtrl: updateResult');
      console.log('match: ', match);
      console.log('match: ', match.IdNumber);
      console.log('match result: ', match.Result1, ' ', match.Result2);
      console.log('indexes: match, day, phase ', matchindex, dayindex, phaseindex);
      // build result string to send to API
      var result = match.Result1 + ' ' + match.Result2;
      $scope.updatedMatch = Tournament.updateMatchResult({ id:$routeParams.id, matchId:match.IdNumber, result:result});
      // update current match view
      $scope.updatedMatch.$promise.then(function(result){
	  console.log('result: ', result);
	  console.log('matchdata: ', $scope.matchesData);
	  console.log('matchdatamatches: ', $scope.matchesData.Phases[phaseindex].Days[dayindex].Matches[matchindex]);
	  $scope.matchesData.Phases[phaseindex].Days[dayindex].Matches[matchindex] = result;
    });
  };

  // simulate a phase of a tournament.
  $scope.simulatePhase = function(phaseName, phaseindex){
    console.log('TournamentSetResultsCtrl: simulatePhase:', phaseName);
    Tournament.simulatePhase({id:$routeParams.id, phaseName:phaseName},
			     function(result){
			       console.log('success in simulation!');
			       $scope.matchesData.Phases[phaseindex].Days = result.Phase.Days;
			     },
			     function(err) {
			       console.log('failure in  simulation! ', err.data);
			       $scope.messageDanger = err.data;
			     });
  };

    // block preduction
    $scope.blockPrediction = function(match, matchindex, dayindex, phaseindex){
	console.log('block prediction', match, matchindex, dayindex, phaseindex);
	Tournament.blockMatchPrediction({id:$routeParams.id, matchId:match.IdNumber},
				       function(response){
					   $scope.matchesData.Phases[phaseindex].Days[dayindex].Matches[matchindex].CanPredict = response.CanPredict;
				       },
				      function(response){
					  console.log('block match predictions failed.')
				      });
    };

}]);

// TournamentSetTeamsCtrl (admin): change teams.
// ToDo: Should only be available if you are admin
tournamentControllers.controller('TournamentSetTeamsCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament set teams controller:');
    console.log('route params', $routeParams);
    $scope.tournamentData = Tournament.get({ id:$routeParams.id });

    $scope.teamsData = Tournament.teams({id:$routeParams.id, groupby:"phase"});

    $scope.edit = function(index, parentIndex) {
      console.log('edit team: ', index, ' ', parentIndex );
      $scope.teamsData.Phases[parentIndex].Teams[index].wantToEdit = true;
    };

    $scope.save = function(index, parentIndex, oldName, newName, phaseName) {
    	console.log('save team:');
    	console.log('team index:', index);
    	console.log('phase index:', parentIndex);
    	console.log('old name:', oldName);
    	console.log('new name:', newName);

      $scope.teamsData.Phases[parentIndex].Teams[index].wantToEdit = false;
      if(newName === undefined) {
	       return;
      }
      Tournament.updateTeamInPhase({id:$routeParams.id, phaseName:phaseName, oldName:oldName, newName:newName},
        function(result) {
          console.log('success setting team');
          $scope.teamsData.Phases = result.Phases;
        }, function(err) {
          console.log('failure setting team ', err.data);
          $scope.messageDanger = err.data;
        });
    };

    Tournament.calendar({id:$routeParams.id, groupby:"phase"},
      function(result) {
        $scope.isActivated = function(phaseindex, phaseName) {
          // get all the matches of the given phase
          var phase = result.Phases[phaseindex];
          // true if all the matches are not ready otherwise false
          for(var i = 0; i < phase.Days.length; i++) {
            var day = phase.Days[i];
            for(var j = 0; j < day.Matches.length; j++) {
              var match = day.Matches[j];
              if(match.Ready === false) {
                return false;
              }
            }
          }
          return true;
        };
      }, function(err) {
        return false;
      });

    $scope.activate = function(phaseName) {
      Tournament.activatePhase({id:$routeParams.id, phaseName:phaseName},
          function(result) {
            console.log('success activating phase');
          }, function(err) {
            console.log('failure activating phase ', err.data);
            $scope.messageDanger = err.data;
          });
    };
}]);

// TournamentFirstStageCtrl: fetch first stage data of a specific tournament.
tournamentControllers.controller('TournamentFirstStageCtrl',  ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
  console.log('Tournament first stage controller:');
  $scope.tournamentData = Tournament.get({ id:$routeParams.id });

  // list of tournament groups
  $scope.groupsData = Tournament.groups({id:$routeParams.id});
  // predicate is updated for ranking tables
  $scope.predicate = 'Points';

}]);

// TournamentSecondStageCtrl: fetch second stage data of a specific tournament.
tournamentControllers.controller('TournamentSecondStageCtrl',  ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
  console.log('Tournament second stage controller:');
  $scope.tournamentData = Tournament.get({ id:$routeParams.id });
  $scope.matchesData = Tournament.matches({id:$routeParams.id, filter:"second"});
}]);

// TournamentPredictCtrl: fetch predicts of a specific tournament.
tournamentControllers.controller('TournamentPredictCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
  console.log('Tournament predict controller:');
  console.log('route params', $routeParams);
  $scope.tournamentData = Tournament.get({ id:$routeParams.id });

  $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:$routeParams.groupby});

  $scope.activatePredict = function(matchIdNumber, index, parentIndex){
    console.log('TournamentPredictCtrl: activate predict:', matchIdNumber);
    $scope.matchesData.Days[parentIndex].Matches[index].wantToPredict = true;
  };

  $scope.predict = function(matchIdNumber, index, parentIndex, result1, result2){
    console.log('TournamentPredictCtrl: predict:', matchIdNumber);

    $scope.matchesData.Days[parentIndex].Matches[index].wantToPredict = false;
    $scope.matchesData.Days[parentIndex].Matches[index].HasPredict = true;

    Tournament.predict({id:$routeParams.id, matchId:matchIdNumber, result1:result1, result2:result2},
      function(result){
        console.log('success in setting prediction!');
        $scope.matchesData.Days[parentIndex].Matches[index].Predict = result.Predict.Result1 + ' - ' + result.Predict.Result2;
        $scope.messageInfo = result.MessageInfo;
        console.log('match result: ', result.Predict.Result1 + ' - ' + result.Predict.Result2);
      },
      function(err) {
        console.log('failure setting prediction! ', err.data);
        $scope.messageDanger = err.data;
      });
  };
}]);

// TournamentRankingCtrl: used to rank participants of the current tournament.
//
tournamentControllers.controller('TournamentRankingCtrl', ['$scope', '$routeParams', 'Tournament', 'Team', '$location',function($scope, $routeParams, Tournament, Team, $location) {
    console.log('Tournament ranking controller:');
    console.log('route params', $routeParams)

    $scope.tournamentData = Tournament.get({ id:$routeParams.id });

    // get the teams that the user belongs to and that are part of the current tournament.
    $scope.tournamentData.$promise.then(function(response) {
	if($scope.currentUser.Teams !== undefined && $scope.currentUser.Teams.length > 0) {
            $scope.teams = filter($scope.currentUser.Teams, response.Teams);
	}
	
    });
    
    function filter(a, b) {
	var c = [];
	for(var i = 0; i < a.length; i++) {
            for(var j = 0; j < b.length; j++) {
		if(a[i].Id == b[j].Id) {
		    c.push(a[i]);
		}
            }
	}
	return c;
    };

    $scope.rankingData = Tournament.ranking({id:$routeParams.id, rankby:'users'});
    $scope.rankingData.$promise.then(function(response){
    	$scope.selectedParticipants = response.Users;
    });
    
    $scope.update = function() {
	$scope.selectedParticipants = $scope.rankingData.Users;
	if ($scope.selectedTeamId == 0){
	    return;
	}

	Team.members({ id:$scope.selectedTeamId }).$promise.then(function(response){
	    if($scope.selectedParticipants !== undefined && $scope.selectedParticipants.length > 0) {
		$scope.selectedParticipants = filter($scope.selectedParticipants, response.Members);
	    }
	});
    };

    $scope.predicate = ''; // predicate is udate for ranking tables.
}]);

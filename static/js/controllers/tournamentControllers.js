'use strict';

// Tournament controllers manage tournament entities (creation, update, deletion) by getting
// data from REST service (resource).
// Handle also user subscription to a tournament (join/leave and join as team/leave as team).
var tournamentControllers = angular.module('tournamentControllers', []);
// TournamentListCtrl: fetch all tournaments data 
tournamentControllers.controller('TournamentListCtrl', ['$scope', 'Tournament', '$location', function($scope, Tournament, $location) {
  console.log('Tournament list controller:');
  $scope.tournaments = Tournament.query();

  $scope.tournaments.$promise.then(function(result){
    if(!$scope.tournaments || ($scope.tournaments && !$scope.tournaments.length))
      $scope.noTournamentsMessage = 'You have no tournaments';
  });

  $scope.searchTournament = function(){
    console.log('TournamentListCtrl: searchTournament');
    console.log('keywords: ', $scope.keywords)
    $location.search('q', $scope.keywords).path('/tournaments/search');
  };

  // experimental: sar
  // start world cup create action 
  $scope.createWorldCup = function(){
    console.log('Creating world cup');
    Tournament.saveWorldCup($scope.tournament,
		    function(tournament) {
		      console.log('World Cup Tournament: ', tournament);
		      $location.path('/tournaments/' + tournament.Id);
		    },
		    function(err) {
		      console.log('save failed: ', err.data);
		      $scope.messageDanger = err.data;
		    });
  };
  // end world cup create action
}]);

// TournamentCardCtrl: fetch data of a particular tournament.
tournamentControllers.controller('TournamentCardCtrl', ['$scope', 'Tournament',function($scope, Tournament) {
    console.log('Tournament card controller:');
    console.log('tournament ID: ', $scope.$parent.tournament.Id);
    $scope.tournamentData = Tournament.get({ id:$scope.$parent.tournament.Id});
    
    $scope.tournamentData.$promise.then(function(tournamentData){
	$scope.tournament = tournamentData.Tournament;
	console.log('tournament card controller, tournamentData = ', tournamentData);
	$scope.participantsCount = tournamentData.Participants.length;
	$scope.teamsCount = tournamentData.Teams.length;
	$scope.progress = tournamentData.Progress;
    });
}]);

// TournamentSearchCtrl: returns an array of tournaments based on a search query.
tournamentControllers.controller('TournamentSearchCtrl', ['$scope', '$routeParams', 'Tournament', '$location', function($scope, $routeParams, Tournament, $location) {
  console.log('Tournament search controller');
  console.log('routeParams: ', $routeParams);
  // get tournaments data result from search query
  $scope.tournamentsData = Tournament.search( {q:$routeParams.q});
  
  $scope.tournamentsData.$promise.then(function(result){
    $scope.tournaments = result.Tournaments;
    $scope.messageInfo = result.MessageInfo;
  });
  
  $scope.query = $routeParams.q;
  // use the isSearching mode to differientiate:
  // no tournaments in app AND no tournaments found using query search
  $scope.isSearching = true;
  $scope.searchTournament = function(){
    console.log('TournamentSearchCtrl: searchTournament');
    console.log('keywords: ', $scope.keywords)
    $location.search('q', $scope.keywords).path('/tournaments/search');
  };
}]);

// TournamentNewCtrl: use this controller to create a new tournament.
tournamentControllers.controller('TournamentNewCtrl', ['$rootScope', '$scope', 'Tournament', '$location', function($rootScope, $scope, Tournament, $location) {
  console.log('Tournament New controller');
  
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
    
    // get message info from redirects.
    $scope.messageInfo = $rootScope.messageInfo;
    // reset to nil var message info in root scope.
    $rootScope.messageInfo = undefined;
    
    // get candidates data from tournament id
    $scope.candidatesData = Tournament.candidates({id:$routeParams.id});
    
    // do we really need theses lines?
    $scope.candidatesData.$promise.then(function(result){
	$scope.candidates = result.Candidates;
    });
    
    // list of tournament groups
    $scope.groupsData = Tournament.groups({id:$routeParams.id});
    // admin function: reset tournament
    $scope.resetTournament = function(){
	Tournament.reset({id:$routeParams.id},
			 function(result){
			     console.log('reset succeed.');
			     $scope.messageInfo = result.MessageInfo;
			     $scope.groupsData.Groups = result.Groups;
			 },
			 function(err){
			     console.log('reset failed: ', err.data);
			     $scope.messageDanger = err.data;
			 });
    };
    
    
    $scope.deleteTournament = function() {
	if(confirm('Are you sure?')){
	    Tournament.delete({ id:$routeParams.id },
			      function(response){
				  $rootScope.messageInfo = response.MessageInfo;
				  $location.path('/');
			      },
			      function(err) {
				  console.log('delete failed: ', err.data);
				  $scope.messageDanger = err.data;
			      });
	}
    };
    
    $scope.joinTournament = function(){
	Tournament.join({ id:$routeParams.id }).$promise.then(function(response){
	    $scope.joinButtonName = 'Leave';
	    $scope.joinButtonMethod = $scope.leaveTournament;
	    $scope.messageInfo = response.MessageInfo;
	    Tournament.participants({ id:$routeParams.id }).$promise.then(function(participantsResult){
		$scope.tournamentData.Participants = participantsResult.Participants;
	    });
	});
    };
    
    $scope.leaveTournament = function(){
	Tournament.leave({ id:$routeParams.id }).$promise.then(function(response){
	    $scope.joinButtonName = 'Join';
	    $scope.joinButtonMethod = $scope.joinTournament;
	    $scope.messageInfo = response.MessageInfo;
	    Tournament.participants({ id:$routeParams.id }).$promise.then(function(participantsResult){
		$scope.tournamentData.Participants = participantsResult.Participants;
	    });
	});
    };
    
    $scope.joinTournamentAsTeam = function(teamId){
	Tournament.joinAsTeam({id:$routeParams.id, teamId:teamId}).$promise.then(function(response){
	    $scope.joinAsTeamButtonName[teamId] = 'Leave';
	    $scope.joinAsTeamButtonMethod[teamId] = $scope.leaveTournamentAsTeam;
	    $scope.messageInfo = response.MessageInfo;
	    Tournament.get({ id:$routeParams.id }).$promise.then(function(tournamentResult){
		$scope.tournamentData.Teams = tournamentResult.Teams;
	    });
	});
    };
    
    $scope.leaveTournamentAsTeam = function(teamId){
	Tournament.leaveAsTeam({id:$routeParams.id, teamId:teamId}).$promise.then(function(response){
	    $scope.joinAsTeamButtonName[teamId] = 'Join';
	    $scope.joinAsTeamButtonMethod[teamId] = $scope.joinTournamentAsTeam;
	    $scope.messageInfo = response.MessageInfo;
	    Tournament.get({ id:$routeParams.id }).$promise.then(function(tournamentResult){
		$scope.tournamentData.Teams = tournamentResult.Teams;
	    });
	});
    };

    // tab at undefined means you are in tournament/:id url
    if($routeParams.tab == undefined){ 
	$scope.tournamentData.$promise.then(function(result){
	    if(result.Tournament.IsFirstStageComplete){
		$scope.tab = 'secondstageh';
	    }else{
		$scope.tab = 'firststage';
	    }
	});
    } else {
	// Initialize tab variable to handle views:
	$scope.tab = $routeParams.tab;
    }

    // Is tournament admin flag identifies if current user is also the admin of the tournament.
    // Use this flag to show specific information of tournament admin.
    $scope.isTournamentAdmin = false;
    $scope.tournamentData.$promise.then(function(result){
	if(result.Tournament.AdminId == $rootScope.currentUser.User.Id){
	    console.log('tournament is admin TRUE');
	    $scope.isTournamentAdmin = true;
	}else{
	    console.log('tournament is admin FALSE');
	    $scope.isTournamentAdmin = false;
	}
    });
    
    // Checks if user has joined a tournament
    $scope.joined = $scope.tournamentData.$promise.then(function(result){
	console.log('tournament joined ready!');
	return result.Joined;
    });
    
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
    
    $scope.candidatesData.$promise.then(function(candidatesResult){
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
    
    // Action triggered when 'Create new button' is clicked, modal window will be hidden.
    // We also set flag 'redirectToNewTeam' to true for listener to know if redirection is needed.
    $scope.newTeam = function(){
	$('#tournament-modal').modal('hide');
	$scope.redirectToNewTeam = true;
    };
    
    // listen 'hidden.bs.modal' event to redirect to new team page
    // Only redirect if flag 'redirectToNewTeam' is set.
    $('#tournament-modal').on('hidden.bs.modal', function (e) {
	// need to have scope for $location to work. So add 'apply' function
	// inside js listener
	if($scope.redirectToNewTeam == true){
	    $scope.$apply(function(){
		$location.path('/teams/new/');
	    });
	}
    })

    // $locationChangeSuccess event is triggered when url changes.
    // note: this event is not triggered when page is refreshed.
    $scope.$on('$locationChangeSuccess', function(event) {
	console.log('tournament show: location changed:');
	console.log('tournament show: routeparams', $routeParams);
	// set tab paramter to new tab parameter.
	// use this to have to last state in order to display the proper view.
	if($routeParams.tab != undefined){
	    $scope.tab = $routeParams.tab;
	}
    });

}]);

// TournamentEditCtrl: collects data to update an existing tournament.
tournamentControllers.controller('TournamentEditCtrl', ['$rootScope', '$scope', '$routeParams', 'Tournament', '$location',function($rootScope, $scope, $routeParams, Tournament, $location) {
  $scope.tournamentData = Tournament.get({ id:$routeParams.id });
  
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
tournamentControllers.controller('TournamentCalendarCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament calendar controller');
    console.log('route params', $routeParams)
    $scope.tournamentData = Tournament.get({ id:$routeParams.id });
    
    $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:$routeParams.groupby});
    
    $scope.byPhaseOnClick = function(){
	$scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'phase'});
    };

    $scope.byDateOnClick = function(){
	$scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'date'});	
    };
    
  $scope.activatePredict = function(matchIdNumber, index, parentIndex){
    console.log('Tournament calendar controller: activate predict:', matchIdNumber);
    $scope.matchesData.Days[parentIndex].Matches[index].wantToPredict = true;
  };

  $scope.predict = function(matchIdNumber, index, parentIndex, result1, result2){
    console.log('Tournament calendar controller: predict:', matchIdNumber);
    console.log('Tournament calendar controller: predict: index:', index);
    console.log('Tournament calendar controller: predict: parent parent index:', parentIndex);

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
    console.log('match result: ', result1, ' ', result2);

  };  

}]);

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
    
}]);

// TournamentSetTeamsCtrl (admin): change teams.
// ToDo: Should only be available if you are admin
tournamentControllers.controller('TournamentSetTeamsCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament set teams controller:');
    console.log('route params', $routeParams)
    $scope.tournamentData = Tournament.get({ id:$routeParams.id });
    
    $scope.teamsData = Tournament.teams({id:$routeParams.id, groupby:"phase"});
    
    $scope.edit = function(index, parentIndex){
	console.log('edit team: ', index, ' ', parentIndex );
	$scope.teamsData.Phases[parentIndex].Teams[index].wantToEdit = true;
    }

    $scope.save = function(index, parentIndex, oldName, newName, phaseName){
	console.log('save team:');
	console.log('team index:', index);
	console.log('phase index:', parentIndex);
	console.log('old name:', oldName);
	console.log('new name:', newName);

	$scope.teamsData.Phases[parentIndex].Teams[index].wantToEdit = false;
	if(newName == undefined){
	    return;
	}
	Tournament.updateTeamInPhase({id:$routeParams.id, phaseName:phaseName, oldName:oldName, newName:newName},
				     function(result){
					 console.log('success setting team');
					 $scope.teamsData.Phases = result.Phases;
				     },
				     function(err) {
					 console.log('failure setting team ', err.data);
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
  // predicate is udate for ranking tables
  $scope.predicate = '';

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
  console.log('route params', $routeParams)
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

// TournamentRankingCtrl: fetch ranking data of a specific tournament.
tournamentControllers.controller('TournamentRankingCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament ranking controller:');
    console.log('route params', $routeParams)
    $scope.tournamentData = Tournament.get({ id:$routeParams.id });
    $scope.rankBy = 'users'
    $scope.rankingData = Tournament.ranking({id:$routeParams.id, rankby:$routeParams.rankby});

    // predicate is udate for ranking tables
    $scope.predicate = '';
    
    $scope.byUsersRankOnClick = function(){
	if($scope.rankBy == 'user'){
	    return;
	}
	$scope.rankBy = 'users';
	$scope.rankingData = Tournament.ranking({id:$routeParams.id, rankby:$scope.rankBy});
	return;
    };

    $scope.byTeamsRankOnClick = function(){
	if($scope.rankBy == 'teams'){
	    return;
	}
	$scope.rankBy = 'teams';
	$scope.rankingData = Tournament.ranking({id:$routeParams.id, rankby:$scope.rankBy});
	return;
    };

}]);

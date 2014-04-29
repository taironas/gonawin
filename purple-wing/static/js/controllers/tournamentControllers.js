'use strict';

// Tournament controllers manage tournament entities (creation, update, deletion) by getting
// data from REST service (resource).
// Handle also user subscription to a tournament (join/leave and join as team/leave as team).
var tournamentControllers = angular.module('tournamentControllers', []);
// TournamentListCtrl: fetch all tournaments data
tournamentControllers.controller('TournamentListCtrl', ['$scope', 'Tournament', '$location', function($scope, Tournament, $location) {
  console.log('Tournament list controller:');

    $scope.count = 25;  // counter for the number of tournaments to display in view.
    $scope.page = 1;    // page counter for tournaments, to know which page to display next.

    // main query to /j/tournaments to get all tournaments.
    $scope.tournaments = Tournament.query({count:$scope.countTeams, page:$scope.pageTeams});

    $scope.tournaments.$promise.then(function(response){
	if(!$scope.tournaments || ($scope.tournaments && !$scope.tournaments.length)){
	    $scope.noTournamentsMessage = 'You have no tournaments';
	}else if($scope.tournaments != undefined){
	    $scope.showMoreTournaments = (response.length == $scope.count);
	}
    });
    
    // show more tournaments function:
    // retreive tournaments by page and increment page.
    $scope.moreTournaments = function(){
	console.log('more tournaments');
	$scope.page = $scope.page + 1;
	Tournament.query({count:$scope.count, page:$scope.page}).$promise.then(function(response){
	    console.log('response: ', response);
	    $scope.tournaments = $scope.tournaments.concat(response);
	    $scope.showMoreTournaments = (response.length == $scope.count);
	});
    };

    $scope.searchTournament = function(){
	console.log('TournamentListCtrl: searchTournament');
	console.log('keywords: ', $scope.keywords)
	$location.search('q', $scope.keywords).path('/search');
    };
    
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
    $scope.candidateTeamsData = Tournament.candidates({id:$routeParams.id});

    // do we really need theses lines?
    $scope.candidateTeamsData.$promise.then(function(result){
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

  // tab at undefined means you are in tournament/:id url.
  // So calendar should be active by default
  if($routeParams.tab == undefined){
    $scope.tab = 'calendar';
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

    // set up join as a team buttons.
    $scope.candidateTeamsData.$promise.then(function(candidatesResult){
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
    Tournament.participants({ id:$routeParams.id }).$promise.then(function(participantsResult){
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
    
  $scope.tabs = {
    "calendar":         { title: 'Calendar',                url: 'templates/tournaments/tab_calendar.html' },
    "firststage":       { title: 'First Stage',             url: 'templates/tournaments/tab_firststage.html' }, 
    "secondstageh":     { title: 'Second Stage Horizontal', url: 'templates/tournaments/tab_bracketH.html' },
    "secondstagev":     { title: 'Second Stage Vertical',   url: 'templates/tournaments/tab_bracketV.html' },
    "ranking":          { title: 'Ranking',                 url: 'templates/tournaments/tab_ranking.html' },
    "admin.setresults": { title: 'Set Results',             url: 'templates/tournaments/tab_setresults.html' },
    "admin.setteams":   { title: 'Set Teams',               url: 'templates/tournaments/tab_setteams.html' }
  };
  
  // set the current tab based on the 'tab' parameter
  if($scope.tab == undefined) {
    $scope.currentTab = $scope.tabs["calendar"].url;
  } else {
    $scope.currentTab = $scope.tabs[$scope.tab].url;
  }

  $scope.onClickTab = function (tab) {
    $scope.currentTab = tab.url;
  }

  // $locationChangeSuccess event is triggered when url changes.
  // note: this event is not triggered when page is refreshed.
  $scope.$on('$locationChangeSuccess', function(event) {
    console.log('tournament show: location changed:');
    console.log('tournament show: routeparams', $routeParams);
    // set tab parameter to new tab parameter.
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
    console.log('routeparams!!!!!', $routeParams);

    $scope.groupby = $routeParams.groupby;

  $scope.updateMatchesView = function(){
    if($scope.groupby != undefined){
      if($scope.groupby == 'phase'){
        $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'phase'});
      } else if($scope.groupby == 'date'){
        $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'date'});
      }
    } else{
      $scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'date'});
    }
  };
  $scope.updateMatchesView();
    // $scope.byPhaseOnClick = function(){
    // 	$scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'phase'});
    // };

    // $scope.byDateOnClick = function(){
    // 	$scope.matchesData = Tournament.calendar({id:$routeParams.id, groupby:'date'});
    // };

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

  // $locationChangeSuccess event is triggered when url changes.
  // note: this event is not triggered when page is refreshed.
  $scope.$on('$locationChangeSuccess', function(event) {
    console.log('tournament calendar!!!: location changed:');
    console.log('tournament calendar!!!!!!: routeparams', $routeParams);
    // set tab parameter to new tab parameter.
    // use this to have to last state in order to display the proper view.
    if($routeParams.tab != undefined){
      $scope.groupby = $routeParams.groupby;
      $scope.updateMatchesView();
    }
  });
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

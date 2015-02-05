'use strict'
var dataServices = angular.module('dataServices', ['ngResource']);

dataServices.factory('User', function($http, $resource, $cookieStore) {
  $http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');

  var User = $resource('j/users/:id', {id:'@id', including:'@including', teamId:'@teamId'}, {
    get: { method: 'GET', params: {including: '@including'}, url: 'j/users/show/:id' },
    update: { method: 'POST', url: 'j/users/update/:id' },
    delete: { method: 'POST', url: 'j/users/destroy/:id' },
    scores: {method: 'GET', url: 'j/users/:id/scores'},
    search: { method: 'GET', url: 'j/users/search?q=:q', cache : true},
    teams: {method : 'GET', url: 'j/users/:id/teams'},
    tournaments: {method : 'GET', url: 'j/users/:id/tournaments'},
    allowInvitation : {method: 'POST', url: 'j/users/allow/:teamId'},
    denyInvitation : {method: 'POST', url: 'j/users/deny/:teamId'},
  })
  // define display name to handle alias or user name.
  // Note: There is another displayName function definition in the Session ressource as we handle users via User and Session.
  User.prototype.displayName = function() {
    if(this.User == undefined) return;
    if(this.User.Alias.length > 0){
      return this.User.Alias;
    } else{
      return this.User.Username;
    }
  };
  return User;
});

dataServices.factory('Tournament', function($http, $resource, $cookieStore) {
  $http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');

  return $resource('j/tournaments/:id',
		   {
		     id:'@id',
		     q:'@q',
		     teamId:'@teamId',
		     groupby: '@groupby',
		     filter: '@filter',
		     matchId: '@matchId',
		     result: '@result',
		     result1: '@result1',
		     result2: '@result2',
		     phaseName: '@phaseName',
		     rankby: '@rankby',
		     limit: '@limit',
		     oldName: '@oldName',
		     newName: '@newName',
		     userId: '@userId'
		   },
		   {
		     get: { method: 'GET', url: 'j/tournaments/show/:id' },
		     save: { method: 'POST', url: 'j/tournaments/new' },
		     update: { method: 'POST', url: 'j/tournaments/update/:id' },
		     delete: { method: 'POST', url: 'j/tournaments/destroy/:id' },
		     search: { method: 'GET', url: 'j/tournaments/search?q=:q', cache : true},
		     participants: { method: 'GET', url:'j/tournaments/:id/participants' },
		     join: {method: 'POST', url: 'j/tournaments/join/:id' },
		     leave: {method: 'POST', url: 'j/tournaments/leave/:id' },
		     joinAsTeam: {method: 'POST', url: 'j/tournaments/joinasteam/:id/:teamId' },
		     leaveAsTeam: {method: 'POST', url: 'j/tournaments/leaveasteam/:id/:teamId' },
		     candidates: {method: 'GET', url: 'j/tournaments/:id/candidates' },
		     saveWorldCup: {method: 'POST', url: 'j/tournaments/newwc'},
		     getWorldCup: {method: 'GET', url: 'j/tournaments/getwc', cache : true},
		     groups: {method: 'GET', url: 'j/tournaments/:id/groups', cache : true},
		     calendar: {method: 'GET', url: 'j/tournaments/:id/calendar?groupby=:groupby'},
		     calendarWithPrediction: {method: 'GET', url: 'j/tournaments/:id/:teamId/calendarwithprediction?groupby=:groupby'},
		     matches: {method: 'GET', url: 'j/tournaments/:id/matches?filter=:filter'},
		     updateMatchResult: {method: 'POST', url: '/j/tournaments/:id/matches/:matchId/update?result=:result'},
		     blockMatchPrediction: {method: 'POST', url: '/j/tournaments/:id/matches/:matchId/blockprediction'},
		     simulatePhase: {method: 'POST', url: '/j/tournaments/:id/matches/simulate?phase=:phaseName'},
		     reset: {method: 'POST', url: '/j/tournaments/:id/admin/reset'},
		     predict: {method: 'POST', url: '/j/tournaments/:id/matches/:matchId/predict?result1=:result1&result2=:result2'},
		     ranking: {method: 'GET', url: 'j/tournaments/:id/ranking?rankby=:rankby&limit=:limit' },
		     teams: {method: 'GET', url: 'j/tournaments/:id/teams?rankby=:rankby' },
		     updateTeamInPhase: {method: 'POST', url: 'j/tournaments/:id/admin/updateteam?phase=:phaseName&old=:oldName&new=:newName'},
		     addAdmin: {method: 'POST', url: 'j/tournaments/:id/admin/add/:userId'},
		     removeAdmin: {method: 'POST', url: 'j/tournaments/:id/admin/remove/:userId'},
		     syncScores: {method: 'POST', url: 'j/tournaments/:id/admin/syncscores'}
		   })
});

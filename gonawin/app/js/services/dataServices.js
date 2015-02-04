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

dataServices.factory('Team', function($http, $resource, $cookieStore) {
    $http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');

    return $resource('j/teams/:id', {
      id:'@id',
      q:'@q',
      requestId: '@requestId',
      rankby: '@rankby',
      tournamentId: '@tournamentId',
      limit: '@limit',
      userId: '@userId',
      count: '@count',
      page: '@page'
    },
    {
      get: { method: 'GET', url: 'j/teams/show/:id' },
      save: { method: 'POST', url: 'j/teams/new' },
      update: { method: 'POST', url: 'j/teams/update/:id' },
      delete: { method: 'POST', url: 'j/teams/destroy/:id' },
      search: { method: 'GET', url: 'j/teams/search?q=:q', cache : true},
      members: { method: 'GET', url:'j/teams/:id/members'},
      join: {method: 'POST', url: 'j/teams/join/:id'},
      leave: {method: 'POST', url: 'j/teams/leave/:id'},
      requestInvite: {method: 'POST', url: 'j/teams/requestinvite/:id'},
      sendInvite: {method: 'POST', url: 'j/teams/sendinvite/:id/:userId'},
      invited: {method: 'GET', url: 'j/teams/invited/:id'},
      allowRequest : {method: 'POST', url: 'j/teams/allow/:requestId'},
      denyRequest : {method: 'POST', url: 'j/teams/deny/:requestId'},
      ranking: {method: 'GET', url: 'j/teams/:id/ranking?rankby=:rankby&limit=:limit'},
      accuracies: {method: 'GET', url: 'j/teams/:id/accuracies'},
      accuracy: {method: 'GET', url: 'j/teams/:id/accuracies/:tournamentId'},
      prices: {method: 'GET', url: 'j/teams/:id/prices', cache : true},
      price: {method: 'GET', url: 'j/teams/:id/prices/:tournamentId', cache : true},
      updatePrice: {method: 'POST', url: 'j/teams/:id/prices/update/:tournamentId'},
      addAdmin: {method: 'POST', url: 'j/teams/:id/admin/add/:userId'},
      removeAdmin: {method: 'POST', url: 'j/teams/:id/admin/remove/:userId'}
    })
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

dataServices.factory('Invite', function($http, $cookieStore, $resource){
  $http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');

  return $resource('j/invite', {emails: '@emails'}, {
    send: {method: 'POST', params: {emails: '@emails'}, url: 'j/invite'}
  })
});

dataServices.factory('Activity', function($http, $cookieStore, $resource){
  $http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');

  return $resource('j/activities', {count: '@count', page: '@page'})
});

dataServices.factory('Session', function($cookieStore, $resource) {

  var Session = $resource('/j/auth/', {access_token:'@access_token', id:'@id', name:'@name', email:'@email'}, {
    fetchUserInfo: { method:'GET', params: {access_token:'@access_token'}, url: 'https://www.googleapis.com/plus/v1/people/me' },
    fetchUser: { method:'GET', params: {access_token:'@access_token', provider:'@provider', id:'@id', name:'@name', email:'@email'}, url: '/j/auth' },
    logout: { method:'JSONP', params: {token:'@token', callback: 'JSON_CALLBACK'}, url: 'https://accounts.google.com/o/oauth2/revoke' },
    authenticateWithTwitter: { method:'GET', url: '/j/auth/twitter' },
    fetchTwitterUser: { method:'GET', params: { oauth_token: '@oauth_token', oauth_verifier: '@oauth_verifier' }, url: '/j/auth/twitter/user/' },
    fetchGoogleLoginUrl: { method:'GET', url: '/j/auth/googleloginurl' },
    authenticateWithGoogle: { method:'GET', url: '/j/auth/google' },
    fetchGoogleUser: { method:'GET', params: { auth_token: '@auth_token' }, url: '/j/auth/google/user/' },
    DeleteGoogleCookie: { method: 'GET', url: '/j/auth/google/deletecookie'},
    serviceIds: { method:'GET', url: '/j/auth/serviceids/' },
  });

  // Need to define displayname function here again as User can be either returned by the server or the session.
  Session.prototype.displayName = function() {
    if(this.User == undefined) return;
    if(this.User.Alias.length > 0){
      return this.User.Alias;
    } else{
      return this.User.Username;
    }
  };
  return Session;
});
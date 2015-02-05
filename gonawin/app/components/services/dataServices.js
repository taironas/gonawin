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

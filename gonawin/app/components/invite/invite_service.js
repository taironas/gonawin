'use strict'
var inviteService = angular.module('inviteService', ['ngResource']);

inviteService.factory('Invite', function($http, $cookieStore, $resource){
  $http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');

  return $resource('j/invite', {emails: '@emails'}, {
    send: {method: 'POST', params: {emails: '@emails'}, url: 'j/invite'}
  })
});

'use strict'
var teamService = angular.module('teamService', ['ngResource']);

teamService.factory('Team', function($http, $resource, $cookieStore) {
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

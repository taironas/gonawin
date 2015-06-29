'use strict';

var tournamentService = angular.module('tournamentService', ['ngResource']);

tournamentService.factory('Tournament', function($http, $resource, $cookieStore) {
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
       activatePhase: {method: 'POST', url: 'j/tournaments/:id/admin/activatephase?phase=:phaseName'},
       addAdmin: {method: 'POST', url: 'j/tournaments/:id/admin/add/:userId'},
       removeAdmin: {method: 'POST', url: 'j/tournaments/:id/admin/remove/:userId'},
       saveChampionsLeague: {method: 'POST', url: 'j/tournaments/newcl'},
       getChampionsLeague: {method: 'GET', url: 'j/tournaments/getcl', cache : true},
       saveCopaAmerica: {method: 'POST', url: 'j/tournaments/newca'},
       getCopaAmerica: {method: 'GET', url: 'j/tournaments/getca', cache : true},
     });
});

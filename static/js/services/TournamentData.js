'use strict'
purpleWingApp.factory('tournamentData', function($http, $log, $q){
    return {
	getTournament:function(tournamentId){
	    var deferred = $q.defer();
            $http({method: 'GET', url:'/j/tournaments/'+tournamentId}).
                success(function(data,status,headers,config){
                    deferred.resolve(data);
                    $log.info(data, status, headers() ,config);
                }).
                error(function (data, status, headers, config){
                    $log.warn(data, status, headers(), config);
                    deferred.reject(status);
                });
            return deferred.promise;
	}
    };
});

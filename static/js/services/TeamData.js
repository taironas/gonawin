'use strict'
purpleWingApp.factory('teamData', function($http, $log, $q){
    return {
	getTeam:function(teamId){
	    var deferred = $q.defer();
            $http({method: 'GET', url:'/j/teams/'+teamId}).
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

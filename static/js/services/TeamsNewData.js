'use strict'
purpleWingApp.factory('teamsNewData', function($http, $log, $q){
    return {
	getData:function(){
	    var deferred = $q.defer();
            $http({method: 'GET', url:'/j/teams/new/'}).
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

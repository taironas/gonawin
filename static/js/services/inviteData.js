'use strict'
purpleWingApp.factory('inviteData', function($http, $log, $q){
    return {
	getData:function(){
	    var deferred = $q.defer();
            $http({method: 'GET', url:'/j/invite/'}).
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

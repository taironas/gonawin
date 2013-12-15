'use strict'
purpleWingApp.factory('userData', function($http, $log, $q){
    return {
	getUser:function(){
	    var deferred = $q.defer();
            $http({method: 'GET', url:'/j/users/1'}).
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

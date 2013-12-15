'use strict';

var purpleWingApp = angular.module('purpleWingApp', ['ngSanitize', 'directive.g+signin', 'ngCookies'])
    .config(function($routeProvider){
	$routeProvider.when('/',
			    {
				templateUrl: 'templates/main.html', 
				controller: 'MainController'
			    });
	$routeProvider.when('/about',
			    {
				templateUrl: 'templates/about.html',
				controller: 'AboutController'
			    });
	$routeProvider.otherwise( {redirectTo: '/'});
    })
    .factory('myCache', function($cacheFactory){
	return $cacheFactory('myCache', {capacity:3})
    });

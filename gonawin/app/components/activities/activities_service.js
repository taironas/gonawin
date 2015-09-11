'use strict';
var activitiesService = angular.module('activitiesService', ['ngResource']);

activitiesService.factory('Activity', function($http, $cookieStore, $resource){
  $http.defaults.headers.common.Authorization = $cookieStore.get('auth');

  return $resource('j/activities', {count: '@count', page: '@page'});
});

'use strict';

var navigationControllers = angular.module('navigationControllers', []);

navigationControllers.controller('NavigationCtrl', ['$scope', '$location', '$cookieStore', 'Session', 'sAuth', 
  function ($scope, $location, $cookieStore, Session, sAuth) {
    console.log('NavigationCtrl module');
    
    $scope.disconnect = function(){
      console.log('NavigationCtrl module:: disconnect');

      Session.logout({ token: $cookieStore.get('access_token') });
      sAuth.logout();
      
      $cookieStore.remove('auth');
      $cookieStore.remove('access_token');
      $cookieStore.remove('user_id');
      $cookieStore.remove('logged_in');
      
      $location.path('/welcome');
    };
}]);

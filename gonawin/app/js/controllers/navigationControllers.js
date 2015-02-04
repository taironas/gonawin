'use strict';

var navigationControllers = angular.module('navigationControllers', []);

navigationControllers.controller('NavigationCtrl', ['$scope', '$rootScope', '$location', '$cookieStore', 'Session', 'sAuth',
  function ($scope, $rootScope, $location, $cookieStore, Session, sAuth) {
    console.log('NavigationCtrl module');
    console.log('NavigationCtrl, current user = ', $rootScope.currentUser);
    console.log('NavigationCtrl, isLoggedIn = ', $rootScope.isLoggedIn);

    $scope.disconnect = function(){
      console.log('NavigationCtrl module:: disconnect');
      
      var providerCookie = $cookieStore.get('provider')
      // logout from Google+
      if(providerCookie == 'google_plus') {
        Session.logout({ token: $cookieStore.get('access_token') });
      }
      // logout from Facebook
      if(providerCookie == 'facebook') {
        sAuth.FBlogout();
      }
      // delete cookie created by Google account
      if(providerCookie == 'google') {
        Session.DeleteGoogleCookie();
      }
      // reset rootScope variables
      $rootScope.currentUser = undefined;
      $rootScope.isLoggedIn = false;

      sAuth.clearCookies();

      $location.path('/welcome');
    };
}]);

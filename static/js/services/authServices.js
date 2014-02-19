'use strict'
var authServices = angular.module('authServices', ['ngResource']);

authServices.factory('sAuth', function($rootScope, $cookieStore, $location, Session) {
  return {
    watchLoginChange: function() {
      var _self = this;

      FB.Event.subscribe('auth.authResponseChange', function(response) {
        console.log('auth.authResponseChange, response = ', response);
        if (response.status === 'connected') {
          /* 
           The user is already logged, 
           is possible retrieve his personal info
          */
          _self.getUserInfo(response.authResponse.accessToken);
          /*
           This is also the point where you should create a 
           session for the current user.
           For this purpose you can use the data inside the 
           response.authResponse object.
          */
        } 
        else {
          /*
           The user is not logged to the app, or into Facebook:
           destroy the session on the server.
          */
        }
      });
    },
    
    getUserInfo: function(accessToken) {
      var _self = this;
      FB.api('/me', function(userInfo) {
        console.log('/me = ', userInfo);
        Session.fetchUser({  access_token: accessToken,
                             provider: 'facebook',
                             id:userInfo.id, 
                             name:userInfo.name, 
                             email:userInfo.email } ).$promise.then(function(userData) {
         $rootScope.currentUser = userData.User;
         console.log('current user: ', $rootScope.currentUser);
         
         $cookieStore.put('access_token', accessToken);
         $cookieStore.put('auth', $rootScope.currentUser.Auth);
         $cookieStore.put('user_id', $rootScope.currentUser.Id);
         
         $rootScope.loggedIn = true;
         
         $location.path('/home');
        });
      });
    },
    
    logout: function() {
      var _self = this;
      FB.logout(function(response) {
        console.log('Facebook logout = ', response);	
      });
    }
  }
});

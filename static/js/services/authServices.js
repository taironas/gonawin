'use strict'
var authServices = angular.module('authServices', ['ngResource']);

authServices.factory('sAuth', function($rootScope, $cookieStore, $location, $q, User, Session) {
  return {
    getCurrentUser: function() {
      var deferred = $q.defer();
      
      if(!$cookieStore.get('access_token') || !$cookieStore.get('user_id')) {
        deferred.resolve(undefined);
      }
      else {
        User.get({ id:$cookieStore.get('user_id') }).$promise.then(function(userData) {
          if(userData.User.Auth == $cookieStore.get('auth')){
            deferred.resolve(userData.User);
          } else {
            deferred.resolve(undefined);
          }
        });
      }

      return deferred.promise;
    },
  
    watchLoginChange: function() {
      var _self = this;

      FB.Event.subscribe('auth.authResponseChange', function(response) {
        console.log('auth.authResponseChange, response = ', response);
        if (response.status === 'connected') {
          _self.getUserInfo(response.authResponse.accessToken);
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
          $rootScope.currentUser = _self.currentUser = userData.User;
          console.log('current user: ', _self.currentUser);
         
          _self.storeCookies(accessToken, _self.currentUser.Auth, _self.currentUser.Id);
         
          //$location.path('/');
        });
      });
    },
    
    storeCookies: function(accessToken, auth, userId) {
      $cookieStore.put('access_token', accessToken);
      $cookieStore.put('auth', auth);
      $cookieStore.put('user_id', userId);
      $cookieStore.put('logged_in', true);
    },
    
    logout: function() {
      var _self = this;
      FB.getLoginStatus(function(response) {
        if (response.status === 'connected') {
          FB.logout(function(response) {
            console.log('Facebook logout = ', response);
            $rootScope.$apply(function() {
              $rootScope.currentUser = _self.currentUser = {};
            });
          });
        }
      });
    },
    
    signinWithTwitter: function(oauthToken, oauthVerifier) {
      // User successfully authorized via Twitter!
      console.log('User successfully authorized via Twitter!');

      Session.fetchTwitterUser({ oauth_token: oauthToken, oauth_verifier: oauthVerifier }).$promise.then(function(userData) {
        $rootScope.currentUser = userData.User;
         console.log('current user: ', $rootScope.currentUser);
         
         _self.storeCookies(oauthToken, _self.currentUser.Auth, _self.currentUser.Id);
         
         $location.path('/');
      });
    }
  }
});

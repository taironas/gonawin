'use strict'
var authServices = angular.module('authServices', ['ngResource']);

authServices.factory('sAuth', function($rootScope, $cookieStore, $location, $q, $timeout, User, Session) {
  return {
    /* returns true when user is logged in based on cookies */
    isLoggedIn: function() {
      var _isLoggedIn = false;
      if($cookieStore.get('access_token') && $cookieStore.get('auth') && $cookieStore.get('user_id') && $cookieStore.get('logged_in')) {
        _isLoggedIn = true;
      }
      return _isLoggedIn;
    },
    getUserID: function() {
      return $cookieStore.get('user_id');
    },
    /* event Facebook handle. Called when login has been detected */
    watchLoginChange: function() {
      var _self = this;

      FB.Event.subscribe('auth.authResponseChange', function(response) {
        console.log('auth.authResponseChange, response = ', response);
        if (response.status === 'connected' && $rootScope.isLoggedIn == false) {
          _self.getFBUserInfo(response.authResponse.accessToken);
        }
      });
    },
    /* get info of Facebook user */
    getFBUserInfo: function(accessToken) {
      var _self = this;
      FB.api('/me', function(userInfo) {
        console.log('/me = ', userInfo);
        $rootScope.currentUser = Session.fetchUser({ 
          access_token: accessToken,
          provider: 'facebook',
          id:userInfo.id,
          name:userInfo.name,
          email:userInfo.email } );
        $rootScope.currentUser.$promise.then(function(currentUser){
          console.log('authServices.getFBUserInfo: current user = ', currentUser);
          _self.storeCookies(accessToken, currentUser.User.Auth, currentUser.User.Id);
          $rootScope.isLoggedIn = true;
          $location.path('/');
        });
      });
    },
    /* store cookies which will be used to dertermine if a user is logged 
     * and to add authentication data in API requests */
    storeCookies: function(accessToken, auth, userId) {
      $cookieStore.put('access_token', accessToken);
      $cookieStore.put('auth', auth);
      $cookieStore.put('user_id', userId);
      $cookieStore.put('logged_in', true);
    },
    /* logout the user who was logged in via Facebook */
    FBlogout: function() {
      var _self = this;
      FB.getLoginStatus(function(response) {
        if (response.status === 'connected') {
          FB.logout(function(response) {
            console.log('Facebook logout = ', response);
          });
        }
      });
    },
    /* Complete signin with Twitter.
     * Fetch Twitter user info then set the current user
     * and store the cookies */
    signinWithTwitter: function(oauthToken, oauthVerifier) {
      var _self = this;
      // User successfully authorized via Twitter!
      console.log('User successfully authorized via Twitter!');

      $rootScope.currentUser = Session.fetchTwitterUser({ oauth_token: oauthToken, oauth_verifier: oauthVerifier });
      $rootScope.currentUser.$promise.then(function(currentUser){
        console.log('signinWithTwitter: current user = ', currentUser);
        _self.storeCookies(oauthToken, currentUser.User.Auth, currentUser.User.Id);
        $rootScope.isLoggedIn = true;
        $location.path('/');
      });
    },
    /* Complete signin with Google.
     * Fetch Google user info then set the current user
     * and store the cookies */
    signinWithGoogle: function(oauthToken, oauthVerifier) {
      var _self = this;
      $rootScope.currentUser = Session.fetchGoogleUser({ oauth_token: oauthToken });
      $rootScope.currentUser.$promise.then(function(currentUser){
        console.log('signinWithGoogle: current user = ', currentUser);
        _self.storeCookies(oauthToken, currentUser.User.Auth, currentUser.User.Id);
        $rootScope.isLoggedIn = true;
        $location.path('/');
      }, function(error){
        $rootScope.currentUser = undefined;
      });
    }
  }
});

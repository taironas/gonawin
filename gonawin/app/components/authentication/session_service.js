'use strict'
var sessionServices = angular.module('sessionServices', ['ngResource']);

sessionServices.factory('Session', function($cookieStore, $resource) {

  var Session = $resource('/j/auth/', {access_token:'@access_token', id:'@id', name:'@name', email:'@email'}, {
    fetchUserInfo: { method:'GET', params: {access_token:'@access_token'}, url: 'https://www.googleapis.com/plus/v1/people/me' },
    fetchUser: { method:'GET', params: {access_token:'@access_token', provider:'@provider', id:'@id', name:'@name', email:'@email'}, url: '/j/auth' },
    logout: { method:'JSONP', params: {token:'@token', callback: 'JSON_CALLBACK'}, url: 'https://accounts.google.com/o/oauth2/revoke' },
    authenticateWithTwitter: { method:'GET', url: '/j/auth/twitter' },
    fetchTwitterUser: { method:'GET', params: { oauth_token: '@oauth_token', oauth_verifier: '@oauth_verifier' }, url: '/j/auth/twitter/user/' },
    fetchGoogleLoginUrl: { method:'GET', url: '/j/auth/googleloginurl' },
    authenticateWithGoogle: { method:'GET', url: '/j/auth/google' },
    fetchGoogleUser: { method:'GET', params: { auth_token: '@auth_token' }, url: '/j/auth/google/user/' },
    DeleteGoogleCookie: { method: 'GET', url: '/j/auth/google/deletecookie'},
    serviceIds: { method:'GET', url: '/j/auth/serviceids/' },
  });

  // Need to define displayname function here again as User can be either returned by the server or the session.
  Session.prototype.displayName = function() {
    if(this.User == undefined) return;
    if(this.User.Alias.length > 0){
      return this.User.Alias;
    } else{
      return this.User.Username;
    }
  };
  return Session;
});

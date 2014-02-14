'use strict';

var sessionControllers = angular.module('sessionControllers', []);

sessionControllers.controller('SessionCtrl', ['$scope', '$location', '$cookieStore', '$q', 'Session', 'User', function ($scope, $location, $cookieStore, $q, Session, User) {
  console.log('SessionController module');
  $scope.currentUser  = undefined;
  $scope.loggedIn = false;
  
  $scope.initSession = function() {
    console.log('SessionController module:: initSession');
    var deferred = $q.defer();
    if(!$cookieStore.get('access_token') || !$cookieStore.get('user_id')) {
      $scope.loggedIn = false;
      deferred.resolve(false);
    } 
    else if($scope.currentUser && $scope.loggedIn){
      deferred.resolve(true);
    }
    else {
      User.get({ id:$cookieStore.get('user_id') }).$promise.then(function(userData){
        $scope.currentUser = userData.User;

        if($scope.currentUser.Auth == $cookieStore.get('auth')){
          $scope.loggedIn = true;
          deferred.resolve(true);
        } 
        else{
          $scope.loggedIn = false;
          deferred.resolve(false);
        }
      });
    }
    return deferred.promise;
  }
  
  $scope.$on('event:google-plus-signin-success', function (event, authResult) {
    // User successfully authorized the G+ App!
    console.log('SessionController module:: User successfully authorized the G+ App!');
    Session.fetchUserInfo({ access_token: authResult.access_token }).$promise.then(function(userInfo) {
      Session.fetchUser({  access_token: authResult.access_token, 
                           id:userInfo.id, 
                           name:userInfo.displayName, 
                           email:userInfo.emails[0].value } ).$promise.then(function(userData) {
       $scope.currentUser = userData.User;
       console.log('current user: ', $scope.currentUser);
       
       $cookieStore.put('access_token', authResult.access_token);
       $cookieStore.put('auth', $scope.currentUser.Auth);
       $cookieStore.put('user_id', $scope.currentUser.Id);
       
       $scope.loggedIn = true;
       
       $location.path('/home');
      });
    });
  });
  $scope.$on('event:google-plus-signin-failure', function (event, authResult) {
    // User has not authorized the G+ App!
    console.log('Not signed into Google Plus.');
  });
  
  $scope.disconnect = function(){
    console.log('SessionController module:: disconnect');

    Session.logout({ token: $cookieStore.get('access_token') });
    
    $cookieStore.remove('auth');
    $cookieStore.remove('access_token');
    $cookieStore.remove('user_id');
    
    $scope.currentUser = undefined;
    $scope.loggedIn = false;
    
    $location.path('/');
  };
  
  $scope.signinWithTwitter = function(oauthToken, oauthVerifier) {
    // User successfully authorized via Twitter!
    console.log('User successfully authorized via Twitter!');

    Session.fetchTwitterUser({ oauth_token: oauthToken, oauth_verifier: oauthVerifier }).$promise.then(function(userData) {
      $scope.currentUser = userData.User;
       console.log('current user: ', $scope.currentUser);
       
       $cookieStore.put('access_token', oauthToken);
       $cookieStore.put('auth', $scope.currentUser.Auth);
       $cookieStore.put('user_id', $scope.currentUser.Id);
       
       $scope.loggedIn = true;
       
       $location.path('/home');
    });
  }
}]);

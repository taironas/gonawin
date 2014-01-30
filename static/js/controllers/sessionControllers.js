'use strict';

var sessionControllers = angular.module('sessionControllers', []);

sessionControllers.controller('SessionCtrl', ['$scope', '$location', '$cookieStore', '$q', 'Session', 
  function ($scope, $location, $cookieStore, $q, Session) {
    console.log('SessionController module');
    $scope.currentUser  = undefined;
    $scope.loggedIn = false;
    
    function getCurrentUser() {
			var deferred = $q.defer();
			if($scope.loggedIn && !$scope.currentUser)
			{
				$scope.currentUser = User.get({ id:$cookieStore.get('user_id') });
			}
			deferred.resolve(currentUser);
			
			return deferred.promise;
		}
		function getUserLoggedIn() {
			var deferred = $q.defer();

			if($scope.loggedIn) {
				deferred.resolve(true);
			}

			if(!$cookieStore.get('access_token') || !$cookieStore.get('user_id')) {
				$scope.loggedIn = false;
				deferred.resolve(false);
			} 
			else {
				User.get({ id:$cookieStore.get('user_id') }).$promise.then(function(result){
					$scope.currentUser = result;

					if($scope.currentUser.User.Auth == $cookieStore.get('auth'))
					{
						$scope.loggedIn = true;
						deferred.resolve(true);
					} 
					else 
					{
						$scope.loggedIn = false;
						deferred.resolve(false);
					}
				});
			}
			
			return deferred.promise;
		}
    
    $scope.initSession = function(){
      console.log('SessionController module:: initSession');
      
      getUserLoggedIn().then(function(result) {
        $scope.loggedIn = result;
        if($scope.loggedIn) {
          getCurrentUser().then(function(result) {
            console.log('current user: ', $scope.currentUser);
            $scope.currentUser = result;
          });
        }
      });
    }
    
    $scope.$on('event:google-plus-signin-success', function (event, authResult) {
      // User successfully authorized the G+ App!
      console.log('SessionController module:: User successfully authorized the G+ App!');
      Session.fetchUserInfo({ access_token: authResult.access_token }).$promise.then(function(userInfo) {
        Session.fetchUser({  access_token: authResult.access_token, 
                                    id:userInfo.id, 
                                    name:userInfo.displayName, 
                                    email:userInfo.emails[0].value } ).$promise.then(function(user) {
          $scope.currentUser = user;
          console.log('current user: ', $scope.currentUser);
          
          $cookieStore.put('access_token', authResult.access_token);
					$cookieStore.put('auth', $scope.currentUser.Auth);
					$cookieStore.put('user_id', $scope.currentUser.Id);
          
          $scope.loggedIn = true;
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
}]);

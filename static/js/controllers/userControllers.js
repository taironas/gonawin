'use strict';

var userControllers = angular.module('userControllers', []);

userControllers.controller('UserListCtrl', ['$scope', 'User', function($scope, User) {
    $scope.users = User.query();
}]);

userControllers.controller('UserShowCtrl', ['$scope', '$routeParams', 'User', '$location', function($scope, $routeParams, User, $location) {
    $scope.user = User.get({ id:$routeParams.id });

    $scope.acceptTeamRequest = function(){
	console.log('User show controller:: accept team Request');
    };
    $scope.denyTeamRequest = function(){
	console.log('User show controller:: deny team Request');
    };
}]);

userControllers.controller('UserEditCtrl', ['$scope', '$routeParams', 'User', 'SessionService', '$location', function($scope, $routeParams, User, SessionService, $location) {
    $scope.user = undefined;
    
    $scope.loadCurrentUser = function() {
	$scope.user = SessionService.getCurrentUser();
    }
    
    $scope.updateUser = function() {
	User.update({ id:$scope.user.Id }, $scope.user,
		    function(){
			$location.path('/settings/edit-profile/');
		    },
		    function(err) {
			console.log('update failed: ', err.data);
		    });
    }
}]);

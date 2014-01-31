'use strict';

var userControllers = angular.module('userControllers', []);

userControllers.controller('UserListCtrl', ['$scope', 'User', function($scope, User) {
    $scope.users = User.query();
}]);

userControllers.controller('UserShowCtrl', ['$scope', '$routeParams', 'User', 'Team', '$location', function($scope, $routeParams, User, Team, $location) {
    $scope.userData = User.get({ id:$routeParams.id, including: "Teams TeamRequests Tournaments" });

    $scope.acceptTeamRequest = function(request){
	console.log('User show controller:: accept team Request');
	console.log('req: ', request);
	Team.allowRequest({requestId:request.Id});
    };
    $scope.denyTeamRequest = function(request){
	console.log('User show controller:: deny team Request');
	console.log('req: ', request);
	Team.denyRequest({requestId:request.Id})
    };
}]);

userControllers.controller('UserEditCtrl', ['$scope', '$routeParams', 'User', 'SessionService', '$location', function($scope, $routeParams, User, SessionService, $location) {
    $scope.userData = undefined;
    
    $scope.loadCurrentUser = function() {
	$scope.userData = SessionService.getCurrentUser();
    }
    
    $scope.updateUser = function() {
	User.update({ id:$scope.user.Id }, $scope.userData.User,
		    function(){
			$location.path('/settings/edit-profile/');
		    },
		    function(err) {
			console.log('update failed: ', err.data);
		    });
    }
}]);

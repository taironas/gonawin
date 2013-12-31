'use strict';

var userControllers = angular.module('userControllers', []);

userControllers.controller('UserListCtrl', ['$scope', 'User', function($scope, User) {
	$scope.users = User.query();
}]);

userControllers.controller('UserShowCtrl', ['$scope', '$routeParams', 'User', '$location', function($scope, $routeParams, User, $location) {
	$scope.user = User.get({ id:$routeParams.id });
}]);

userControllers.controller('UserEditCtrl', ['$scope', '$routeParams', 'User', '$location', function($scope, User, $location) {
	var currentUserId = 1; /* need to retrive the id of the current user*/
	
	$scope.user = User.get({ id:currentUserId  });

	$scope.updateUser = function() {
		var user = User.get({ id:currentUserId });
		User.update({ id:currentUserId }, $scope.user,
			function(){
				$location.path('/users/update/' + currentUserId);
			},
		function(err) {
			console.log('update failed: ', err.data);
		});
	}
}]);

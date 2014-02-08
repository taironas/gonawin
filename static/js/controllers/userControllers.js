'use strict';

var userControllers = angular.module('userControllers', []);

userControllers.controller('UserListCtrl', ['$scope', 'User', function($scope, User) {
  $scope.users = User.query();
}]);

userControllers.controller('UserShowCtrl', ['$scope', '$routeParams', 'User', 'Team', function($scope, $routeParams, User, Team) {
  $scope.userData = User.get({ id:$routeParams.id, including: "Teams TeamRequests Tournaments" },
			     function(data){},
			     function(err){
			       console.log('get user failed: ', err.data);
			       $scope.messageDanger = err.data;
			     });

  $scope.acceptTeamRequest = function(request){
    console.log('User show controller:: accept team Request');
    console.log('req: ', request);
    Team.allowRequest({requestId:request.Id},
		      function(data){},
		      function(err){
			console.log('allow request failed: ', err.data);
			$scope.messageDanger = err.data;
		      });
  };
  $scope.denyTeamRequest = function(request){
    console.log('User show controller:: deny team Request');
    console.log('req: ', request);
    Team.denyRequest({requestId:request.Id},
		     function(data){},
		     function(err){
		       console.log('deny request failed: ', err.data);
		       $scope.messageDanger = err.data;
		     });
  };
}]);

userControllers.controller('UserEditCtrl', ['$scope', 'User', '$location', function($scope, User, $location) {
  
  $scope.updateUser = function() {
    User.update({ id:$scope.currentUser.Id }, $scope.currentUser,
		function(){ 
		  $location.path('/settings/edit-profile/'); 
		},
		function(err) { 
		  console.log('update failed: ', err.data);
		  $scope.messageDanger = err.data;
		});
  }
}]);

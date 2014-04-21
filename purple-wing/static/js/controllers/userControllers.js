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

// User edit controller. Use this controller to edit the current user data.
userControllers.controller('UserEditCtrl', ['$scope', '$rootScope', 'User', '$location', function($scope, $rootScope, User, $location) {
    
    $scope.updateUser = function() {
	User.update({ id:$rootScope.currentUser.User.Id}, $scope.currentUser,
		    function(response){
			$rootScope.messageInfo = response.MessageInfo;
		    },
		    function(err) { 
			console.log('update failed: ', err.data);
			$scope.messageDanger = err.data;
		    });
    }
}]);

userControllers.controller('UserScoresCtrl', ['$scope', '$routeParams', 'User', 'Team', function($scope, $routeParams, User, Team) {
  console.log('User Scores controller');
  $scope.userData = User.get({ id:$routeParams.id, including: "Teams TeamRequests Tournaments" },
			     function(data){},
			     function(err){
			       console.log('get user failed: ', err.data);
			       $scope.messageDanger = err.data;
			     });

  $scope.scoreData = User.scores({id:$routeParams.id},
				 function(response){
				   console.log('response: ', response);
				 },
				 function(err){
				   console.log('user scores failed', err.data)
				 });
}]);

// UserCardCtrl: fetch data of a particular user.
userControllers.controller('UserCardCtrl', ['$scope', 'User', function($scope, User) {
    console.log('User card controller:');
    console.log('user ID: ', $scope.$parent.user.Id);
    $scope.userData = User.get({ id:$scope.$parent.user.Id});

    $scope.userData.$promise.then(function(userData){
      $scope.user = userData.User;
    });
}]);

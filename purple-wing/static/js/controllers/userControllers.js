'use strict';

var userControllers = angular.module('userControllers', []);

userControllers.controller('UserListCtrl', ['$scope', 'User', function($scope, User) {
  $scope.users = User.query();
}]);

userControllers.controller('UserShowCtrl', ['$scope', '$routeParams', 'User', 'Team', function($scope, $routeParams, User, Team) {
  $scope.userData = User.get({ id:$routeParams.id, including: "Teams TeamRequests Tournaments Invitations" },
			     function(data){},
			     function(err){
			       console.log('get user failed: ', err.data);
			       $scope.messageDanger = err.data;
			     });
    $scope.userData.$promise.then(function(response){
	var lenInvite = 0;
	if($scope.userData.Invitations != undefined){
	    lenInvite = $scope.userData.Invitations.length;
	}
	for(var i = 0; i < lenInvite; i++){
	    $scope.userData.Invitations[i].handled = true;
	}
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
    
    $scope.acceptInvitation = function(invitation, index){
	if(!$scope.userData.Invitations[index].handled){
	    return;
	}
	User.allowInvitation({ id:$routeParams.Id, teamId:invitation.Id},
			     function(data){
				 console.log('user allow invitation');
				 $scope.messageInfo = data.MessageInfo;
				 $scope.userData.Invitations[index].handled = false;
			     },
			     function(err){
				 console.log('allow invitation failed: ', err.data);
				 $scope.messageDanger = err.data;
			     });
    };
    
    $scope.denyInvitation = function(invitation, index){
	if(!$scope.userData.Invitations[index].show){
	    return;
	}
	User.denyInvitation({ id:$routeParams.Id, teamId:invitation.Id},
			    function(data){
				console.log('user deny invitation');
				$scope.messageInfo = data.MessageInfo;
				$scope.userData.Invitations[index].handled = false;
			    },
			    function(err){
				console.log('deny invitation failed: ', err.data);
				$scope.messageDanger = err.data;
			    });
    };
    
}]);

// User edit controller. Use this controller to edit the current user data.
userControllers.controller('UserEditCtrl', ['$scope', '$rootScope', '$location', 'User', 'sAuth', 
  function($scope, $rootScope, $location, User, sAuth) {
    $scope.updateUser = function() {
      User.update({ id:$rootScope.currentUser.User.Id}, $scope.currentUser,
        function(response){
          $rootScope.messageInfo = response.MessageInfo;
        },
        function(err) { 
        console.log('update failed: ', err.data);
        $scope.messageDanger = err.data;
      });
    };
    
    $scope.deleteUser = function() {
      if(confirm('Are you sure?')){
        User.delete({ id:$rootScope.currentUser.User.Id},
        function(response){
          // reset rootScope variables
          $rootScope.currentUser = undefined;
          $rootScope.isLoggedIn = false;
        
          sAuth.clearCookies();
          $location.path('/welcome');
        },
        function(err) {
          $scope.messageDanger = err.data;
          console.log('delete failed: ', err.data);
        });
      }
    };
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

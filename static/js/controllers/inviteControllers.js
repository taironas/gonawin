'use strict';

var inviteControllers = angular.module('inviteControllers', []);

inviteControllers.controller('InviteCtrl', ['$scope', 'Invite', 'SessionService', function($scope, Invite, SessionService) {
    
    $scope.inviteFriends = function() {
	console.log('invite friends');
	console.log($scope.currentUser.User);
	console.log($scope.invite.emails);
	Invite.send({emails: $scope.invite.emails},
		    function(result){
			console.log('invite successfull: ');
		    }, function(err){
			console.log('error: ',err);
		    });
    };
}]);

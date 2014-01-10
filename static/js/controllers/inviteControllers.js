'use strict';

var inviteControllers = angular.module('inviteControllers', []);

inviteControllers.controller('InviteCtrl', ['$scope', 'Invite', 'SessionService', function($scope, Invite, SessionService) {
	$scope.currentUser = SessionService.getCurrentUser();
	
	$scope.inviteFriends = function() {
		Invite.send($scope.currentUser, $scope.invite.emails);
	}
}]);
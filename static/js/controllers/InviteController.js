'use strict';

var inviteController = angular.module('inviteController', []);

inviteController.controller('InviteCtrl', ['$scope', 'Invite', 'SessionService', function($scope, Invite, SessionService) {
	$scope.currentUser = SessionService.getCurrentUser();
	
	$scope.inviteFriends = function() {
		Invite.send($scope.currentUser, $scope.invite.emails);
	}
}]);

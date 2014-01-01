'use strict';

var inviteController = angular.module('inviteController', []);

inviteController.controller('InviteCtrl', ['$scope', '$location', 'SessionService', function($scope, $location, SessionService) {
	$scope.currentUser = SessionService.getCurrentUser();
	
	$scope.inviteFriends = function() {
		console.log('Invite friends');
		console.log('Emails: ', $scope.emails);
	}
}]);

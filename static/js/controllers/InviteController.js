'use strict';

var inviteController = angular.module('inviteController', []);

inviteController.controller('InviteCtrl', ['$scope', '$location', function($scope, $location) {

	$scope.inviteFriends = function() {
		console.log('Invite friends');
	}
}]);

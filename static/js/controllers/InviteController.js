'use strict';

var inviteController = angular.module('inviteController', []);

inviteController.controller('InviteCtrl', ['$scope', '$http', '$location', 'SessionService', function($scope, $http, $location, SessionService) {
	$scope.currentUser = SessionService.getCurrentUser();
	
	$scope.inviteFriends = function() {
		$http({
			method: 'POST',
			url: '/j/invite',
			contentType: 'application/json',
			params:{ emails: $scope.invite.emails, name: $scope.currentUser.Name } }).
			success(function() {
				console.log('invite friends successfully');
			}).
			error(function() { console.log('invite friends failed') });
	}
}]);

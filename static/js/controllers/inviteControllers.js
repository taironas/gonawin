'use strict';

var inviteControllers = angular.module('inviteControllers', []);

inviteControllers.controller('InviteCtrl', ['$scope', 'Invite', function($scope, Invite) {
  $scope.inviteFriends = function() {
    console.log('invite friends');
    Invite.send({emails: $scope.invite.emails},
		function(result){
		  console.log('invite successfull');
		}, 
		function(err){
		  console.log('invite failed: ',err.data);
		  $scope.messageDanger = err.data;
		});
  };
}]);

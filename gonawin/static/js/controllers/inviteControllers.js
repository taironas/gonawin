'use strict';

var inviteControllers = angular.module('inviteControllers', []);

inviteControllers.controller('InviteCtrl', ['$scope', '$rootScope', 'Invite', function($scope, $rootScope, Invite) {
  $rootScope.title = 'gonawin - Invite';
  
  $scope.inviteFriends = function() {
    console.log('invite friends');
    Invite.send({emails: $scope.invite.emails},
		function(response){
		  $scope.messageInfo = response.MessageInfo;
		  console.log('invite successfull');
		}, 
		function(err){
		  console.log('invite failed: ',err.data);
		  $scope.messageDanger = err.data;
		});
  };
}]);

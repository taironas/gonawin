'use strict';

purpleWingApp.controller('InviteController',
		     function InviteController($scope, inviteData, $location, $routeParams){
			 $scope.inviteData = inviteData.getData();
		     }
		    );

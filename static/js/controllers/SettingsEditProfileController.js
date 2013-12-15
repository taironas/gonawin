'use strict';

purpleWingApp.controller('SettingsEditProfileController',
		     function SettingsEditProfileController($scope, profileData, $location, $routeParams){
			 console.log('setting edit profile controller');
			 $scope.profileData = profileData.getData();
		     }
		    );

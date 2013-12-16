'use strict';
function google_plus_sign_in_render(){
  angular.element($('#login-button')).scope().render();
}

purpleWingApp.controller('SessionController', function ($scope, $http, $cookieStore) {
	console.log('SessionController module');
	$scope.loggedIn = false;
	$scope.currentUser = undefined;
	render();
	function render(){
		gapi.signin.render('login-button', {
	    'callback' : function(authResult){
				$scope.checkAuth(authResult);
	    },
	    'clientid': '36368233290',
	    'cookiepolicy': 'single_host_origin',
	    'requestvisibleactions': 'http://schemas.google.com/AddActivity',
	    'scope': 'https://www.googleapis.com/auth/plus.login https://www.googleapis.com/auth/userinfo.email',
	    width: 'wide'
		});
  };
  $scope.checkAuth = function(authResult){
		if(authResult.access_token){
	    $scope.loggedIn = true;
	    $cookieStore.put('access_token', authResult.access_token );
	    $scope.login();
	    $scope.$apply();

	    return true;
		}
		else {
	    $scope.loggedIn = false;
			$scope.currentUser = undefined;
	    $scope.$apply();
	    return false;
		}
  };
  $scope.login = function(){
		var peopleUrl = 'https://www.googleapis.com/plus/v1/people/me?access_token='+$cookieStore.get('access_token');
		$.ajax({
	    type: 'GET',
	    url: peopleUrl,
	    contentType: 'application/json',
	    dataType: 'jsonp',
	    success: function(result) {
				$scope.completeLogin(result)
	    }
		});
  };
  $scope.completeLogin = function(userInfo){
		var authUrl = window.location.href + 'ng/auth/google'

		$.ajax({
	    type: 'GET',
	    url: authUrl,
	    contentType: 'application/json',
	    dataType: 'jsonp',
	    success: function(result) {
				console.log('completeLogin successfully');
				console.log(result);
	    },
	    data: {
				access_token: $cookieStore.get('access_token'),
				id: userInfo.id,
				name: userInfo.displayName,
				email: userInfo.emails[0].value
	    }
		});
  };

  $scope.disconnect = function(access_token){
		var revokeUrl = 'https://accounts.google.com/o/oauth2/revoke?token='+$cookieStore.get('access_token');
		$.ajax({
	    type: 'GET',
	    url: revokeUrl,
	    async: false,
	    contentType: "application/json",
	    dataType: 'jsonp',
	    success: function(nullResponse){
				$scope.loggedIn = false;
				$scope.currentUser = undefined;
				$scope.$apply();
	    }
		})
  };
});

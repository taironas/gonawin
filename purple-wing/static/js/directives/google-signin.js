'use strict';

angular.module('directive.googlesignin', []).
  directive('googleSignin', function (Session, $rootScope) {
    return {
      restrict: 'E',
      template: '<button class="btn btn-danger">Sign in with Google</button>',
      replace: true,
      link: function (scope, element, attrs) {
        element.bind("click", function(){
          console.log("Sign in with Google has started...");
          Session.authenticateWithGoogle().$promise.then(function(data){
            console.log('authenticateWithGoogle: data = ', data);
            $rootScope.$broadcast('event:google-signin-success', data);
          }).then(function(error) {
            console.log('authenticateWithGoogle: error = ', error);
            $rootScope.$broadcast('event:google-signin-failure', error);
          });
        })
      }
    };
  });

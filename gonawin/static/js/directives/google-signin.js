'use strict';

angular.module('directive.googlesignin', []).
  directive('googleSignin', function (Session, $rootScope) {
    return {
      restrict: 'E',
      template: '<button class="btn btn-block btn-social btn-google-plus btn-lg"><i class="fa fa-google-plus"></i> Sign in with Google</button>',
      replace: true,
      link: function (scope, element, attrs) {
        element.bind("click", function(){
          console.log("Sign in with Google has started...");
          Session.fetchGoogleLoginUrl().$promise.then(function(data){
            window.location.replace(data.Url);
          }).then(function(error) {
            console.log('fetchGoogleLoginUrl: error = ', error);
          });
        })
      }
    };
  });

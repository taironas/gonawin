'use strict';

angular.module('directive.facebooksignin', []).
  directive('facebookSignin', function (Session, $location) {
    return {
      restrict: 'E',
      template: '<button class="btn btn-block btn-social btn-facebook"><i class="fa fa-facebook"></i> Signin with Facebook</button>',
      replace: true,
      link: function (scope, element, attrs) {
        element.bind("click", function(){
          console.log("Sign in with Facebook has started...");
          FB.login(function(response) {
            if (response.session) {
              //watchLoginChange event has been triggered and
              //will be handle in auth service.
            }
          }, {scope:'email'});
        })
      }
    };
  });

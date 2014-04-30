'use strict';

/*
 * angular-google-plus-directive v0.0.1
 * ♡ CopyHeart 2013 by Jerad Bitner http://jeradbitner.com
 * Copying is an act of love. Please copy.
 *
 * Modified version to be able to set the client id via a promise
 */


angular.module('directive.googleplussignin', []).
  directive('googlePlusSignin', ['Session', function (Session) {
    return {
      restrict: 'E',
      transclude: true,
      template: '<span></span>',
      replace: true,
      link: function (scope, element, attrs, ctrl, linker) {

        Session.serviceIds().$promise.then(function(response){
          attrs.clientid = response.GooglePlusClientId + '.apps.googleusercontent.com';
          attrs.$set('data-clientid', attrs.clientid);
          
          // Some default values, based on prior versions of this directive
          var defaults = {
            callback: 'signinCallback',
            cookiepolicy: 'single_host_origin',
            scope: 'https://www.googleapis.com/auth/plus.login https://www.googleapis.com/auth/userinfo.email',
            width: 'wide'
          };
          
          defaults.clientid = attrs.clientid;

          // Provide default values if not explicitly set
          angular.forEach(Object.getOwnPropertyNames(defaults), function(propName) {
            if (!attrs.hasOwnProperty('data-' + propName)) {
              attrs.$set('data-' + propName, defaults[propName]);
            }
          });
          
          // Asynchronously load the G+ SDK when service IDs are ready.
          var po = document.createElement('script'); po.type = 'text/javascript'; po.async = true;
          po.src = 'https://apis.google.com/js/client:plusone.js';
          var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(po, s);
          
          linker(function(el, tScope){
            po.onload = function() {
              element.addClass('customGPlusSignIn');
              if (el.length) {
                element.append(el);
              }
              gapi.signin.render(element[0], defaults);
            };
          });
        });
      }
    };
  }]).run(['$window','$rootScope',function($window,$rootScope){
    $window.signinCallback = function (authResult) {
      if(authResult && authResult.access_token){
        $rootScope.$broadcast('event:google-plus-signin-success',authResult);
      }
      else{
        $rootScope.$broadcast('event:google-plus-signin-failure',authResult);
      }
    };
  }]);


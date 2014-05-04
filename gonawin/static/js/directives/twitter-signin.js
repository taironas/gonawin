'use strict';

angular.module('directive.twittersignin', []).
  directive('twitterSignin', function (Session, $location) {
    return {
      restrict: 'E',
      template: '<button class="btn btn-social btn-block btn-twitter btn-lg"><i class="fa fa-twitter"></i> Sign in with Twitter</button>',
      replace: true,
      link: function (scope, element, attrs) {
        element.bind("click", function(){
          console.log("Sign in with Twitter has started...");
          Session.authenticateWithTwitter().$promise.then(function(data){
            console.log('authenticateWithTwitter: data = ', data);
            window.location.replace('https://api.twitter.com/oauth/authorize?oauth_token='+data.OAuthToken);
          });
        })
      }
    };
  });

'use strict';

var purpleWingApp = angular.module('purpleWingApp', [
  'ngSanitize',
  'ngRoute',
  'ngResource',
  'ngCookies',
  'directive.g+signin',
  'directive.twittersignin',
  'directive.formValidation',
  'directive.joinButton',
  '$strap.directives',
  'filter.fromNow',
  'filter.reverse',

  'rootControllers',
  'navigationControllers',
  'dashboardControllers',
  'activitiesControllers',
  'userControllers',
  'teamControllers',
  'tournamentControllers',
  'inviteControllers',
  
  'dataServices',
  'authServices'
]);

purpleWingApp.factory('notFoundInterceptor', ['$q', '$location', function($q, $location){
  return {
    response: function(response) {
      return response || $q.when(response);
    },

    responseError: function(response) {
      if (response && response.status === 404) {
        $location.path('/404');
      }
      return $q.reject(response);
    }
  };
}]);

purpleWingApp.config(['$routeProvider', '$httpProvider',
  function($routeProvider, $httpProvider) {
    $routeProvider.
      when('/welcome', { templateUrl: 'templates/welcome.html', requireLogin: false }).
      when('/', { templateUrl:  'templates/home.html', controller: 'RootCtrl', requireLogin: true }).
      when('/about', { templateUrl: 'templates/about.html', requireLogin: false }).
      when('/contact', { templateUrl: 'templates/contact.html', requireLogin: false }).
      when('/users/', { templateUrl: 'templates/users/index.html', controller: 'UserListCtrl', requireLogin: true }).
      when('/users/:id', { templateUrl: 'templates/users/show.html', controller: 'UserShowCtrl', requireLogin: true }).
      when('/users/:id/scores', {templateUrl: 'templates/users/scores.html', controller: 'UserScoresCtrl', requireLogin: true}).
      when('/teams', { templateUrl: 'templates/teams/index.html', controller: 'TeamListCtrl', requireLogin: true }).
      when('/teams/new', { templateUrl: 'templates/teams/new.html', controller: 'TeamNewCtrl', requireLogin: true }).
      when('/teams/:id', { templateUrl: 'templates/teams/show.html', controller: 'TeamShowCtrl', requireLogin: true }).
      when('/teams/edit/:id', { templateUrl: 'templates/teams/edit.html', controller: 'TeamEditCtrl', requireLogin: true }).
      when('/teams/search', { templateUrl: 'templates/teams/index.html', controller: 'TeamSearchCtrl', requireLogin: true}).
      when('/teams/:id/ranking', { templateUrl: 'templates/teams/ranking.html', controller: 'TeamRankingCtrl', requireLogin: true }).
      when('/teams/:id/accuracies', {templateUrl: 'templates/teams/accuracies.html', controller: 'TeamAccuraciesCtrl', requireLogin: true}).
      when('/teams/:id/accuracies/:tournamentId', {templateUrl: 'templates/teams/accuracy.html', controller: 'TeamAccuracyByTournamentCtrl', requireLogin: true}).
      when('/teams/:id/prices', {templateUrl: 'templates/teams/prices.html', controller: 'TeamPricesCtrl', requireLogin: true}).
      when('/teams/:id/prices/:tournamentId', {templateUrl: 'templates/teams/price.html', controller: 'TeamPriceByTournamentCtrl', requireLogin: true}).
      when('/teams/:id/prices/edit/:tournamentId', {templateUrl: 'templates/teams/priceedit.html', controller: 'TeamPriceEditByTournamentCtrl', requireLogin: true}).

      when('/tournaments', { templateUrl: 'templates/tournaments/index.html', controller: 'TournamentListCtrl', requireLogin: true }).
      when('/tournaments/new', { templateUrl: 'templates/tournaments/new.html', controller: 'TournamentNewCtrl', requireLogin: true }).
      when('/tournaments/:id', { templateUrl: 'templates/tournaments/show.html', controller: 'TournamentShowCtrl', requireLogin: true }).
      when('/tournaments/edit/:id', { templateUrl: 'templates/tournaments/edit.html', controller: 'TournamentEditCtrl', requireLogin: true }).
      when('/tournaments/edit/:id', { templateUrl: 'templates/tournaments/edit.html', controller: 'TournamentEditCtrl', requireLogin: true }).
      when('/tournaments/search', { templateUrl: 'templates/tournaments/index.html', controller: 'TournamentSearchCtrl', requireLogin: true }).
      when('/tournaments/:id/calendar', { templateUrl: 'templates/tournaments/calendar.html', controller: 'TournamentCalendarCtrl', requireLogin: true }).
      when('/tournaments/:id/matches/firststage', { templateUrl: 'templates/tournaments/firststage.html', controller: 'TournamentFirstStageCtrl', requireLogin: true }).
      when('/tournaments/:id/matches/secondstage', { templateUrl: 'templates/tournaments/secondstage.html', controller: 'TournamentSecondStageCtrl', requireLogin: true }).
      when('/tournaments/:id/matches/secondstage2', { templateUrl: 'templates/tournaments/secondstage2.html', controller: 'TournamentSecondStageCtrl', requireLogin: true }).
      when('/tournaments/:id/predict/', { templateUrl: 'templates/tournaments/predict.html', controller: 'TournamentPredictCtrl', requireLogin: true }).
      when('/tournaments/:id/ranking', { templateUrl: 'templates/tournaments/ranking.html', controller: 'TournamentRankingCtrl', requireLogin: true }).
      // this should be an admin page
      when('/tournaments/:id/matches/setresults', { templateUrl: 'templates/tournaments/setresults.html', controller: 'TournamentSetResultsCtrl', requireLogin: true }).

      when('/settings/edit-profile', { templateUrl: 'templates/users/edit.html', controller: 'UserEditCtrl', requireLogin: true }).
      when('/settings/networks', { templateUrl: 'templates/settings/networks.html', requireLogin: true }).
      when('/settings/email', { templateUrl: 'templates/settings/email.html', requireLogin: true }).
      when('/invite', { templateUrl: 'templates/invite.html', controller: 'InviteCtrl', requireLogin: true }).
      when('/404', { templateUrl: 'static/templates/404.html' }).
      otherwise( {redirectTo: '/'});
    
    $httpProvider.interceptors.push('notFoundInterceptor');
}]);

purpleWingApp.run(['$rootScope', '$location', '$window', 'sAuth', 'Session', 'User', function($rootScope, $location, $window, sAuth, Session, User) {
  $rootScope.currentUser = undefined;
  $rootScope.isLoggedIn = false;

  $window.fbAsyncInit = function() {
    // Executed when the SDK is loaded
    FB.init({
      appId: '232160743609875',
      channelUrl: 'static/templates/channel.html',
      status: true, /*Set if you want to check the authentication status at the start up of the app */
      cookie: true,
      xfbml: true
    });

    sAuth.watchLoginChange();
  };
  
  (function(d){
    // load the Facebook javascript SDK
    var js,
    id = 'facebook-jssdk',
    ref = d.getElementsByTagName('script')[0];

    if (d.getElementById(id)) {
      return;
    }

    js = d.createElement('script');
    js.id = id;
    js.async = true;
    js.src = "//connect.facebook.net/en_US/all.js";

    ref.parentNode.insertBefore(js, ref);

  }(document));
  
  $rootScope.$on("$routeChangeStart", function(event, next, current) {
    console.log('routeChangeStart, requireLogin = ', next.requireLogin);
    console.log('routeChangeStart, current user = ', $rootScope.currentUser);
    console.log('routeChangeStart, isLoggedIn = ', $rootScope.isLoggedIn);
    
    $rootScope.isLoggedIn = sAuth.isLoggedIn();
    if($location.$$path === '/auth/twitter/callback')
    {
      sAuth.signinWithTwitter(($location.search()).oauth_token, ($location.search()).oauth_verifier);
    } else {
      // Everytime the route in our app changes check authentication status.
      // Get current user only if we are logged in.
      if( $rootScope.isLoggedIn && (undefined == $rootScope.currentUser) ) {
        $rootScope.currentUser = User.get({ id:sAuth.getUserID() });
        console.log('routeChangeStart, current user = ', $rootScope.currentUser);
      }
      // Redirect user to root if he tries to go on welcome page and he is logged in.
      if( $location.path() === '/welcome' && $rootScope.isLoggedIn ) {
        console.log('routeChangeStart, redirect to root');
        $location.path('/');
      }
      // Redidrect to welcome if route requires to be logged in and user is not logged in.
      if ( next.requireLogin && (undefined == $rootScope.currentUser) ) {
        console.log('routeChangeStart, redirect to welcome');
        $location.path('/welcome');
      }
    }
    console.log('end of routeChangeStart');
  });
  
  $rootScope.$on('event:google-plus-signin-success', function (event, authResult) {
    // User successfully authorized the G+ App!
    console.log('event:google-plus-signin-success');
    Session.fetchUserInfo({ access_token: authResult.access_token }).$promise.then(function(userInfo) {
      $rootScope.currentUser = Session.fetchUser({  
        access_token: authResult.access_token,
        provider: 'google',
        id:userInfo.id,
        name:userInfo.displayName,
        email:userInfo.emails[0].value } );
      $rootScope.currentUser.$promise.then(function(currentUser){
        console.log('event:google-plus-signin-success: current user = ', currentUser);
        sAuth.storeCookies(authResult.access_token, currentUser.User.Auth, currentUser.User.Id);
        $rootScope.isLoggedIn = true;
        $location.path('/');
      });
    });
  });
  $rootScope.$on('event:google-plus-signin-failure', function (event, authResult) {
    // User has not authorized the G+ App!
  });
}]);

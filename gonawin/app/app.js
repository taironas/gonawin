'use strict';

angular.module('gonawingApp', [
  'ngSanitize',
  'ngRoute',
  'ngResource',
  'ngCookies',
  /* directives */
  'directive.googleplussignin',
  'directive.twittersignin',
  'directive.googlesignin',
  'directive.facebooksignin',
  'directive.formValidation',
  'directive.formFocus',
  'directive.joinButton',
  'directive.addButton',
  'directive.activities',
  'mgcrea.ngStrap',
  /* filters */
  'filter.fromNow',
  'filter.moment',
  'filter.reverse',
  /* controllers */
  'rootControllers',
  'navigationControllers',
  'dashboardControllers',
  'activitiesControllers',
  'userControllers',
  'teamControllers',
  'tournamentControllers',
  'inviteControllers',
  'searchControllers',
  /* services */
  'userService',
  'authService',
  'sessionService',
  'activitiesService',
  'inviteService',
  'teamService',
  'tournamentService',
  'themeSelectorService'
]);

angular.module('gonawingApp').factory('notFoundInterceptor', ['$q', '$location', function($q, $location) {
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

angular.module('gonawingApp').config(['$routeProvider', '$httpProvider',
  function($routeProvider, $httpProvider) {
    $routeProvider.
      when('/welcome', { templateUrl: 'components/home/welcome.html', requireLogin: false }).
      when('/getting-started', { templateUrl: 'components/home/getting-started.html', requireLogin: false }).
      when('/', { templateUrl:  'components/home/home.html', controller: 'RootCtrl', requireLogin: true }).
      when('/about', { templateUrl: 'shared/about.html', requireLogin: false }).
      when('/search', { templateUrl: 'components/search/search.html', controller: 'SearchCtrl', requireLogin: true }).

      when('/users/', { templateUrl: 'components/user/index.html', controller: 'UserListCtrl', requireLogin: true }).
      when('/users/:id', { templateUrl: 'components/user/show.html', controller: 'UserShowCtrl', requireLogin: true }).

      when('/teams', { templateUrl: 'components/team/index.html', controller: 'TeamListCtrl', requireLogin: true }).
      when('/teams/new', { templateUrl: 'components/team/new.html', controller: 'TeamNewCtrl', requireLogin: true }).
      when('/teams/:id', { templateUrl: 'components/team/show.html', controller: 'TeamShowCtrl', requireLogin: true }).
      when('/teams/edit/:id', { templateUrl: 'components/team/edit.html', controller: 'TeamEditCtrl', requireLogin: true }).
      when('/teams/invite/:id', { templateUrl: 'components/team/invite.html', controller: 'TeamInviteCtrl', requireLogin: true }).

      when('/tournaments', { templateUrl: 'components/tournament/index.html', controller: 'TournamentListCtrl', requireLogin: true }).
      when('/tournaments/new', { templateUrl: 'components/tournament/new.html', controller: 'TournamentNewCtrl', requireLogin: true }).
      when('/tournaments/:id', { templateUrl: 'components/tournament/show.html', controller: 'TournamentShowCtrl', requireLogin: true, reloadOnSearch: false }).
      when('/tournaments/edit/:id', { templateUrl: 'components/tournament/edit.html', controller: 'TournamentEditCtrl', requireLogin: true }).

      when('/settings/edit-profile', { templateUrl: 'components/user/edit.html', controller: 'UserEditCtrl', requireLogin: true }).
      when('/settings/networks', { templateUrl: 'components/settings/networks.html', requireLogin: true }).
      when('/settings/email', { templateUrl: 'components/settings/email.html', requireLogin: true }).
      when('/invite', { templateUrl: 'components/invite/invite.html', controller: 'InviteCtrl', requireLogin: true }).
      when('/404', { templateUrl: 'shared/404.html' }).
      otherwise( {redirectTo: '/'});

    $httpProvider.interceptors.push('notFoundInterceptor');
}]);

angular.module('gonawingApp').run(['$rootScope', '$location', '$window', '$cookieStore', 'sAuth', 'Session', 'User', function($rootScope, $location, $window, $cookieStore, sAuth, Session, User) {
    $rootScope.title = 'gonawin';

    $rootScope.currentUser = undefined;
    $rootScope.isLoggedIn = false;
    $rootScope.serviceIds = Session.serviceIds();

    $window.fbAsyncInit = function() {
      // Executed when the SDK is loaded
      $rootScope.serviceIds.$promise.then(function(response){
          FB.init({
        appId: response.FacebookAppId,
        channelUrl: 'app/shared/channel.html',
        status: true, /*Set if you want to check the authentication status at the start up of the app */
        cookie: true,
        xfbml: true
          });

          sAuth.watchLoginChange();
      });
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

    setPageTitle();
    $rootScope.originalPath = $location.$$path;

    $rootScope.isLoggedIn = sAuth.isLoggedIn();
    $rootScope.isLoginRequired = next.requireLogin;

    if($location.$$path === '/auth/twitter/callback') {
      sAuth.signinWithTwitter(($location.search()).oauth_token, ($location.search()).oauth_verifier);
    } else if($location.$$path === '/auth/google/callback') {
      sAuth.signinWithGoogle(($location.search()).auth_token);
    } else {
      // Everytime the route in our app changes check authentication status.
      // Get current user only if we are logged in.
      if( $rootScope.isLoggedIn && (undefined === $rootScope.currentUser) ) {
        $rootScope.currentUser = User.get({ id:sAuth.getUserID(), including: "Teams TeamRequests Invitations" });
        console.log('routeChangeStart, current user = ', $rootScope.currentUser);
      }
      // Redirect user to root if he tries to go on welcome page or signin page and he is logged in.
      if( ($location.path() === '/welcome' || $location.path() === '/signin') && $rootScope.isLoggedIn ) {
        console.log('routeChangeStart, redirect to root');
        $location.path('/');
      }
      // Redidrect to welcome if route requires to be logged in and user is not logged in.
      if ( next.requireLogin && ((undefined === $rootScope.currentUser) || !$rootScope.isLoggedIn) ) {
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
        $cookieStore.put('provider', 'google_plus');
        $rootScope.isLoggedIn = true;
        $location.path('/');
      });
    });
  });
  $rootScope.$on('event:google-plus-signin-failure', function (event, authResult) {
    // User has not authorized the G+ App!
  });

  // search function:
  function setPageTitle() {
    if( $location.$$url == '/welcome') {
      $rootScope.title = 'gonawin';
    } else if($location.$$url == '/about') {
      $rootScope.title = 'gonawin - About';
    } else if($location.$$url == '/getting-started') {
      $rootScope.title = 'gonawin - Getting Started';
    }
  }
}]);

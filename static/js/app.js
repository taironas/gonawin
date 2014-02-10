'use strict';

var purpleWingApp = angular.module('purpleWingApp', [
  'ngSanitize',
  'ngRoute',
  'ngResource',
  'ngCookies',
  'directive.g+signin',
  'directive.formValidation',
  'directive.joinButton',
  '$strap.directives',
  
  'homeControllers',
  'sessionControllers',
  'userControllers',
  'teamControllers',
  'tournamentControllers',
  'inviteControllers',
  
  'dataServices'
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

purpleWingApp.factory('loginInterceptor', ['$rootScope', '$location', function($rootScope, $location){
  return {
    request: function(request) {
      console.log('request = ', request);
      return request || $q.when(request);
    }
  };
}]);

purpleWingApp.config(['$routeProvider', '$httpProvider',
  function($routeProvider, $httpProvider) {
    $routeProvider.
      when('/', { templateUrl: 'templates/welcome.html', requireLogin: false }).
      when('/home', { templateUrl: 'templates/home.html', controller: 'HomeCtrl', requireLogin: true }).
      when('/about', { templateUrl: 'templates/about.html', requireLogin: false }).
      when('/contact', { templateUrl: 'templates/contact.html', requireLogin: false }).
      when('/users/', { templateUrl: 'templates/users/index.html', controller: 'UserListCtrl', requireLogin: true }).
      when('/users/show/:id', { templateUrl: 'templates/users/show.html', controller: 'UserShowCtrl', requireLogin: true }).
      when('/teams', { templateUrl: 'templates/teams/index.html', controller: 'TeamListCtrl', requireLogin: true }).
      when('/teams/new', { templateUrl: 'templates/teams/new.html', controller: 'TeamNewCtrl', requireLogin: true }).
      when('/teams/show/:id', { templateUrl: 'templates/teams/show.html', controller: 'TeamShowCtrl', requireLogin: true }).
      when('/teams/edit/:id', { templateUrl: 'templates/teams/edit.html', controller: 'TeamEditCtrl', requireLogin: true }).
      when('/teams/search', { templateUrl: 'templates/teams/index.html', controller: 'TeamSearchCtrl', requireLogin: true}).
      when('/tournaments', { templateUrl: 'templates/tournaments/index.html', controller: 'TournamentListCtrl', requireLogin: true }).
      when('/tournaments/new', { templateUrl: 'templates/tournaments/new.html', controller: 'TournamentNewCtrl', requireLogin: true }).
      when('/tournaments/show/:id', { templateUrl: 'templates/tournaments/show.html', controller: 'TournamentShowCtrl', requireLogin: true }).
      when('/tournaments/edit/:id', { templateUrl: 'templates/tournaments/edit.html', controller: 'TournamentEditCtrl', requireLogin: true }).
      when('/tournaments/edit/:id', { templateUrl: 'templates/tournaments/edit.html', controller: 'TournamentEditCtrl', requireLogin: true }).
      when('/tournaments/search', { templateUrl: 'templates/tournaments/index.html', controller: 'TournamentSearchCtrl', requireLogin: true }).
      when('/tournaments/:id/calendar', { templateUrl: 'templates/tournaments/calendar.html', controller: 'TournamentCalendarCtrl', requireLogin: true }).
      when('/tournaments/:id/firststage', { templateUrl: 'templates/tournaments/firststage.html', controller: 'TournamentFirstStageCtrl', requireLogin: true }).
      when('/tournaments/:id/secondstage', { templateUrl: 'templates/tournaments/secondstage.html', controller: 'TournamentSecondStageCtrl', requireLogin: true }).
      when('/settings/edit-profile', { templateUrl: 'templates/users/edit.html', controller: 'UserEditCtrl', requireLogin: true }).
      when('/settings/networks', { templateUrl: 'templates/settings/networks.html', requireLogin: true }).
      when('/settings/email', { templateUrl: 'templates/settings/email.html', requireLogin: true }).
      when('/invite', { templateUrl: 'templates/invite.html', controller: 'InviteCtrl', requireLogin: true }).
      when('/404', { templateUrl: 'static/templates/404.html' }).
      otherwise( {redirectTo: '/'});
    
    $httpProvider.interceptors.push('notFoundInterceptor');
}]);

purpleWingApp.run(['$rootScope', '$location', function($rootScope, $location) {
  $rootScope.$on("$routeChangeStart", function(event, next, current) {
    console.log('routeChangeStart');
    // Everytime the route in our app changes check authentication status
    $rootScope.$$childHead.initSession().then(function(loggedIn){
      if (next.requireLogin) {
        console.log('requireLogin');
        if(!loggedIn) {
          // if you're logged out send to home page.
          $location.path('/');
        }
      } else {
        console.log('not requireLogin');
        console.log('path = ', $location.path());
        console.log('loggedIn = ', loggedIn);
        if(loggedIn && $location.path() === '/') {
          $location.path('/home');
        }
      }
    });
  });
}]);

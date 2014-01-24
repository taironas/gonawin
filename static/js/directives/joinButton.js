'use strict'

angular.module('directive.joinButton', []).directive('joinbutton', [
  'Team','Tournament', '$compile', '$routeParams', function(Team, Tournament, $compile, $routeParams) {
    var IsTeamJoined = function(teamId, teams) {
      if(!teamId || !teams) {
        return false;
      }
      for (var i=0 ; i<teams.length; i++)
      {
        if(teams[i].Id == teamId)
        {
          return true;
        }
      }
    };
    
    return {
      restrict: 'E',
      scope: {
        target: '=',      // Bind the target to the object given
        resource: '@',    // Store the string associated by resource
        teamid: '='
      },
      replace: true,
      link: function(scope, element, attrs) {
        var createJoinBtn, createLeaveBtn, join_btn, leave_btn;
        join_btn = null;
        leave_btn = null;
        createJoinBtn = function() {
          join_btn = angular.element('<button class="btn btn-primary" ng-disabled="submitting">Join</button>');
          $compile(join_btn)(scope);
          element.append(join_btn);
          return join_btn.bind('click', function(e) {
            scope.submitting = true;
            if(scope.resource == 'tournament')
            {
              if(scope.teamid) {
                Tournament.joinAsTeam({id:$routeParams.id, teamId:scope.teamid}, function(response) {
                  Tournament.get({id:response.Id}, function(response){
                    scope.$parent.tournamentData = response;
                    scope.submitting = false;
                    join_btn.remove();
                  });
                }, function(error) {
                  return scope.submitting = false;
                });
              } else {
                Tournament.join({id:$routeParams.id}, function(response) {
                  Tournament.get({id:response.Id}, function(response){
                    scope.$parent.tournamentData = response;
                    scope.submitting = false;
                    join_btn.remove();
                  });
                }, function(error) {
                  return scope.submitting = false;
                });
              }
            } else if(scope.resource == 'team') {
              Team.join({id:$routeParams.id}, function(response) {
                scope.submitting = false;
                join_btn.remove();
                return createLeaveBtn();
              }, function(error) {
                return scope.submitting = false;
              });
            }

            return scope.$apply();
          });
        };
        createLeaveBtn = function() {
          leave_btn = angular.element('<button class="btn btn-primary" ng-disabled="submitting">Leave</button>');
          $compile(leave_btn)(scope);
          element.append(leave_btn);
          return leave_btn.bind('click', function(e) {
            scope.submitting = true;
            if(scope.resource == 'tournament')
            {
              if(scope.teamid) {
                Tournament.leaveAsTeam({id:$routeParams.id, teamId:scope.teamid}, function(response) {
                  Tournament.get({id:response.Id}, function(response){
                    scope.$parent.tournamentData = response;
                    scope.submitting = false;
                    leave_btn.remove();
                  });
                }, function(error) {
                  return scope.submitting = false;
                });
              } else {
                Tournament.leave({id:$routeParams.id}, function(response) {
                  Tournament.get({id:response.Id}, function(response){
                    scope.$parent.tournamentData = response;
                    scope.submitting = false;
                    leave_btn.remove();
                  });
                }, function(error) {
                  return scope.submitting = false;
                });
              }
            } else if(scope.resource == 'team') {
              Team.leave({id:$routeParams.id}, function(response) {
                scope.submitting = false;
                leave_btn.remove();
                return createJoinBtn();
              }, function(error) {
                return scope.submitting = false;
              });
            }

            return scope.$apply();
          });
        };
        return scope.$watch('target', function(val) {
          if (scope.target) {
            return scope.target.$promise.then(function(result){
              if (result.Joined || IsTeamJoined(scope.teamid, scope.target.Teams)) {
                return createLeaveBtn();
              } else {
                return createJoinBtn();
              }
             });
          }
        });
      }
    };
  }
]);
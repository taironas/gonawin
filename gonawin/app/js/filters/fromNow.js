'use strict'

angular.module('filter.fromNow', []).filter('fromNow', function() {
  return function(dateString) {
    return moment(dateString).fromNow();
  };
});

// filter to have all dates with same time zone.
angular.module('filter.moment', []).filter('moment', function() {
    return function(input, format) {
	// no need to parseint input as input is already with ISO format.
        return moment(input).utc().format(format); 
    };  
});

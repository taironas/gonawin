'use strict'

var themeSelectorService = angular.module('themeSelectorService', []);

themeSelectorService.factory('ThemeSelector', [function() {
   var themes = ["frogideas", "bythepool", "heatwave", "summerwarmth", "duskfalling"];
   return {
     theme: function(userId) {
       var themeId = (userId % 4);

       return themes[themeId];
     },
     numColors: function(userId) {
       return (userId%3)+2;
     }
   };
 }]);

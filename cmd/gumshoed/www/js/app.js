(function() {
	var app = angular.module('gumshoe', []);

  app.controller('ShowController', ['$log', '$http', function($log, $http){
    var showCtrl = this;
    showCtrl.shows = [];

    $http.get("/api/shows").success(function(data){
      showCtrl.shows = data.Shows;
    });

    this.deleteShow = function(index) {
      $http.delete("/api/show/delete/" + showCtrl.shows[index].ID).success(function(data){
        showCtrl.shows.splice(index, 1);
      });
    };

  } ] );

  app.directive("gumshoeTabs", function() {
     return {
       restrict: "E",
       templateUrl: "gumshoe-tabs.html",
       controller: function() {
         this.tab = 1;

         this.isSet = function(checkTab) {
           return this.tab === checkTab;
         };

         this.setTab = function(activeTab) {
           this.tab = activeTab;
         };
       },
       controllerAs: "tab"
     };
   });

	app.directive("gumshoeStatus", function() {
		return {
      restrict: 'E',
      templateUrl: "gumshoe-status.html"
    };
  });
	app.directive("gumshoeShows", function() {
		return {
      restrict: 'E',
      templateUrl: "gumshoe-shows.html"
    };
  });
	app.directive("gumshoeQueue", function() {
		return {
      restrict: 'E',
      templateUrl: "gumshoe-queue.html"
    };
  });
	app.directive("gumshoeSettings", function() {
		return {
      restrict: 'E',
      templateUrl: "gumshoe-settings.html"
    };
  });
})();

(function() {
	var app = angular.module('gumshoe', []);

  app.controller('ShowController', ['$log', '$http', function($log, $http){
    var showCtrl = this;
    showCtrl.shows = [];
    showCtrl.newShow = {};
    var showAddForm = false;

    $http.get("/api/shows").success(function(data){
      showCtrl.shows = data.Shows;
    });

    this.addShow = function() {
      $log.log("addShow");
      switch(this.newShow.episodal) {
        case "true":
          this.newShow.episodal = true;
          break;
        case "false":
          this.newShow.episodal = false;
          break;
      }
      $http.post("/api/show/new", this.newShow).success(function(data){
        showCtrl.shows.push(showCtrl.newShow);
        showCtrl.newShow = {};
        showCtrl.showAddForm = false;
      }).error(function(data, status, headers, config){
        $log.log(data, status, headers, config);
      });
    };

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

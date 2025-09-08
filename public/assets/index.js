var app = (() => {
  // app/assets/index.js
  var app = (() => {
    window.carousel = function(images) {
      return {
        currentIndex: 0,
        images: images || [],
        start() {
          setInterval(() => {
            this.currentIndex = (this.currentIndex + 1) % this.images.length;
          }, 3e3);
        },
        goTo(index) {
          this.currentIndex = index;
        }
      };
    };
  })();
})();

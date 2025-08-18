// app/assets/index.js
var app = (() => {
  window.carousel = function (imageList = []) {
    return {
      currentIndex: 0,
      images:
        Array.isArray(imageList) && imageList.length > 0
          ? imageList
          : [
              // fallback images if none provided
              "https://thewowstyle.com/wp-content/uploads/2015/01/nature-images.jpg",
              "https://c.wallhere.com/photos/95/06/1920x1200_px_landscape_nature_sky-1107228.jpg!d",
              "https://thewowstyle.com/wp-content/uploads/2015/01/nature-wallpaper-27.jpg",
            ],
      start() {
        if (this.images.length > 1) {
          this._timer = setInterval(() => {
            this.currentIndex = (this.currentIndex + 1) % this.images.length;
          }, 3000);
        }
      },
      goTo(index) {
        this.currentIndex = index;
      },
      destroy() {
        this._timer && clearInterval(this._timer);
      },
    };
  };
})();

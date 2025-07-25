var app = (() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

  // app/assets/index.js
  var index_exports = {};
  __export(index_exports, {
    bookingCalculator: () => bookingCalculator
  });
  window.carousel = function() {
    return {
      currentIndex: 0,
      images: [
        "https://thewowstyle.com/wp-content/uploads/2015/01/nature-images.jpg",
        "https://c.wallhere.com/photos/95/06/1920x1200_px_landscape_nature_sky-1107228.jpg!d",
        "https://thewowstyle.com/wp-content/uploads/2015/01/nature-wallpaper-27.jpg"
      ],
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
  function bookingCalculator() {
    const servicePrices = JSON.parse(
      document.getElementById("service-data").textContent
    );
    const baseRate = parseFloat(document.getElementById("base-rate").value);
    const nights = parseInt(document.getElementById("nights-count").value);
    return {
      guests: 1,
      services: [],
      get baseTotal() {
        return baseRate * this.guests * nights;
      },
      get servicesTotal() {
        return this.services.reduce((sum, id) => {
          return sum + (servicePrices[id] || 0);
        }, 0);
      },
      get total() {
        return (this.baseTotal + this.servicesTotal).toFixed(2);
      }
    };
  }
  window.bookingCalculator = bookingCalculator;
  return __toCommonJS(index_exports);
})();

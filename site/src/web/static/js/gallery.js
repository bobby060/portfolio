// Lightbox for album pages. Without JS, each photo link simply opens the full image.
import PhotoSwipeLightbox from '/static/vendor/photoswipe/photoswipe-lightbox.esm.min.js';

const lightbox = new PhotoSwipeLightbox({
  gallery: '#photos',
  children: 'a', // reads href + data-pswp-width/height from each link
  // The main PhotoSwipe code only downloads the first time a photo is clicked
  pswpModule: () => import('/static/vendor/photoswipe/photoswipe.esm.min.js'),
  bgOpacity: 0.5, // the old site's semi-transparent backdrop
});
lightbox.init();

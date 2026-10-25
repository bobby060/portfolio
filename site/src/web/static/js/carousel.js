// Prev/next buttons for scroll-snap carousels. The buttons start hidden in the HTML
// and are revealed here, so without JS nobody sees buttons that do nothing.
for (const carousel of document.querySelectorAll('[data-carousel]')) {
  const track = carousel.querySelector('[data-carousel-track]');

  const step = (dir) => {
    const max = track.scrollWidth - track.clientWidth;
    const atEnd = dir > 0 ? track.scrollLeft >= max - 1 : track.scrollLeft <= 1;
    // Wrap around at either end, like the old carousel
    const left = atEnd ? (dir > 0 ? 0 : max) : track.scrollLeft + dir * track.clientWidth;
    track.scrollTo({ left }); // the track's scroll-smooth class animates this
  };

  for (const [selector, dir] of [['[data-carousel-prev]', -1], ['[data-carousel-next]', 1]]) {
    const button = carousel.querySelector(selector);
    if (!button) continue;
    button.hidden = false;
    button.addEventListener('click', () => step(dir));
  }
}

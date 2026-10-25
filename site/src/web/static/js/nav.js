// Auto-hiding navbar: slides up when scrolling down, comes back when scrolling up.
const nav = document.getElementById('navbar');
const HIDDEN = '-translate-y-full'; // listed in input.css via @source inline()
let lastY = window.scrollY;
let ticking = false;

window.addEventListener('scroll', () => {
  if (ticking) return; // at most one update per animation frame
  ticking = true;
  requestAnimationFrame(() => {
    const y = window.scrollY;
    const delta = y - lastY;
    if (Math.abs(delta) > 5) { // ignore small jitters
      // Only hide once we've scrolled past the bar itself, and never while its menu is open
      const hide = delta > 0 && y > nav.offsetHeight && !nav.querySelector('details[open]');
      nav.classList.toggle(HIDDEN, hide);
      lastY = y;
    }
    ticking = false;
  });
}, { passive: true });

// Keyboard users tabbing into the hidden bar should be able to see where focus is
nav.addEventListener('focusin', () => nav.classList.remove(HIDDEN));

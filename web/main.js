import { createStars } from './render.js';
import { showMoonDay } from './moon.js';
import { initCalendar } from './calendar.js';

document.addEventListener('DOMContentLoaded', () => {
    const starsContainer = document.querySelector('.stars');
    createStars(starsContainer, 150);
    updateTheme
    setInterval(updateTheme, 60000);

    initCalendar();
    
    showMoonDay(new Date(), true);
});

  function updateTheme() {
    const hour = new Date().getHours();
    const body = document.body;

    body.style.transition = 'background-color 0.5s ease, color 0.5s ease';
    body.classList.remove('body--day', 'body-morning', 'body--evening', 'body--night');

    if (hour >= 6 && hour < 12) {
      body.classList.add('body--morning');
    } else if (hour >= 12 && hour < 18) {
      body.classList.add('body--day');
    } else if (hour >= 18 && hour < 23) {
      body.classList.add('body--evening');
    } else {
      body.classList.add('body--night');
    }
  }

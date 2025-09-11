import { createStars } from './render.js';
import { showMoonDay } from './moon.js';
import './calendar.js';

document.addEventListener('DOMContentLoaded', () => {
    const starsContainer = document.querySelector('.stars');
    createStars(starsContainer, 150);

    showMoonDay(new Date(), true);
});
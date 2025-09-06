import { createStars } from './render.js';
import './moon.js';

document.addEventListener('DOMContentLoaded', () => {
    const starsContainer = document.querySelector('.stars');
    createStars(starsContainer, 150);
});

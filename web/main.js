import { createStars } from './render.js';
import { showMoonDay } from './moon.js';
import { initCalendar } from './calendar.js';

document.addEventListener('DOMContentLoaded', () => {
    const starsContainer = document.querySelector('.stars');
    createStars(starsContainer, 150);

    // Инициализируем календарь
    initCalendar();
    
    // Показываем данные для текущей даты
    showMoonDay(new Date(), true);
});
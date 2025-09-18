import { showMoonDay } from './moon.js';

const moonDateInput = document.getElementById('moonDateInput');
const prevMonth = document.getElementById('prevMonth');
const nextMonth = document.getElementById('nextMonth');
const currentMonthYear = document.getElementById('currentMonthYear');
const calendarDays = document.getElementById('calendarDays');

let currentDate = new Date();
let selectedDate = null;

export function initCalendar() {
    renderCalendar();
    
    prevMonth.addEventListener('click', () => {
        currentDate.setMonth(currentDate.getMonth() - 1);
        renderCalendar();});
    
    nextMonth.addEventListener('click', () => {
        currentDate.setMonth(currentDate.getMonth() + 1);
        renderCalendar();});
    
    calendarDays.addEventListener('click', (e) => {
        if (e.target.classList.contains('calendar-day') && !e.target.classList.contains('other-month')) {
            const day = parseInt(e.target.textContent);
            selectDate(new Date(currentDate.getFullYear(), currentDate.getMonth(), day));
            showMoonDay(new Date(currentDate.getFullYear(), currentDate.getMonth(), day), false)}});
    
    moonDateInput.addEventListener('change', (e) => {
        const date = new Date(e.target.value);
        selectDate(date);
        currentDate = new Date(date);
        renderCalendar();});
    
    moonDateInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && moonDateInput.value) {
            const date = new Date(moonDateInput.value);
            selectDate(date);
            showMoonDay(date, false);}});
    
    selectDate(new Date());
}

// Рендер календаря
function renderCalendar() {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    
    currentMonthYear.textContent = `${currentDate.toLocaleString('en', { month: 'long' })} ${year}`;
    currentMonthYear.classList.add('no-wrap');
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    
    let startDay = firstDay.getDay() === 0 ? 6 : firstDay.getDay() - 1;
    
    calendarDays.innerHTML = '';
    
    const prevMonthLastDay = new Date(year, month, 0).getDate();
    for (let i = 0; i < startDay; i++) {
        const day = document.createElement('div');
        day.className = 'calendar-day other-month';
        day.textContent = prevMonthLastDay - startDay + i + 1;
        calendarDays.appendChild(day);}
    
    // Дни текущего месяца
    for (let i = 1; i <= daysInMonth; i++) {
        const day = document.createElement('div');
        day.className = 'calendar-day';
        day.textContent = i;
        
        if (selectedDate && 
            selectedDate.getDate() === i && 
            selectedDate.getMonth() === month && 
            selectedDate.getFullYear() === year) {
            day.classList.add('selected');}
        
        const today = new Date();
        if (i === today.getDate() && month === today.getMonth() && year === today.getFullYear()) {
            day.classList.add('today');}
        
        calendarDays.appendChild(day);}
    
    const totalCells = 42; 
    const remainingCells = totalCells - (startDay + daysInMonth);
    for (let i = 1; i <= remainingCells; i++) {
        const day = document.createElement('div');
        day.className = 'calendar-day other-month';
        day.textContent = i;
        calendarDays.appendChild(day);
    }
}

function selectDate(date) {
    selectedDate = date;

    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');

    moonDateInput.value = `${year}-${month}-${day}`;
    renderCalendar();
}
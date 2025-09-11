import { showMoonDay } from './moon.js';

// Получаем элементы
const chooseDateBtn = document.getElementById('chooseDateBtn'); // если есть
const datePickerContainer = document.getElementById('datePickerContainer');
const moonDateInput = document.getElementById('moonDateInput');

// Убираем старую кнопку выбора даты (если она есть в DOM)
if (chooseDateBtn) {
    chooseDateBtn.style.display = 'none';
}

// Показываем контейнер сразу
datePickerContainer.style.display = 'block';

// Переменные для календаря
let currentDate = new Date();
let selectedDate = null;

// Получаем элементы календаря
const prevMonth = document.getElementById('prevMonth');
const nextMonth = document.getElementById('nextMonth');
const currentMonthYear = document.getElementById('currentMonthYear');
const calendarDays = document.getElementById('calendarDays');
const submitDateBtn = document.getElementById('submitDateBtn');

// Инициализация календаря
export function initCalendar() {
    renderCalendar();
    
    // Навигация по месяцам
    prevMonth.addEventListener('click', () => {
        currentDate.setMonth(currentDate.getMonth() - 1);
        renderCalendar();
    });
    
    nextMonth.addEventListener('click', () => {
        currentDate.setMonth(currentDate.getMonth() + 1);
        renderCalendar();
    });
    
    // Выбор даты из календаря
    calendarDays.addEventListener('click', (e) => {
        if (e.target.classList.contains('calendar-day') && !e.target.classList.contains('other-month')) {
            const day = parseInt(e.target.textContent);
            selectDate(new Date(currentDate.getFullYear(), currentDate.getMonth(), day));
        }
    });
    
    // Синхронизация input date с календарем
    moonDateInput.addEventListener('change', (e) => {
        const date = new Date(e.target.value);
        selectDate(date);
        currentDate = new Date(date);
        renderCalendar();
    });
    
    // Отправка формы
    submitDateBtn.addEventListener('click', () => {
        if (selectedDate) {
            showMoonDay(selectedDate, false);
        }
    });
    
    // Отправка через Enter (оставляем для удобства)
    moonDateInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && moonDateInput.value) {
            const date = new Date(moonDateInput.value);
            selectDate(date);
            showMoonDay(date, false);
        }
    });
}

// Рендер календаря
function renderCalendar() {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    
    // Обновляем заголовок
    currentMonthYear.textContent = `${currentDate.toLocaleString('ru', { month: 'long' })} ${year}`;
    
    // Получаем первый день месяца и количество дней
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    
    // Начинаем с понедельника
    let startDay = firstDay.getDay() === 0 ? 6 : firstDay.getDay() - 1;
    
    // Очищаем дни
    calendarDays.innerHTML = '';
    
    // Дни предыдущего месяца
    const prevMonthLastDay = new Date(year, month, 0).getDate();
    for (let i = 0; i < startDay; i++) {
        const day = document.createElement('div');
        day.className = 'calendar-day other-month';
        day.textContent = prevMonthLastDay - startDay + i + 1;
        calendarDays.appendChild(day);
    }
    
    // Дни текущего месяца
    for (let i = 1; i <= daysInMonth; i++) {
        const day = document.createElement('div');
        day.className = 'calendar-day';
        day.textContent = i;
        
        // Проверяем, выбрана ли эта дата
        if (selectedDate && 
            selectedDate.getDate() === i && 
            selectedDate.getMonth() === month && 
            selectedDate.getFullYear() === year) {
            day.classList.add('selected');
        }
        
        // Проверяем, сегодня ли это
        const today = new Date();
        if (i === today.getDate() && month === today.getMonth() && year === today.getFullYear()) {
            day.style.border = '2px solid #4A5B9C';
        }
        
        calendarDays.appendChild(day);
    }
    
    // Дни следующего месяца
    const totalCells = 42; // 6 строк по 7 дней
    const remainingCells = totalCells - (startDay + daysInMonth);
    for (let i = 1; i <= remainingCells; i++) {
        const day = document.createElement('div');
        day.className = 'calendar-day other-month';
        day.textContent = i;
        calendarDays.appendChild(day);
    }
}

// Выбор даты
function selectDate(date) {
    selectedDate = date;
    moonDateInput.value = date.toISOString().split('T')[0];
    renderCalendar();
}

// Инициализируем календарь
initCalendar();

// Устанавливаем сегодняшнюю дату по умолчанию
selectDate(new Date());
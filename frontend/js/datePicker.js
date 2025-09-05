import { getMoonDataByDate } from './api.js';
import { revealMoon } from './animation.js';

export function initializeDatePicker() {
  const dateBtn = document.getElementById('getMoonDayByDateBtn');
  const datePicker = document.getElementById('moonDatePicker');
  const submitBtn = document.getElementById('submitDateBtn');
  const datePickerContainer = document.getElementById('datePickerContainer');
  const resultDiv = document.getElementById('result');

  datePicker.max = new Date().toISOString().split('T')[0];

  dateBtn.addEventListener('click', () => {
    datePickerContainer.hidden = !datePickerContainer.hidden;
  });

  submitBtn.addEventListener('click', () => {
    if (datePicker.value) {
      const selectedDate = new Date(datePicker.value + 'T12:00:00');
      handleDateSelection(selectedDate, resultDiv);
    }
  });

  datePicker.addEventListener('keypress', e => {
    if (e.key === 'Enter' && datePicker.value) {
      const selectedDate = new Date(datePicker.value + 'T12:00:00');
      handleDateSelection(selectedDate, resultDiv);
    }
  });
}

async function handleDateSelection(date, resultDiv) {
  resultDiv.innerHTML = '<div class="loading-text">Загрузка данных о луне...</div>';
  try {
    const moonData = await getMoonDataByDate(date);
    await revealMoon(moonData.CurrentState.Illumination);
    displayMoonDataForDate(moonData, date, resultDiv);
  } catch (error) {
    resultDiv.innerHTML = `<div class="error">Ошибка: ${error.message}</div>`;
  }
}

function displayMoonDataForDate(moonData, date, resultDiv) {
  const formattedDate = date.toLocaleDateString('ru-RU');
  const moonDay = Math.floor(moonData.CurrentState.MoonDays);
  const illumination = moonData.CurrentState.Illumination.toFixed(1);
  const phase = moonData.CurrentState.Phase;
  const zodiac = moonData.CurrentState.Zodiac;

  resultDiv.innerHTML = `
    <div class="moon-day">Лунный день на ${formattedDate}</div>
    <div class="moon-details">
      <div>Лунный день: <b>${moonDay}</b></div>
      <div>Фаза: ${phase.Emoji} ${phase.Name}</div>
      <div>Знак: ${zodiac.Emoji} ${zodiac.Name}</div>
      <div>Освещение: ${illumination}%</div>
    </div>`;
  resultDiv.className = 'result success';
}

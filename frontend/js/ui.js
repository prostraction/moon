import { getMoonData, getMoonDataByDate } from './api.js';
import { revealMoon } from './render.js';

export function initializeApp() {
  const moonBtn = document.getElementById('getMoonDayBtn');
  moonBtn.addEventListener('click', async () => {
    const moonData = await getMoonData();
    revealMoon(moonData, "Лунный день сегодня");
  });

  const moonByDateBtn = document.getElementById('getMoonDayByDateBtn');
  const datePickerContainer = document.getElementById('datePickerContainer');
  const submitDateBtn2 = document.getElementById('submitDateBtn2');
  const dateInput = document.getElementById('dateInput');

  moonByDateBtn.addEventListener('click', () => {
    datePickerContainer.hidden = !datePickerContainer.hidden;
  });

  submitDateBtn2.addEventListener('click', async () => {
    if (!dateInput.value) return;
    const selectedDate = new Date(dateInput.value);
    const moonData = await getMoonDataByDate(selectedDate);
    revealMoon(moonData, `Лунный день на ${dateInput.value}`);
    datePickerContainer.hidden = true;
  });

  document.addEventListener("keydown", (e) => {
    if (e.key === "Enter" && !datePickerContainer.hidden && dateInput.value) {
      const selectedDate = new Date(dateInput.value);
      getMoonDataByDate(selectedDate).then(data => {
        revealMoon(data, `Лунный день на ${dateInput.value}`);
        datePickerContainer.hidden = true;
      });
    }
  });
}

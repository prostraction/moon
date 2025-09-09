const resultDiv = document.getElementById('result');
const todayBtn = document.getElementById('todayBtn');
const chooseDateBtn = document.getElementById('chooseDateBtn');

const datePickerInline = document.getElementById('datePickerInline');
const moonDateInput = document.getElementById('moonDateInput');


// Ограничение выбора даты
//moonDateInput.max = new Date().toISOString().split('T')[0];

// --- API ---
async function getMoonData(date = new Date()) {
    const params = {
        utc: -Math.round(date.getTimezoneOffset() / 60),
        day: date.getDate(),
        month: date.getMonth() + 1,
        year: date.getFullYear(),
        hour: date.getHours(),
        minute: date.getMinutes(),
        second: date.getSeconds(),
        lang: 'ru'
    };

    const url = new URL(CONFIG.API_URL);
    Object.entries(params).forEach(([key, value]) => url.searchParams.append(key, value.toString()));

    const response = await fetch(url.toString(), {
        method: 'GET',
        headers: { 'Accept': 'application/json' }
    });

    if (!response.ok) throw new Error(`Ошибка API: ${response.status} ${response.statusText}`);
    return await response.json();
}

// --- Отображение ---
async function showMoonDay(date, isCurrent) {
    resultDiv.innerHTML = '<div class="loading-text">Подключаемся к лунному API...</div>';

    try {
        const data = await getMoonData(date);
        let moonDay = Math.floor(data.EndDay.MoonDays);
        let illumination = data.EndDay.Illumination;
        let phase = data.EndDay.Phase;
        let zodiac = data.EndDay.Zodiac;
        
        if (isCurrent) {
            moonDay = Math.floor(data.CurrentState.MoonDays);
            illumination = data.CurrentState.Illumination;
            phase = data.CurrentState.Phase;
            zodiac = data.CurrentState.Zodiac;
        }

        // Удаляем приветственный блок, если он есть
        const initialMessage = resultDiv.querySelector('.initial-message');
        if (initialMessage) initialMessage.remove();

        resultDiv.innerHTML = `
            <div class="moon-day">Лунный день: <span class="highlight">${moonDay}</span></div>
            <div class="moon-details">
                <div class="detail-item"><span class="detail-label">Фаза луны:</span> <span class="detail-value">${phase.Emoji} ${phase.Name}</span></div>
                <div class="detail-item"><span class="detail-label">Освещённость:</span> <span class="detail-value">${illumination}%</span></div>
                <div class="detail-item"><span class="detail-label">Знак зодиака:</span> 
                <span class="detail-value">
                    <img src="icons/${zodiac.Name.toLowerCase()}.svg" alt="${zodiac.Name}" class="zodiac-icon"> 
                    ${zodiac.NameLocalized}
                </span>
            </div>
        `;
    } catch (err) {
        resultDiv.innerHTML = `
            <div class="error-title">Ошибка получения данных</div>
            <div class="error-detail">${err.message || 'Неизвестная ошибка'}</div>
        `;
    }
}


// --- Кнопки ---
todayBtn.addEventListener('click', () => showMoonDay(new Date(), true));

chooseDateBtn.addEventListener('click', () => {
    datePickerInline.classList.toggle('show'); // показываем/скрываем контейнер
    if (datePickerInline.classList.contains('show')) {
        moonDateInput.focus();
    }
});


// Отправка через Enter
moonDateInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter' && moonDateInput.value) {
        const date = new Date(moonDateInput.value);// + 'T12:00:00');
        showMoonDay(date, false);
        datePickerInline.classList.remove('show');
    }
});

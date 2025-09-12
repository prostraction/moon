export { getMoonData, showMoonDay };
import CONFIG from './CONFIG.js';
const resultDiv = document.getElementById('result');

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
        const today = new Date();
        const formattedDate = today.toLocaleDateString("ru-RU");
        if (initialMessage) initialMessage.remove();

        resultDiv.innerHTML = `
            <div class="moon-day">
                <span class="detail-value">${phase.Emoji} Moon day: ${moonDay}
            </div>

            <div class="moon-details">  <br>
                <div class="detail-item"><span class="detail-label">Date:</span><span class="detail-value">${formattedDate}</span></div>
                <div class="detail-item"><span class="detail-label">Moon phaze: </span> <span class="detail-value">${phase.Name}</span></div>
                <div class="detail-item"><span class="detail-label">Illumination:</span> <span class="detail-value">${illumination}%</span></div>
                <div class="detail-item"><span class="detail-label">Zodiac sign:</span> <span class="detail-value">${zodiac.Name}</span></div>
            </div>
        `;
    } catch (err) {
        resultDiv.innerHTML = `
            <div class="error-title">Ошибка получения данных</div>
            <div class="error-detail">${err.message || 'Неизвестная ошибка'}</div>
        `;
    }
}
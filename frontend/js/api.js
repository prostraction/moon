export async function getMoonData() {
    const now = new Date();

    const params = {
        utc: -Math.round(now.getTimezoneOffset() / 60),
        day: now.getDate(),
        month: now.getMonth() + 1,
        year: now.getFullYear(),
        hour: now.getHours(),
        minute: now.getMinutes(),
        second: now.getSeconds()
    };

    const url = new URL('https://moon.qoph.org/v1/moonPhaseDate');
    Object.entries(params).forEach(([key, value]) => {
        url.searchParams.append(key, value.toString());
    });

    console.log('Запрос к API:', url.toString());

    const response = await fetch(url.toString(), {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
        }
    });

    if (!response.ok) {
        throw new Error(`Ошибка API: ${response.status} ${response.statusText}`);
    }
    return await response.json();
}

export async function getMoonDataByDate(date) {
    const params = {
        utc: -Math.round(date.getTimezoneOffset() / 60),
        day: date.getDate(),
        month: date.getMonth() + 1,
        year: date.getFullYear(),
        hour: date.getHours(),
        minute: date.getMinutes(),
        second: date.getSeconds()
    };

    const url = new URL('https://moon.qoph.org/v1/moonPhaseDate');
    Object.entries(params).forEach(([key, value]) => {
        url.searchParams.append(key, value.toString());
    });

    console.log('Запрос к API:', url.toString());

    const response = await fetch(url.toString(), {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
        }
    });

    if (!response.ok) {
        throw new Error(`Ошибка API: ${response.status} ${response.statusText}`);
    }
    return await response.json();
}
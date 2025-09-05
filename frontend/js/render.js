// render.js

// Создание звёзд на фоне
export function createStars() {
  const stars = document.querySelector('.stars');
  if (!stars) return;

  for (let i = 0; i < 100; i++) {
    const star = document.createElement('div');
    star.className = 'star';
    star.style.left = Math.random() * 100 + '%';
    star.style.top = Math.random() * 100 + '%';
    stars.appendChild(star);
  }
}

// Создание облаков на фоне
export function createClouds() {
  const sky = document.querySelector('.sky');
  if (!sky) return;

  const clouds = document.createElement('div');
  clouds.className = 'clouds';
  sky.appendChild(clouds);

  const positions = [
    { left: '10%', top: '20%' },
    { left: '70%', top: '30%' },
    { left: '30%', top: '60%' }
  ];

  positions.forEach(pos => {
    const cloud = document.createElement('div');
    cloud.className = 'cloud';
    cloud.style.left = pos.left;
    cloud.style.top = pos.top;
    cloud.style.width = '80px';
    cloud.style.height = '80px';
    clouds.appendChild(cloud);
  });
}

// Отобразить результат в контейнере
export function revealMoon(moonData, label = "Лунный день") {
  const result = document.getElementById("result");

  if (!result) {
    console.error("Нет контейнера #result для отображения");
    return;
  }

  const current = moonData.CurrentState;

  result.innerHTML = `
    <div class="moon-result">
      <h3>${label}</h3>
      <p><strong>День:</strong> ${current.MoonDays}</p>
      <p><strong>Фаза:</strong> ${current.Phase.Emoji} ${current.Phase.Name}</p>
      <p><strong>Освещённость:</strong> ${current.Illumination}%</p>
      <p><strong>Знак зодиака:</strong> ${current.Zodiac.Emoji} ${current.Zodiac.Name}</p>
    </div>
  `;

  // показать плавно
  result.classList.remove("hidden");
  requestAnimationFrame(() => result.classList.add("show"));
}

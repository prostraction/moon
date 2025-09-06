export function createStars(container, count = 150) {
    const colors = ['#aaf0ff', '#ffffaa', '#ffaaaa', '#ddaaff']; // голубой, жёлтый, розовый, фиолетовый

    for (let i = 0; i < count; i++) {
        const star = document.createElement('div');
        star.className = 'star';

        // Размер 1–3px
        const size = Math.random() * 2 + 1;
        star.style.width = star.style.height = size + 'px';

        // Цвет
        const color = colors[Math.floor(Math.random() * colors.length)];
        star.style.backgroundColor = color;

        // Позиция
        star.style.left = Math.random() * 100 + '%';
        star.style.top = Math.random() * 100 + '%';

        // Начальная прозрачность
        star.style.opacity = Math.random() * 0.8 + 0.2;

        // Анимация
        const duration = Math.random() * 3 + 1.5; // 1.5–4.5s
        const delay = Math.random() * 3;          // 0–3s
        star.style.animationDuration = `${duration}s`;
        star.style.animationDelay = `${delay}s`;

        // Добавляем лёгкое свечение
        star.style.boxShadow = `0 0 ${size * 4}px ${color}`;

        container.appendChild(star);
    }
}

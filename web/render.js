export function createStars(container, count = 300) {
    const colors = ['#ffffffff', '#aaf0ff', '#ffffaa', '#ffaaaa', '#ddaaff'];

    for (let i = 0; i < count; i++) {
        const star = document.createElement('div');
        star.className = 'star';

        // Размер 1–3px
        const size = Math.random() * 3 + 1;
        star.style.width = size + 'px';
        star.style.height = size + 'px';

        // Цвет
        const color = colors[Math.floor(Math.random() * colors.length)];
        star.style.backgroundColor = color;

        // Позиция
        star.style.left = Math.random() * 100 + '%';
        star.style.top = Math.random() * 100 + '%';

        // Начальная прозрачность
        star.style.opacity = Math.random() * 0.8 + 0.2;

        // Анимация
        const duration = Math.random() * 3 + 1.5;
        const delay = Math.random() * 3;
        star.style.animationDuration = duration + 's';
        star.style.animationDelay = delay + 's';

        // Добавляем лёгкое свечение
        star.style.boxShadow = `0 0 ${size * 1}px ${color}`;

        container.appendChild(star);
    }
}


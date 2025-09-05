export function createStars(container, count = 150) {
    const colors = ['#aaf0ff', '#ffffaa', '#ffaaaa', '#ddaaff']; // светло-голубой, желтый, красный, фиолетовый

    for (let i = 0; i < count; i++) {
        const star = document.createElement('div');
        star.className = 'star';

        // Случайный размер
        const size = Math.random() * 2 + 1; // 1-3px
        star.style.width = star.style.height = size + 'px';

        // Случайный цвет
        star.style.backgroundColor = colors[Math.floor(Math.random() * colors.length)];

        // Случайное положение
        star.style.left = Math.random() * 100 + '%';
        star.style.top = Math.random() * 100 + '%';

        // Начальная прозрачность
        star.style.opacity = Math.random() * 0.8 + 0.2; // 0.2 - 1.0

        // Анимация мерцания с разной скоростью
        const duration = Math.random() * 3 + 1.5; // 1.5 - 4.5s
        const delay = Math.random() * 3; // 0-3s
        star.style.animation = `twinkle ${duration}s infinite alternate`;
        star.style.animationDelay = `${delay}s`;

        container.appendChild(star);
    }
}

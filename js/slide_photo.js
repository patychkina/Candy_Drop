const groups = document.querySelectorAll('.group');
const leftArrow = document.querySelector('.arrow.left');
const rightArrow = document.querySelector('.arrow.right');
let activeIndex = 1; // Начинаем с центральной карточки

function updateActiveGroup(newIndex) {
    groups[activeIndex].classList.remove('active');
    activeIndex = newIndex;
    groups[activeIndex].classList.add('active');
}

leftArrow.addEventListener('click', () => {
    const newIndex = (activeIndex - 1 + groups.length) % groups.length; // Переход на предыдущую карточку
    updateActiveGroup(newIndex);
});

rightArrow.addEventListener('click', () => {
    const newIndex = (activeIndex + 1) % groups.length; // Переход на следующую карточку
    updateActiveGroup(newIndex);
});

// Пример на jQuery для обработки кликов по стрелкам
$('.arrow').on('click', function() {
    var direction = $(this).hasClass('left') ? 'left' : 'right';
    var currentGroup = $('.group.active');
    var nextGroup;

    if (direction === 'left') {
        nextGroup = currentGroup.prev('.group');
    } else {
        nextGroup = currentGroup.next('.group');
    }

    if (nextGroup.length > 0) {
        currentGroup.removeClass('active'); // Убираем класс active с текущей группы
        nextGroup.addClass('active'); // Добавляем класс active на следующую группу
    }
});

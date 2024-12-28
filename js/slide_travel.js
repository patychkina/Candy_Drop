const contests = document.querySelectorAll('.contest');
const leftArrow = document.querySelector('.arrow.left');
const rightArrow = document.querySelector('.arrow.right');
let activeIndex = 1; // Начинаем с центральной карточки

function updateActiveContest(newIndex) {
    contests[activeIndex].classList.remove('active');
    activeIndex = newIndex;
    contests[activeIndex].classList.add('active');
}

leftArrow.addEventListener('click', () => {
    const newIndex = (activeIndex - 1 + contests.length) % contests.length; // Переход на предыдущую карточку
    updateActiveContest(newIndex);
});

rightArrow.addEventListener('click', () => {
    const newIndex = (activeIndex + 1) % contests.length; // Переход на следующую карточку
    updateActiveContest(newIndex);
});

// Пример на jQuery для обработки кликов по стрелкам
$('.arrow').on('click', function() {
    var direction = $(this).hasClass('left') ? 'left' : 'right';
    var currentContest = $('.contest.active');
    var nextContest;

    if (direction === 'left') {
        nextContest = currentContest.prev('.contest');
    } else {
        nextContest = currentContest.next('.contest');
    }

    if (nextContest.length > 0) {
        currentContest.removeClass('active'); // Убираем класс active с текущей карточки
        nextContest.addClass('active'); // Добавляем класс active на следующую карточку
    }
});

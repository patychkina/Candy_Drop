document.addEventListener('DOMContentLoaded', function() {
    const authForm = document.getElementById('authForm');
    const errorMessage = document.getElementById('errorMessage');
    const loginInput = document.getElementById('login');
    const passwordInput = document.getElementById('password');

    // Проверка наличия сессии
    fetch('http://localhost:8080/protected', {
        method: 'GET',
        credentials: 'include' // Включите эту строку для отправки куки
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('No session found');
        }
    })
    .then(data => {
        switch (data.privilege) {
            case '0':
                window.location.href = 'admin_page.html';
                break;
            case '1':
                window.location.href = 'teacher_page.html';
                break;
            case '2':
                window.location.href = 'child_page.html';
                break;
            case '3':
                window.location.href = 'parent_page.html';
                break;
            default:
                alert('Unknown privilege level');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });

    authForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const login = loginInput.value;
        const password = passwordInput.value;

        fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ login, password }),
            credentials: 'include' // Включите эту строку для отправки куки
        })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                return response.json().then(errorData => {
                    throw new Error(JSON.stringify(errorData));
                });
            }
        })
        .then(data => {
            errorMessage.style.display = 'none';
            loginInput.style.borderColor = '';
            passwordInput.style.borderColor = '';

            switch (data.privilege) {
                case '0':
                    window.location.href = 'admin_page.html';
                    break;
                case '1':
                    window.location.href = 'teacher_page.html';
                    break;
                case '2':
                    window.location.href = 'child_page.html';
                    break;
                case '3':
                    window.location.href = 'parent_page.html';
                    break;
                default:
                    alert('Unknown privilege level');
            }
        })
        .catch(error => {
            const errorData = JSON.parse(error.message);
            errorMessage.style.display = 'block';
            errorMessage.textContent = error.message;

            if (errorData.message === 'Такого пользователя не существует') {
                loginInput.style.borderColor = 'red';
                passwordInput.style.borderColor = '';
                alert("User");
            } else if (errorData.message === 'Пароль неверный') {
                loginInput.style.borderColor = '';
                passwordInput.style.borderColor = 'red';
                alert("Pass");
            }
        });
    });
});

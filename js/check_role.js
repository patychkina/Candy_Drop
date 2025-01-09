 // Проверка наличия сессии
 fetch('http://localhost:8080/protected', {
    method: 'GET',
    credentials: 'include' // Включаем куки для отправки с запросом
})
.then(response => {
    if (response.ok) {
        return response.json(); 
    } else {
        throw new Error('No session found');
    }
})
.then(data => {
    const { user_id, privilege } = data; 
    
    const authLink = document.querySelector('.auth-icon a');
    switch (privilege) {
        case "0":
            authLink.setAttribute('href', 'admin_page.html');
            break;
        case "1":
            authLink.setAttribute('href', 'teacher_page.html');
            break;
        case "2":
            authLink.setAttribute('href', 'child_page.html');
            break;
        case "3":
            authLink.setAttribute('href', 'parent_page.html');
            break;
        default:
            authLink.setAttribute('href', 'authorization.html');
            break;
    }
})
.catch(error => {
    console.error('Error:', error);
});

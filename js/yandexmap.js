ymaps.ready(init);
function init() {
    var map = new ymaps.Map("map", {
        center: [47.215572, 38.916527],
        zoom: 17,
        controls: ['zoomControl', 'geolocationControl']
    });

    // Добавляем метку
    var placemark = new ymaps.Placemark(
        [47.215572, 38.916527],
        {
            hintContent: "Кузнечная, 6",
            balloonContent: "Ростовская область, г. Таганрог, ул. Кузнечная, 6"
        },
        {
            iconLayout: 'default#image',
            iconImageHref: 'assets/icons/custom-marker.png',
            iconImageSize: [40, 40],
            iconImageOffset: [-20, -40]
        }
    );

    map.geoObjects.add(placemark);

    // Добавляем слой пробок
    var trafficControl = new ymaps.control.TrafficControl({ state: { providerKey: 'traffic#actual', providerStateName: 'provision' } });
    map.controls.add(trafficControl);

}

//////////////////////////////////////////////////////////////
///                ХОТИТЕ ТАКОГО ЖЕ БОТА?                  ///
///             ХОТИТЕ УРОКОВ ПРОГРАММИРОВАНИЯ?            ///
///                https://t.me/Tichomirov2003             ///
//////////////////////////////////////////////////////////////
//tg data
let tg = window.Telegram.WebApp
//Иницилизация
function initialization() {
    let intervalId = setInterval(function() {
        if (tg.initData && tg.initDataUnsafe && tg.initDataUnsafe.user && tg.initDataUnsafe.user.id) {
            getGroupList()
            button_tg()
            clearInterval(intervalId);
            return
        } else {
            console.error('tg.initData is not available yet.');
        }
    }, 100)
}
//Запрос данных (вместе с проверкой)
function getGroupList(){
    var xhr = new XMLHttpRequest();
    var url = "/getGroupList";
    var params ="&id=" + tg.initDataUnsafe.user.id+"&token ="+tg.initData
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.onload = function() {
        if (xhr.status == 200) {
            setGroupList(xhr.responseText)
        } else {
            alert("Не удалось получить список групп.");
        }
    };
    xhr.send(params);
}
//Отправка запроса, получение информации, формирование данных.
function setGroupList(groupList){
    //Обработка json
    var jsonObject = JSON.parse(groupList);
    var groupsArray = jsonObject.groups;
    if (groupsArray.length==0){
        window.location.href = "/noAdmin"
        return
    }
    //Создание динамического контента
    groupsArray.forEach(function(group) {
        var newBlock = document.createElement('a');   
        newBlock.innerHTML = '<a class="data-card" id="id_'+group+'"><h3>'+group+'</h3><button class="button-28" role="button" type="'+group+'" id="deleteButton" name="popup-button">Удалить группу</button><button class="button-28" role="button" type="'+group+'" id="scheduleButton" name="popup-button">Редактировать расписание</button></a>'
        var body = document.getElementById("page-contain")
        body.appendChild(newBlock);
    });
    setButtons()
    
}
function button_tg() {
    tg.MainButton.setText("Создать группу");
    tg.MainButton.textColor = "#FFFFFF";
    tg.MainButton.color = "#753BBD";

    Telegram.WebApp.onEvent('mainButtonClicked', function() {
        window.location.href = "/new-group"
    });
    tg.MainButton.show()
}
//Слушатели и кнопки.
function setButtons() {
    //Добавление слушателей к кнопкам удаления.
    var i = 0
    var buttons = document.querySelectorAll("#deleteButton");
    buttons.forEach(function(button) {
        i += 1
        button.addEventListener("click", function() {
            const type_ = this.getAttribute('type');
            popupDeletingGroup(type_);
        });
    });
    //Если кнопок 0, то вставляется предупреждение.
    if (i === 0) {
        var containerElement = document.getElementById("page-contain");
        containerElement.innerHTML = '<a class="data-card" id="noGroup">' +
            '<h2>У вас нет ни одной группы.</h2>' +
            '<button class="button-28" role="button" name="popup-button" onclick="button_action()">Создать группу</button>' +
            '</a>';
    }
    //Добавление слушателей к кнопкам изменения расписания.
    var scheduleButtons = document.querySelectorAll("#scheduleButton");
    scheduleButtons.forEach(function(button) {
        button.addEventListener("click", function() {
            const type_ = this.getAttribute('type');
            change_schedule(type_);
        });
    });

    Telegram.WebApp.onEvent('backButtonClicked', function() {
        tg.close()
    });
    tg.BackButton.show()
}
//Подтверждение на удаление
function popupDeletingGroup(type_) {
    // Функция для отображения всплывающего окна подтверждения
    var confirmation = confirm("Вы уверены, что хотите удалить группу?");
    if (confirmation) {
        deleteGroup(type_);
    }
}
//Кнопка "нет групп"
function button_action() {
    window.location.href = "/new-group"
}
function change_schedule(name) {
    // Создание формы
    var form = document.createElement('form');
    form.method = 'post';
    form.action = '/changeSchedule';
    // Добавление полей в форму
    var groupNameField = document.createElement('input');
    groupNameField.type = 'hidden';
    groupNameField.name = 'groupName';
    groupNameField.value = (name);
    form.appendChild(groupNameField);


    // Добавление формы в body и ее автоматическая отправка
    document.body.appendChild(form);

    // Подписка на событие onload для обработки ответа после отправки формы
    form.onload = function() {
        alert("Функция завершила работу " + name);
    };

    form.submit();
}
//Запрос на удаление группы с переданным именем
function deleteGroup(name) {
    var xhr = new XMLHttpRequest();
    var url = "/deleteGroup";
    var params ="&id=" + tg.initDataUnsafe.user.id+"&groupname=" + name+"&token ="+tg.initData
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.onload = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var response = xhr.responseText;
            console.log(response);
            // Удаляем элемент с id "name"
            var element = document.getElementById("id_" + name);
            if (element) {
                element.remove();
            }
            alert("Успешно")
        } else {}
    };
    xhr.send(params);
}
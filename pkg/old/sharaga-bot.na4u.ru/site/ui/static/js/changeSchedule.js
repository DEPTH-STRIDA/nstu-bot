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
            //Действие при нажатии кнопки назад (приложение в телеграм)
            button_back_tg()
            getSchedule()
            clearInterval(intervalId);
            return
        } else {
            console.error('tg.initData is not available yet.');
        }
    }, 100)
}
//Действие при нажатии кнопки назад (приложение в телеграм)
function button_back_tg() {
    Telegram.WebApp.onEvent("backButtonClicked", function () {
        window.location.href = "/"
    });
    tg.BackButton.show();
}
//Запрос данных (вместе с проверкой)
function getSchedule(){
    var xhr = new XMLHttpRequest();
    var url = "/getSchedule";
    var groupNameElement = document.getElementById("GroupName");
    var groupNameContent = groupNameElement.textContent
    var params ="&id=" + tg.initDataUnsafe.user.id+"&token ="+tg.initData+"&groupName="+groupNameContent
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.onload = function() {
        if (xhr.status == 200) {
            //Создание контента
            setSchedule(xhr.responseText)
            //Добавление обработчиков
            handlers()
        } else {
            alert("Не удалось получить список групп.");
        }
    };
    xhr.send(params);
}
//Создание динамического контента.
function setSchedule(schedule){
    //Обработка json
    console.log(schedule)
    var jsonObject = JSON.parse(schedule);
    //Создание html
    var containerElement = document.getElementById("form_schedule");
    toInner= '<button class="button-28_input" role="button" id="first_button_1" onclick="editShedule()">Собрать расписание.</button>'

    toInner+='<h2 >Четная неделя</h2>'
    toInner+='<table><tbody>'
    for (var i = 0; i < 7; i++) {
        toInner+=createDay(i,jsonObject,0)
    }
    toInner+='</table></tbody>'
    toInner+='<h2 id="second_label">Нечетная неделя</h2>'

    toInner+='<table><tbody>'
    for (var i = 0; i < 7; i++) {
        toInner+=createDay(i,jsonObject,1)
    }
    toInner+='</table></tbody>'

    containerElement.innerHTML = toInner
   
    
}
var dayOfWeek = ["Понедельник", "Вторник","Среда", "Четверг", "Пятница","Суббота","Воскресенье"];
var dayOfWeekEng = ["monday","tuesday","wednesday","thursday","friday","saturday","sunday"]
function createDay(number,jsonObject,oddOrEven){
    toInner = '<tr><td><h>'+dayOfWeek[number]+'</h></td></tr>'

    if (oddOrEven==0){
        for (var i = 0; i < jsonObject.EvenWeekSchedule[number].length; i++) {
            toInner+='<tr>'
            toInner+='<td><input class="input_class" type="text" onkeydown="disableEnterKey(event)" placeholder="Предмет" id="even_'+dayOfWeekEng[number]+'_'+i+'"'
            toInner+='value="'+jsonObject.EvenWeekSchedule[number][i]+'"'
            toInner+='/></td>'
            toInner+='</tr>'
        }
    }else{
        for (var i = 0; i < jsonObject.OddWeekSchedule[number].length; i++) {
            toInner+='<tr>'
            toInner+='<td><input class="input_class" type="text" onkeydown="disableEnterKey(event)" placeholder="Предмет" id="odd_'+dayOfWeekEng[number]+'_'+i+'"'
            toInner+='value="'+jsonObject.OddWeekSchedule[number][i]+'"'
            toInner+='/></td>'
            toInner+='</tr>'
        }
    }
   

    toInner+='<tr><td>'
    if (oddOrEven==0){
        toInner+='<button class="button_plus" type="button" id="even_'+dayOfWeekEng[number]+'_plus_button">+</button>'
        toInner+='<button class="button_minus" type="button" id="even_'+dayOfWeekEng[number]+'_minus_button">-</button>'
    }else{
        toInner+='<button class="button_plus" type="button" id="odd_'+dayOfWeekEng[number]+'_plus_button">+</button>'
        toInner+='<button class="button_minus" type="button" id="odd_'+dayOfWeekEng[number]+'_minus_button">-</button>'
    }
    toInner+='</tr></td>'
    return toInner
}
function countSubjects() {
    let subjectsCount = {};

    const days = ["even_monday", "even_tuesday", "even_wednesday", "even_thursday", "even_friday", "even_saturday", "even_sunday", "odd_monday", "odd_tuesday", "odd_wednesday", "odd_thursday", "odd_friday", "odd_saturday", "odd_sunday"];

    days.forEach((day) => {
        const inputs = document.querySelectorAll(`input[id^=${day}]`);
        subjectsCount[day + "Index"] = inputs.length;
    });

    return subjectsCount;
}
//Добавление обработчиков на кнопки.
function handlers() {
    // Пример использования
    const subjectsCount = countSubjects();
    console.log(subjectsCount);
    // Счетчики для индексов полей ввода
    let even_mondayIndex = subjectsCount.even_mondayIndex;
    let even_tuesdayIndex = subjectsCount.even_tuesdayIndex;
    let even_wednesdayIndex = subjectsCount.even_wednesdayIndex;
    let even_thursdayIndex = subjectsCount.even_thursdayIndex;
    let even_fridayIndex = subjectsCount.even_fridayIndex;
    let even_saturdayIndex = subjectsCount.even_saturdayIndex;
    let even_sundayIndex = subjectsCount.even_sundayIndex;

    let odd_mondayIndex = subjectsCount.odd_mondayIndex;
    let odd_tuesdayIndex = subjectsCount.odd_tuesdayIndex;
    let odd_wednesdayIndex = subjectsCount.odd_wednesdayIndex;
    let odd_thursdayIndex = subjectsCount.odd_thursdayIndex;
    let odd_fridayIndex = subjectsCount.odd_fridayIndex;
    let odd_saturdayIndex = subjectsCount.odd_saturdayIndex;
    let odd_sundayIndex = subjectsCount.odd_sundayIndex;
    // Функция для создания нового поля ввода
    function createInput(day, index) {
        const newRow = $("<tr>");
        newRow.append(`<td><input class="input_class" type="text" onkeydown="disableEnterKey(event)" placeholder="Предмет" id="${day}_${index}" /></td>`);
        $("#" + day + "_plus_button")
            .closest("tr")
            .before(newRow);
    }

    // Обработчики событий для понедельника
    $("#even_monday_plus_button").on("click", function () {
        createInput("even_monday", even_mondayIndex);
        even_mondayIndex++;
    });
    $("#even_monday_minus_button").on("click", function () {
        if (even_mondayIndex > 1) {
            even_mondayIndex--;
            $("#even_monday_" + even_mondayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для вторника
    $("#even_tuesday_plus_button").on("click", function () {
        createInput("even_tuesday", even_tuesdayIndex);
        even_tuesdayIndex++;
    });
    $("#even_tuesday_minus_button").on("click", function () {
        if (even_tuesdayIndex > 1) {
            even_tuesdayIndex--;
            $("#even_tuesday_" + even_tuesdayIndex)
                .closest("tr")
                .remove();
        }
    });
    // Обработчики событий для среды
    $("#even_wednesday_plus_button").on("click", function () {
        createInput("even_wednesday", even_wednesdayIndex);
        even_wednesdayIndex++;
    });
    $("#even_wednesday_minus_button").on("click", function () {
        if (even_wednesdayIndex > 1) {
            even_wednesdayIndex--;
            $("#even_wednesday_" + even_wednesdayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для четверга
    $("#even_thursday_plus_button").on("click", function () {
        createInput("even_thursday", even_thursdayIndex);
        even_thursdayIndex++;
    });
    $("#even_thursday_minus_button").on("click", function () {
        if (even_thursdayIndex > 1) {
            even_thursdayIndex--;
            $("#even_thursday_" + even_thursdayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для пятницы
    $("#even_friday_plus_button").on("click", function () {
        createInput("even_friday", even_fridayIndex);
        even_fridayIndex++;
    });
    $("#even_friday_minus_button").on("click", function () {
        if (even_fridayIndex > 1) {
            even_fridayIndex--;
            $("#even_friday_" + even_fridayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для субботы
    $("#even_saturday_plus_button").on("click", function () {
        createInput("even_saturday", even_saturdayIndex);
        even_saturdayIndex++;
    });
    $("#even_saturday_minus_button").on("click", function () {
        if (even_saturdayIndex > 1) {
            even_saturdayIndex--;
            $("#even_saturday_" + even_saturdayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для воскресенья
    $("#even_sunday_plus_button").on("click", function () {
        createInput("even_sunday", even_sundayIndex);
        even_sundayIndex++;
    });
    $("#even_sunday_minus_button").on("click", function () {
        if (even_sundayIndex > 1) {
            even_sundayIndex--;
            $("#even_sunday_" + even_sundayIndex)
                .closest("tr")
                .remove();
        }
    });
    ///////////////////////
    // Обработчики событий для понедельника
    $("#odd_monday_plus_button").on("click", function () {
        createInput("odd_monday", odd_mondayIndex);
        odd_mondayIndex++;
    });
    $("#odd_monday_minus_button").on("click", function () {
        if (odd_mondayIndex > 1) {
            odd_mondayIndex--;
            $("#odd_monday_" + odd_mondayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для вторника
    $("#odd_tuesday_plus_button").on("click", function () {
        createInput("odd_tuesday", odd_tuesdayIndex);
        odd_tuesdayIndex++;
    });
    $("#odd_tuesday_minus_button").on("click", function () {
        if (odd_tuesdayIndex > 1) {
            odd_tuesdayIndex--;
            $("#odd_tuesday_" + odd_tuesdayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для среды
    $("#odd_wednesday_plus_button").on("click", function () {
        createInput("odd_wednesday", odd_wednesdayIndex);
        odd_wednesdayIndex++;
    });
    $("#odd_wednesday_minus_button").on("click", function () {
        if (odd_wednesdayIndex > 1) {
            odd_wednesdayIndex--;
            $("#odd_wednesday_" + odd_wednesdayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для четверга
    $("#odd_thursday_plus_button").on("click", function () {
        createInput("odd_thursday", odd_thursdayIndex);
        odd_thursdayIndex++;
    });
    $("#odd_thursday_minus_button").on("click", function () {
        if (odd_thursdayIndex > 1) {
            odd_thursdayIndex--;
            $("#odd_thursday_" + odd_thursdayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для пятницы
    $("#odd_friday_plus_button").on("click", function () {
        createInput("odd_friday", odd_fridayIndex);
        odd_fridayIndex++;
    });
    $("#odd_friday_minus_button").on("click", function () {
        if (odd_fridayIndex > 1) {
            odd_fridayIndex--;
            $("#odd_friday_" + odd_fridayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для субботы
    $("#odd_saturday_plus_button").on("click", function () {
        createInput("odd_saturday", odd_saturdayIndex);
        odd_saturdayIndex++;
    });
    $("#odd_saturday_minus_button").on("click", function () {
        if (odd_saturdayIndex > 1) {
            odd_saturdayIndex--;
            $("#odd_saturday_" + odd_saturdayIndex)
                .closest("tr")
                .remove();
        }
    });

    // Обработчики событий для воскресенья
    $("#odd_sunday_plus_button").on("click", function () {
        createInput("odd_sunday", odd_sundayIndex);
        odd_sundayIndex++;
    });
    $("#odd_sunday_minus_button").on("click", function () {
        if (odd_sundayIndex > 1) {
            odd_sundayIndex--;
            $("#odd_sunday_" + odd_sundayIndex)
                .closest("tr")
                .remove();
        }
    });
}
/////////////////////////////////////////////////
///             new_schedule_group            ///
/////////////////////////////////////////////////
function collectData() {
  var allData = [];
  for (var weekType of ["even", "odd"]) {
      for (var day = 0; day < 7; day++) {
          var dataArray = [];
          for (var index = 0; ; index++) {
              var inputId = weekType + "_" + getDayName(day) + "_" + index;
              var inputElement = document.getElementById(inputId);

              if (!inputElement) {
                  break;
              }

              var inputValue = inputElement.value;
              dataArray.push(inputValue);
          }
          allData.push(dataArray);
      }
  }

  return allData;
}
function getDayName(dayNumber) {
  var daysOfWeek = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"];
  return daysOfWeek[dayNumber];
}
function splitArraysToJsonString(inputArray) {
  if (inputArray.length !== 14) {
      console.error("Неверный размер входного массива. Ожидалось 14 подмассивов.");
      return null;
  }

  var evenWeek = inputArray.slice(0, 7);
  var oddWeek = inputArray.slice(7);

  var result = {
      Even_week_schedule: evenWeek,
      Odd_week_schedule: oddWeek,
  };

  return JSON.stringify(result);
}
function collect() {
  var allWeekData = collectData();
  console.log("Все данные:", allWeekData);
  var jsonString = splitArraysToJsonString(allWeekData);
  return jsonString;
}
function editShedule() {
  //Имя группы
  var groupNameElement = document.getElementById("GroupName");
  var groupNameContent = groupNameElement.textContent

  var xhr = new XMLHttpRequest();
  var url = "/collect_schedule";

  var schedule_json = collect();
  if (schedule_json == "error") {
      alert("Неизвестная ошибка.");
      return;
  }
  var params = "schedule=" + encodeURIComponent(schedule_json) + "&id=" + tg.initDataUnsafe.user.id + "&groupName=" + groupNameContent+ "&token="+tg.initData
  xhr.open("POST", url, true);
  xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhr.onload = function () {
      if (xhr.status == 200) {
          alert(xhr.responseText);
      } else {
          alert(xhr.responseText);
      }
  };

  xhr.send(params);
}




						
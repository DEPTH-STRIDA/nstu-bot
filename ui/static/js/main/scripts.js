let tg, initData;
// Пользователь является членом группы.
let consist_of;

// main
document.addEventListener("DOMContentLoaded", function () {
  // // Иницилизация телеграмм
  // const result = initializeTg(5, 2000); // 5 попыток с интервалом в 2 секунды
  // if (result) {
  //   tg = result.tg;
  //   initData = result.initData;
  //   // Дальнейшая логика с использованием tg и initData
  // } else {
  //   showAlert("Пожалуйста, запустите форму через телеграмм", 15, true);
  //   return;
  // }

  // Загрузка всех групп
  SetGroups();

  setupGlobalClickListener();
});

function setupGlobalClickListener() {
  document.addEventListener("click", function (event) {
    const clickedElement = event.target;

    // Функция для получения ближайшего родителя с указанным data-type
    function getClosestElementWithType(element) {
      return element.closest("[data-type]");
    }

    const targetElement = getClosestElementWithType(clickedElement);

    if (!targetElement) return; // Если не нашли элемент с data-type, прекращаем выполнение

    const elementType = targetElement.getAttribute("data-type");
    const elementId = targetElement.getAttribute("data-id");

    switch (elementType) {
      case "group-card":
        handleGroupCardClick(targetElement, clickedElement);
        break;
      case "user-profile":
        // handleUserProfileClick(targetElement);
        break;
      case "notification":
        // handleNotificationClick(targetElement);
        break;
      case "chat-message":
        // handleChatMessageClick(targetElement);
        break;
      default:
        console.log("Неизвестный тип элемента:", elementType);
    }
  });
}

function SetGroups() {
  // Вставка html блоков в DOM
  function SetHtmlCards(groups) {
    const targetElement = document.getElementById("cards-placeholder");

    // Очищаем содержимое целевого элемента перед добавлением новых карточек
    targetElement.innerHTML = "";

    for (const group of groups) {
      console.log(
        `ID: ${group.id}, Name: ${group.name}, Title: ${group.title}`
      );

      const htmlString = `
        <div class="group-card" group-id="${group.id}">
          <img class="member-img" src="/static/img/member.svg" alt="" />
          <h2 class="group-name">${group.name}</h2>
          <div class="divider"></div>
          <h2 class="desc">${group.title}</h2>
          <button class="join-group" id="enter-button" action_type="enter" group_button_id="${group.id}">Вступить</button>
        </div>
      `;

      targetElement.insertAdjacentHTML("beforeend", htmlString);
    }
  }

  sendPostRequestAsync("/get/groups", JSON.stringify({ initData: initData }))
    .then((result) => {
      if (result.success) {
        // Успешный запрос
        const responseBody = result.response.responseText;
        console.log("Тело ответа:", responseBody);

        // Если ответ в формате JSON, вы можете распарсить его:
        try {
          const jsonResponse = JSON.parse(responseBody);
          console.log("Распарсенный JSON:", jsonResponse);

          SetHtmlCards(jsonResponse["groups"]);
          consist_of = jsonResponse["onsists-of"];
        } catch (e) {
          console.log("Ответ не является валидным JSON");
        }
      } else {
        // Запрос завершился с ошибкой
        console.log(
          "Ошибка:",
          result.response.status,
          result.response.statusText
        );
      }
    })
    .catch((error) => {
      console.error("Произошла ошибка:", error);
    });
}

// Обработчики на нажатие поделиться
document.addEventListener("DOMContentLoaded", function () {
  // Получаем ссылки на элементы DOM
  const shareBtn = document.getElementById("share-btn");
  const shareVariants = document.getElementById("share-variants");
  let isTooltipVisible = false;
  let isFirstClickAfterShow = false;

  // Функция для переключения видимости подсказки
  function toggleTooltip() {
    if (isTooltipVisible) {
      shareVariants.style.display = "none";
      isTooltipVisible = false;
      isFirstClickAfterShow = false;
    } else {
      shareVariants.style.display = "block";
      isTooltipVisible = true;
      isFirstClickAfterShow = true;
    }
  }

  // Обработчик клика на кнопку
  shareBtn.addEventListener("click", (event) => {
    event.stopPropagation(); // Предотвращаем всплытие события
    toggleTooltip();
  });

  // Обработчик клика на документ
  document.addEventListener("click", (event) => {
    if (isTooltipVisible) {
      if (isFirstClickAfterShow) {
        // Игнорируем первый клик после показа подсказки
        isFirstClickAfterShow = false;
      } else if (!shareVariants.contains(event.target)) {
        // Скрываем подсказку при клике вне её
        toggleTooltip();
      }
    }
  });

  // Предотвращаем скрытие подсказки при клике на неё
  shareVariants.addEventListener("click", (event) => {
    event.stopPropagation();
  });
});

// Анимация открытия контейнера с полем ввода времени
document.addEventListener("DOMContentLoaded", function () {
  const checkbox = document.getElementById("is-emergency");
  const dateInputContainer = document.getElementById("dateInputContainer");

  checkbox.addEventListener("change", function () {
    if (this.checked) {
      dateInputContainer.style.maxHeight =
        dateInputContainer.scrollHeight + "px";
      dateInputContainer.style.opacity = "1";
    } else {
      dateInputContainer.style.maxHeight = "0";
      dateInputContainer.style.opacity = "0";
    }
  });

  // Инициализация начального состояния
  if (checkbox.checked) {
    dateInputContainer.style.maxHeight = dateInputContainer.scrollHeight + "px";
    dateInputContainer.style.opacity = "1";
  } else {
    dateInputContainer.style.maxHeight = "0";
    dateInputContainer.style.opacity = "0";
  }
});

// Переключение выбора недели
document.addEventListener("DOMContentLoaded", function () {
  const container = document.querySelector(".schedule .horizontal");
  const buttons = container.querySelectorAll(".week-button");

  // Создаем и добавляем элемент индикатора выбора
  const indicator = document.createElement("div");
  indicator.classList.add("selection-indicator");
  container.appendChild(indicator);

  function updateSelection(button) {
    const buttonRect = button.getBoundingClientRect();
    const containerRect = container.getBoundingClientRect();

    indicator.style.width = `${buttonRect.width}px`;
    indicator.style.left = `${buttonRect.left - containerRect.left}px`;
  }

  // Инициализируем позицию индикатора
  updateSelection(buttons[0]);

  buttons.forEach((button) => {
    button.addEventListener("click", function () {
      buttons.forEach((btn) => btn.classList.remove("selected"));
      this.classList.add("selected");
      updateSelection(this);
    });
  });
});

// Открытие закрытие секции
document.addEventListener("DOMContentLoaded", function () {
  const weekButtons = document.querySelectorAll(".horizontal button");
  const dayHeaders = document.querySelectorAll(".day-header");

  // Переключение между неделями (оставляем без изменений текста)
  weekButtons.forEach((button) => {
    button.addEventListener("click", function () {
      weekButtons.forEach((btn) => btn.classList.remove("selected"));
      this.classList.add("selected");
    });
  });

  // Сворачивание/разворачивание дней
  dayHeaders.forEach((header) => {
    header.addEventListener("click", function () {
      const weekDay = this.closest(".week-day");
      weekDay.classList.toggle("collapsed");
    });
  });
});

document.addEventListener("DOMContentLoaded", function () {
  const weekSelector = document.getElementById("week-selector");
  const evenWeekButton = document.getElementById("even-week");
  const oddWeekButton = document.getElementById("odd-week");
  const allScheduleButton = document.getElementById("all-schedule");
  const evenWeekSchedule = document.querySelector(".even-week-schedule");
  const oddWeekSchedule = document.querySelector(".odd-week-schedule");
  const titleWeeks = document.querySelectorAll(".title-week");
  const weekAlternationCheckbox = document.getElementById("is-emergency");

  function updateScheduleView() {
    const isAlternating = weekAlternationCheckbox.checked;

    evenWeekButton.style.display = isAlternating ? "block" : "none";
    oddWeekButton.style.display = isAlternating ? "block" : "none";
    allScheduleButton.style.display = isAlternating ? "none" : "block";

    if (!isAlternating) {
      evenWeekSchedule.style.display = "block";
      oddWeekSchedule.style.display = "none";
      titleWeeks.forEach((title) => (title.style.display = "none"));
      allScheduleButton.classList.add("selected");
      evenWeekButton.classList.remove("selected");
      oddWeekButton.classList.remove("selected");
    } else {
      titleWeeks.forEach((title) => (title.style.display = "block"));
      // Проверяем, какая неделя была выбрана ранее
      if (
        !evenWeekButton.classList.contains("selected") &&
        !oddWeekButton.classList.contains("selected")
      ) {
        // Если ни одна неделя не выбрана, выбираем четную по умолчанию
        evenWeekButton.classList.add("selected");
      }
      if (evenWeekButton.classList.contains("selected")) {
        evenWeekSchedule.style.display = "block";
        oddWeekSchedule.style.display = "none";
        oddWeekButton.classList.remove("selected");
      } else if (oddWeekButton.classList.contains("selected")) {
        evenWeekSchedule.style.display = "none";
        oddWeekSchedule.style.display = "block";
        evenWeekButton.classList.remove("selected");
      }
    }

    allScheduleButton.classList.remove("selected");
    updateSelectionIndicator();
  }

  // Обновляем обработчик изменения состояния чекбокса
  weekAlternationCheckbox.addEventListener("change", function () {
    updateScheduleView();
    // Если включено чередование и ни одна неделя не выбрана, выбираем четную
    if (
      this.checked &&
      !evenWeekButton.classList.contains("selected") &&
      !oddWeekButton.classList.contains("selected")
    ) {
      evenWeekButton.classList.add("selected");
      updateScheduleView(); // Вызываем еще раз для обновления отображения
    }
  });

  // Обновляем обработчик кликов по кнопкам выбора недели
  weekSelector.addEventListener("click", function (event) {
    if (event.target.classList.contains("week-button")) {
      weekSelector
        .querySelectorAll(".week-button")
        .forEach((btn) => btn.classList.remove("selected"));
      event.target.classList.add("selected");
      updateScheduleView();
    }
  });

  function updateSelectionIndicator() {
    const indicator = weekSelector.querySelector(".selection-indicator");
    const selectedButton = weekSelector.querySelector(".selected");

    if (selectedButton) {
      const buttonRect = selectedButton.getBoundingClientRect();
      const containerRect = weekSelector.getBoundingClientRect();

      indicator.style.width = `${buttonRect.width}px`;
      indicator.style.left = `${buttonRect.left - containerRect.left}px`;
    }
  }

  weekAlternationCheckbox.addEventListener("change", updateScheduleView);

  weekSelector.addEventListener("click", function (event) {
    if (event.target.classList.contains("week-button")) {
      weekSelector
        .querySelectorAll(".week-button")
        .forEach((btn) => btn.classList.remove("selected"));
      event.target.classList.add("selected");
      updateScheduleView();
    }
  });

  // Инициализация
  updateScheduleView();
});

document.addEventListener("DOMContentLoaded", function () {
  function addSubject(weekType, dayIndex) {
    const weekDay = document.querySelector(`#${weekType}-${dayIndex}`);
    if (!weekDay) {
      console.error(`Week day not found for ${weekType}-${dayIndex}`);
      return;
    }

    const container = weekDay.querySelector(".day-content");
    if (!container) {
      console.error(`Container not found for ${weekType}-${dayIndex}`);
      return;
    }

    let table = container.querySelector("table");
    if (!table) {
      table = document.createElement("table");
      container.insertBefore(table, container.querySelector(".add-subject"));
    }

    const newRow = table.insertRow(-1);
    const cellCount = 5;

    for (let i = 0; i < cellCount; i++) {
      const cell = newRow.insertCell(i);
      if (i === 0) {
        cell.textContent = table.rows.length;
        cell.className = "number";
      } else if (i === cellCount - 1) {
        cell.className = "mini-menu";
        cell.innerHTML = `
          <img src="/static/img/add-menu.svg" alt="" class="open-mini-menu" />
          <div class="opened-menu">
            <img class="arrow" src="/static/img/mini-menu/arrow-0.svg" alt="" />
            <img class="delete" src="/static/img/mini-menu/delete.svg" alt="" />
            <img class="arrow" src="/static/img/mini-menu/arrow-1.svg" alt="" />
          </div>
        `;
      } else {
        const textarea = document.createElement("textarea");
        if (i === 1) {
          cell.className = "time";
          textarea.className = "time-input";
        } else if (i === 2) {
          cell.className = "subject";
          textarea.className = "subject-input";
        } else {
          cell.className = "audience";
          textarea.className = "audience-input";
        }
        cell.appendChild(textarea);
      }
    }

    addMiniMenuHandlers(newRow.querySelector(".mini-menu"));
    recalculateRowNumbers(table);
  }

  function updateRowNumbers(table) {
    const rows = Array.from(table.rows);
    rows.forEach((row, index) => {
      const numberCell = row.cells[0];
      if (numberCell) {
        numberCell.textContent = index + 1;
      }
    });
  }

  function moveRow(row, direction) {
    const table = row.closest("table");
    const index = Array.from(table.rows).indexOf(row);

    if (direction === "up" && index > 0) {
      table.rows[index - 1].before(row);
    } else if (direction === "down" && index < table.rows.length - 1) {
      table.rows[index + 1].after(row);
    }

    recalculateRowNumbers(table);
  }

  function recalculateRowNumbers(table) {
    const rows = Array.from(table.rows);
    rows.forEach((row, index) => {
      const numberCell = row.cells[0];
      if (numberCell && numberCell.classList.contains("number")) {
        numberCell.textContent = index + 1;
      }
    });
  }

  function addMiniMenuHandlers(miniMenu) {
    if (!miniMenu) {
      console.error("Mini menu not found");
      return;
    }

    const openButton = miniMenu.querySelector(".open-mini-menu");
    const openedMenu = miniMenu.querySelector(".opened-menu");
    const upArrow = openedMenu.querySelector(".arrow:first-child");
    const downArrow = openedMenu.querySelector(".arrow:last-child");
    const deleteButton = openedMenu.querySelector(".delete");

    openButton.addEventListener("click", (e) => {
      e.stopPropagation();
      openedMenu.style.display =
        openedMenu.style.display === "flex" ? "none" : "flex";
    });

    upArrow.addEventListener("click", (e) => {
      e.stopPropagation();
      moveRow(miniMenu.closest("tr"), "up");
      openedMenu.style.display = "none";
    });

    downArrow.addEventListener("click", (e) => {
      e.stopPropagation();
      moveRow(miniMenu.closest("tr"), "down");
      openedMenu.style.display = "none";
    });

    deleteButton.addEventListener("click", (e) => {
      e.stopPropagation();
      deleteRow(miniMenu.closest("tr"));
      openedMenu.style.display = "none";
    });

    document.addEventListener("click", () => {
      openedMenu.style.display = "none";
    });

    openedMenu.addEventListener("click", (e) => {
      e.stopPropagation();
    });
  }

  function deleteRow(row) {
    const table = row.closest("table");
    row.remove();
    recalculateRowNumbers(table);
  }

  function updateRowNumbers(table) {
    const rows = table.rows;
    for (let i = 1; i < rows.length; i++) {
      const numberCell = rows[i].cells[0];
      if (numberCell) {
        numberCell.textContent = i;
      }
    }
  }

  function updateRowNumbers(table) {
    Array.from(table.rows).forEach((row, index) => {
      if (index > 0) {
        row.cells[0].textContent = index;
      }
    });
  }

  // Добавляем обработчики для кнопок "добавить запись"
  document.querySelectorAll(".add-subject").forEach((button) => {
    button.addEventListener("click", () => {
      const weekDay = button.closest(".week-day");
      if (weekDay) {
        const weekType =
          weekDay.id.split("-")[0] + "-" + weekDay.id.split("-")[1];
        const dayIndex = weekDay.id.split("-")[2];
        addSubject(weekType, dayIndex);
      } else {
        console.error("Week day not found for button:", button);
      }
    });
  });

  // Добавляем обработчики для существующих мини-меню
  document.querySelectorAll(".mini-menu").forEach(addMiniMenuHandlers);

  // Обработка переключения между неделями
  const weekSelector = document.getElementById("week-selector");
  const evenWeekButton = document.getElementById("even-week");
  const oddWeekButton = document.getElementById("odd-week");
  const allScheduleButton = document.getElementById("all-schedule");
  const evenWeekSchedule = document.querySelector(".even-week-schedule");
  const oddWeekSchedule = document.querySelector(".odd-week-schedule");
  const titleWeeks = document.querySelectorAll(".title-week");
  const weekAlternationCheckbox = document.getElementById("is-emergency");

  function updateScheduleView() {
    const isAlternating = weekAlternationCheckbox.checked;

    evenWeekButton.style.display = isAlternating ? "block" : "none";
    oddWeekButton.style.display = isAlternating ? "block" : "none";
    allScheduleButton.style.display = isAlternating ? "none" : "block";

    if (!isAlternating) {
      evenWeekSchedule.style.display = "block";
      oddWeekSchedule.style.display = "none";
      titleWeeks.forEach((title) => (title.style.display = "none"));
      allScheduleButton.classList.add("selected");
      evenWeekButton.classList.remove("selected");
      oddWeekButton.classList.remove("selected");
    } else {
      titleWeeks.forEach((title) => (title.style.display = "block"));
      if (evenWeekButton.classList.contains("selected")) {
        evenWeekSchedule.style.display = "block";
        oddWeekSchedule.style.display = "none";
      } else if (oddWeekButton.classList.contains("selected")) {
        evenWeekSchedule.style.display = "none";
        oddWeekSchedule.style.display = "block";
      }
    }

    updateSelectionIndicator();
  }

  function updateSelectionIndicator() {
    const indicator = weekSelector.querySelector(".selection-indicator");
    const selectedButton = weekSelector.querySelector(".selected");

    if (selectedButton) {
      const buttonRect = selectedButton.getBoundingClientRect();
      const containerRect = weekSelector.getBoundingClientRect();

      indicator.style.width = `${buttonRect.width}px`;
      indicator.style.left = `${buttonRect.left - containerRect.left}px`;
    }
  }

  weekAlternationCheckbox.addEventListener("change", updateScheduleView);

  weekSelector.addEventListener("click", function (event) {
    if (event.target.classList.contains("week-button")) {
      weekSelector
        .querySelectorAll(".week-button")
        .forEach((btn) => btn.classList.remove("selected"));
      event.target.classList.add("selected");
      updateScheduleView();
    }
  });

  // Инициализация
  updateScheduleView();
});

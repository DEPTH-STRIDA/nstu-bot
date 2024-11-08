let tg;
let initData = {
  user: {
    id: 1, // Временное значение, замените на реальное ID пользователя
    // Добавьте другие необходимые данные пользователя
  },
  // Добавьте другие необходимые поля
};
let consist_of;
let currentGroupId = null;
let isRequestInProgress = false;
let currentUserGroupId = null;

document.addEventListener("DOMContentLoaded", function () {
  hideScheduleContainer();
  SetGroups();
  SetMyGroups();
  setupGlobalClickListener();
  initializeScheduleUI();
  initializeShareButton();
  initializeCreateSchedule();
});

function setupGlobalClickListener() {
  document.addEventListener("click", function (event) {
    const clickedElement = event.target;

    const joinGroupButton = clickedElement.closest(".join-group");
    if (joinGroupButton) {
      handleGroupButton(joinGroupButton);
      return;
    }

    const groupCard = clickedElement.closest(".group-card");
    if (groupCard && !clickedElement.closest(".join-group")) {
      handleGroupCardClick(groupCard);
      return;
    }

    const exitScheduleButton = clickedElement.closest(".save-schedule");
    if (exitScheduleButton) {
      hideScheduleContainer();
      return;
    }

    const addSubjectButton = clickedElement.closest(".add-subject");
    if (addSubjectButton) {
      const weekDay = addSubjectButton.closest(".week-day");
      if (weekDay) {
        const [weekType, , dayIndex] = weekDay.id.split("-");
        addSubject(`${weekType}-week`, dayIndex);
      }
      return;
    }

    const miniMenu = clickedElement.closest(".mini-menu");
    if (miniMenu) {
      handleMiniMenuClick(clickedElement, miniMenu);
      return;
    }
  });
}

function handleGroupCardClick(groupCard) {
  const groupId = groupCard.getAttribute("group-id");
  currentGroupId = groupId;

  sendPostRequestAsync(
    "/get/group-schedule",
    JSON.stringify({
      initData: initData,
      "group-id": groupId,
    })
  )
    .then((result) => {
      if (result.success) {
        try {
          const scheduleData = JSON.parse(result.response.responseText);
          console.log("С сервера получены данные: ", scheduleData);
          displaySchedule(scheduleData);
          showScheduleContainer();
        } catch (error) {
          console.error("Ошибка при разборе ответа:", error);
          alert("Не удалось загрузить расписание группы");
        }
      } else {
        console.error(
          "Ошибка при получении расписания:",
          result.response.statusText
        );
        alert("Не удалось загрузить расписание группы");
      }
    })
    .catch((error) => {
      console.error("Ошибка при отправке запроса:", error);
      alert("Произошла ошибка при загрузке расписания");
    });
}

function displaySchedule(data) {
  if (!data || typeof data !== "object") {
    console.error("Получены некорректные данные расписания:", data);
    alert("Не удалось загрузить расписание группы");
    return;
  }

  console.log("Данные расписания:", data);

  document.querySelector('.schedule-container input[placeholder="Имя"]').value =
    data.Name || "";
  document.querySelector(
    '.schedule-container input[placeholder="Подпись"]'
  ).value = data.Title || "";

  const alternatingCheckbox = document.querySelector(
    ".schedule-container #is-emergency"
  );
  alternatingCheckbox.checked = !!data.IsAlternatingGroup;
  alternatingCheckbox.dispatchEvent(new Event("change"));

  const dateInputContainer = document.getElementById("dateInputContainer");
  if (data.IsAlternatingGroup) {
    dateInputContainer.style.maxHeight = dateInputContainer.scrollHeight + "px";
    dateInputContainer.style.opacity = "1";
  } else {
    dateInputContainer.style.maxHeight = "0";
    dateInputContainer.style.opacity = "0";
  }

  if (data.StartDate) {
    document.querySelector(".schedule-container #even-start").value = new Date(
      data.StartDate
    )
      .toISOString()
      .split("T")[0];
  }

  clearSchedule();

  if (data.EvenWeek) {
    try {
      const evenWeek =
        typeof data.EvenWeek === "string"
          ? JSON.parse(data.EvenWeek)
          : data.EvenWeek;
      fillWeekSchedule(evenWeek, "even-week");
    } catch (e) {
      console.warn("Ошибка при обработке данных для четной недели:", e);
    }
  }

  if (data.IsAlternatingGroup && data.OddWeek) {
    try {
      const oddWeek =
        typeof data.OddWeek === "string"
          ? JSON.parse(data.OddWeek)
          : data.OddWeek;
      fillWeekSchedule(oddWeek, "odd-week");
    } catch (e) {
      console.warn("Ошибка при обработке данных для нечетной недели:", e);
    }
  }

  const evenWeekButton = document.getElementById("even-week");
  const oddWeekButton = document.getElementById("odd-week");
  const allScheduleButton = document.getElementById("all-schedule");

  if (data.IsAlternatingGroup) {
    evenWeekButton.style.display = "block";
    oddWeekButton.style.display = "block";
    allScheduleButton.style.display = "none";
    evenWeekButton.classList.add("selected");
    oddWeekButton.classList.remove("selected");
  } else {
    evenWeekButton.style.display = "none";
    oddWeekButton.style.display = "none";
    allScheduleButton.style.display = "block";
    allScheduleButton.classList.add("selected");
  }

  const joinExitButton = document.querySelector(
    ".schedule-container .join-exit-button"
  );
  if (joinExitButton) {
    const groupId = parseInt(data.ID, 10);
    const isMember = groupId === currentUserGroupId;
    updateButtonState(joinExitButton, isMember, groupId);
  }

  document
    .querySelectorAll(".schedule-container input, .schedule-container textarea")
    .forEach((el) => {
      el.disabled = true;
    });

  document.querySelectorAll(".day-header .toggle-day").forEach((button) => {
    button.addEventListener("click", function () {
      const weekDay = this.closest(".week-day");
      weekDay.classList.toggle("collapsed");
    });
  });
}

function updateButtonState(button, isMember, groupId) {
  button.textContent = isMember ? "Выйти" : "Вступить";
  button.dataset.groupId = groupId;
  button.onclick = (e) => {
    e.preventDefault();
    handleJoinExitClick(groupId);
  };
}

function clearSchedule() {
  ["even-week", "odd-week"].forEach((weekType) => {
    for (let i = 0; i < 7; i++) {
      const dayContent = document.querySelector(
        `#${weekType}-${i} .day-content`
      );
      if (dayContent) {
        dayContent.innerHTML = "";
      }
    }
  });
}

function fillWeekSchedule(weekData, weekType) {
  const days = [
    "monday",
    "tuesday",
    "wednesday",
    "thursday",
    "friday",
    "saturday",
    "sunday",
  ];

  days.forEach((day, index) => {
    const dayContent = document.querySelector(
      `#${weekType}-${index} .day-content`
    );
    if (dayContent) {
      const lessons = weekData[day] || [];
      if (lessons.length > 0) {
        let tableHtml = "<table>";
        lessons.forEach((lesson, lessonIndex) => {
          tableHtml += `
            <tr>
              <td class="number">${lessonIndex + 1}</td>
              <td class="time"><textarea class="time-input" disabled>${
                lesson[0] || ""
              }</textarea></td>
              <td class="subject"><textarea class="subject-input" disabled>${
                lesson[1] || ""
              }</textarea></td>
              <td class="audience"><textarea class="audience-input" disabled>${
                lesson[2] || ""
              }</textarea></td>
            </tr>
          `;
        });
        tableHtml += "</table>";
        dayContent.innerHTML = tableHtml;
      } else {
        dayContent.innerHTML = "<p>Нет занятий</p>";
      }
    }
  });
}

function handleJoinExitClick(groupId) {
  if (isRequestInProgress) return;

  isRequestInProgress = true;

  const joinToId = parseInt(groupId, 10);
  if (isNaN(joinToId)) {
    console.error("Некорректный ID группы:", groupId);
    alert("Ошибка: некорректный ID группы");
    isRequestInProgress = false;
    return;
  }

  const isMember = joinToId === currentUserGroupId;
  const url = isMember ? "/group/exit" : "/group/join";

  const data = JSON.stringify({
    initData: initData,
    "join-to": joinToId,
  });

  console.log("Отправляемые данные:", data);

  const button = document.querySelector(
    ".schedule-container .join-exit-button"
  );
  button.disabled = true;

  sendPostRequestAsync(url, data)
    .then((result) => {
      if (result.success) {
        const response = JSON.parse(result.response.responseText);
        if (response.status) {
          currentUserGroupId = isMember ? null : joinToId;
          updateButtonState(button, !isMember, joinToId);
          alert(
            response.message ||
              (isMember ? "Вы вышли из группы" : "Вы вступили в группу")
          );
          SetGroups();
        } else {
          alert("Ошибка: " + (response.message || "Неизвестная ошибка"));
        }
      } else {
        console.error("Ответ сервера:", result.response);
        alert(
          `Ошибка при отправке запроса: ${result.response.status} ${result.response.statusText}`
        );
      }
    })
    .catch((error) => {
      console.error("Произошла ошибка:", error);
      alert("Произошла ошибка при обработке запроса");
    })
    .finally(() => {
      button.disabled = false;
      isRequestInProgress = false;
    });
}

function handleGroupButton(button) {
  let groupButtonId = parseInt(button.getAttribute("group_button_id"), 10);
  let originalButtonText = button.textContent.trim();
  let action = originalButtonText === "Выйти" ? "exit" : "join";

  function showConfirmDialog(message) {
    return confirm(message);
  }

  function updateAllButtons(newState, groupId) {
    document
      .querySelectorAll(".join-group, .schedule-container .join-exit-button")
      .forEach((btn) => {
        let btnGroupId = parseInt(
          btn.getAttribute("group_button_id") || btn.dataset.groupId,
          10
        );
        if (newState === "join" && btnGroupId === groupId) {
          btn.textContent = "Выйти";
        } else {
          btn.textContent = "Вступить";
        }
      });
  }

  function updateMemberIcon(groupId, addIcon) {
    let allGroupsSection = document.querySelector(".search-group .groups");
    let userGroupsSection = document.querySelector(".yours-group .groups");

    [allGroupsSection, userGroupsSection].forEach((section) => {
      if (section) {
        document.querySelectorAll(".group-card").forEach((card) => {
          let cardGroupId = parseInt(card.getAttribute("group-id"), 10);
          let existingIcon = card.querySelector(".member-img");
          if (addIcon && cardGroupId === groupId && !existingIcon) {
            let icon = document.createElement("img");
            icon.className = "member-img";
            icon.src = "/static/img/member.svg";
            icon.alt = "";
            card.insertBefore(icon, card.firstChild);
          } else if ((!addIcon || cardGroupId !== groupId) && existingIcon) {
            existingIcon.remove();
          }
        });
      }
    });
  }

  function sendRequest(action) {
    button.classList.add("loading");
    button.textContent = "";

    let url = action === "join" ? "/group/join" : "/group/exit";
    let data = JSON.stringify({
      initData: initData,
      "join-to": groupButtonId,
    });

    sendPostRequestAsync(url, data)
      .then((result) => {
        if (result.success) {
          const jsonResponse = JSON.parse(result.response.responseText);
          if (jsonResponse.status) {
            if (action === "join") {
              currentUserGroupId = groupButtonId;
              updateAllButtons("join", groupButtonId);
              updateMemberIcon(groupButtonId, true);
            } else {
              currentUserGroupId = null;
              updateAllButtons("exit", groupButtonId);
              updateMemberIcon(groupButtonId, false);
            }
            alert(
              jsonResponse.message ||
                `Вы успешно ${
                  action === "join" ? "вступили в" : "вышли из"
                } группу`
            );
            SetGroups(); // Обновляем список групп
          } else {
            alert("Ошибка: " + (jsonResponse.message || "Неизвестная ошибка"));
            updateAllButtons(
              action === "join" ? "exit" : "join",
              groupButtonId
            );
          }
        } else {
          alert(
            `Ошибка при отправке запроса: ${result.response.status} ${result.response.statusText}`
          );
          updateAllButtons(action === "join" ? "exit" : "join", groupButtonId);
        }
      })
      .catch((error) => {
        console.error("Произошла ошибка:", error);
        alert("Произошла ошибка при обработке запроса");
        updateAllButtons(action === "join" ? "exit" : "join", groupButtonId);
      })
      .finally(() => {
        button.classList.remove("loading");
      });
  }

  let confirmMessage =
    action === "join"
      ? "Вы уверены, что хотите вступить в эту группу?"
      : "Вы уверены, что хотите выйти из этой группы?";

  if (showConfirmDialog(confirmMessage)) {
    sendRequest(action);
  }
}

function SetGroups() {
  function SetHtmlCards(groups, userGroupId) {
    currentUserGroupId = parseInt(userGroupId, 10);
    const targetElement = document.getElementById("cards-placeholder");
    if (!targetElement) {
      console.error("Элемент cards-placeholder не найден в DOM");
      return;
    }

    targetElement.innerHTML = "";

    if (!Array.isArray(groups) || groups.length === 0) {
      targetElement.innerHTML = "<p>Группы не найдены</p>";
      return;
    }

    for (const group of groups) {
      if (!group || typeof group !== "object") {
        console.error("Некорректные данные группы:", group);
        continue;
      }

      let htmlString = `<div class="group-card" group-id="${group.ID || ""}">`;
      let buttonMsg = "Вступить";

      if (userGroupId == group.ID) {
        htmlString += `<img class="member-img" src="/static/img/member.svg" alt="" />`;
        buttonMsg = "Выйти";
      }

      htmlString += `
        <h2 class="group-name">${group.Name || "Без названия"}</h2>
        <div class="divider"></div>
        <h2 class="desc">${group.Title || "Без описания"}</h2>
        <button class="join-group" id="enter-button" group_button_id="${
          group.ID || ""
        }">${buttonMsg}</button>
      </div>`;

      targetElement.insertAdjacentHTML("beforeend", htmlString);
    }
  }

  sendPostRequestAsync("/get/groups", JSON.stringify({ initData: initData }))
    .then((result) => {
      if (result.success) {
        try {
          const jsonResponse = JSON.parse(result.response.responseText);
          if (jsonResponse && jsonResponse.groups) {
            SetHtmlCards(jsonResponse.groups, jsonResponse["consists-of"]);
          } else {
            console.error("Некорректный формат данных в ответе сервера");
          }
        } catch (e) {
          console.error("Ошибка при парсинге JSON:", e);
        }
      } else {
        console.error(
          "Ошибка запроса:",
          result.response.status,
          result.response.statusText
        );
      }
    })
    .catch((error) => {
      console.error("Произошла ошибка при отправке запроса:", error);
    });
}

function SetMyGroups() {
  function createGroupCard(group, userGroupId) {
    const isMember = group.ID === userGroupId;
    let htmlString = `
      <div class="group-card" group-id="${group.ID}">
        ${
          isMember
            ? '<img class="member-img" src="/static/img/member.svg" alt="" />'
            : ""
        }
        <h2 class="group-name">${group.Name}</h2>
        <div class="divider"></div>
        <h2 class="desc">${group.Title}</h2>
        <button class="join-group" id="enter-button" group_button_id="${
          group.ID
        }">
          ${isMember ? "Выйти" : "Вступить"}
        </button>
        <button class="edit-group" id="edit-button" group_button_id="${
          group.ID
        }">Редактировать</button>
      </div>
    `;
    return htmlString;
  }

  function SetHtmlCards(groups, userGroupId) {
    const targetElement = document.querySelector(".yours-group .groups");
    if (!targetElement) {
      console.error("Элемент .yours-group .groups не найден в DOM");
      return;
    }

    try {
      targetElement.innerHTML = "";

      if (!Array.isArray(groups) || groups.length === 0) {
        targetElement.innerHTML = "<p>У вас нет групп</p>";
        return;
      }

      for (const group of groups) {
        if (!group || typeof group !== "object") {
          console.warn("Некорректные данные группы:", group);
          continue;
        }
        const htmlString = createGroupCard(group, userGroupId);
        targetElement.insertAdjacentHTML("beforeend", htmlString);
      }

      document.querySelectorAll(".edit-group").forEach((button) => {
        button.addEventListener("click", function () {
          const groupId = this.getAttribute("group_button_id");
          console.log(`Редактирование группы с ID: ${groupId}`);
          // Здесь логика для редактирования группы
        });
      });
    } catch (error) {
      console.error("Ошибка при установке HTML карточек:", error);
    }
  }

  sendPostRequestAsync("/get/my-groups", JSON.stringify({ initData: initData }))
    .then((result) => {
      if (result.success) {
        try {
          const jsonResponse = JSON.parse(result.response.responseText);
          console.log("Полученные данные о группах:", jsonResponse);
          if (jsonResponse && Array.isArray(jsonResponse.groups)) {
            SetHtmlCards(jsonResponse.groups, jsonResponse["consists-of"]);
          } else {
            console.error("Некорректный формат данных в ответе сервера");
          }
        } catch (e) {
          console.error("Ошибка при парсинге JSON:", e);
        }
      } else {
        console.error(
          "Ошибка запроса:",
          result.response.status,
          result.response.statusText
        );
      }
    })
    .catch((error) => {
      console.error("Произошла ошибка при отправке запроса:", error);
    });
}

function initializeScheduleUI() {
  const weekSelector = document.getElementById("week-selector");
  const weekAlternationCheckbox = document.getElementById("is-emergency");

  function updateScheduleView() {
    const isAlternating = weekAlternationCheckbox.checked;
    const evenWeekButton = document.getElementById("even-week");
    const oddWeekButton = document.getElementById("odd-week");
    const allScheduleButton = document.getElementById("all-schedule");
    const evenWeekSchedule = document.querySelector(".even-week-schedule");
    const oddWeekSchedule = document.querySelector(".odd-week-schedule");
    const titleWeeks = document.querySelectorAll(".title-week");

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

  updateScheduleView();
}

function initializeShareButton() {
  const shareBtn = document.getElementById("share-btn");
  const shareVariants = document.getElementById("share-variants");
  let isTooltipVisible = false;
  let isFirstClickAfterShow = false;

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

  shareBtn.addEventListener("click", (event) => {
    event.stopPropagation();
    toggleTooltip();
  });

  document.addEventListener("click", (event) => {
    if (isTooltipVisible) {
      if (isFirstClickAfterShow) {
        isFirstClickAfterShow = false;
      } else if (!shareVariants.contains(event.target)) {
        toggleTooltip();
      }
    }
  });

  shareVariants.addEventListener("click", (event) => {
    event.stopPropagation();
  });
}

function initializeCreateSchedule() {
  const createScheduleContainer = document.querySelector(
    ".create-shedule .schedule-container"
  );
  if (!createScheduleContainer) return;

  const groupNameInput = createScheduleContainer.querySelector(
    'input[placeholder="Имя"]'
  );
  const groupDescInput = createScheduleContainer.querySelector(
    'input[placeholder="Подпись"]'
  );
  const alternatingCheckbox =
    createScheduleContainer.querySelector("#is-emergency");
  const dateInput = createScheduleContainer.querySelector("#even-start");
  const saveButton = createScheduleContainer.querySelector(".save-schedule");

  groupNameInput.disabled = false;
  groupDescInput.disabled = false;
  alternatingCheckbox.disabled = false;
  dateInput.disabled = false;

  alternatingCheckbox.addEventListener("change", function () {
    const dateInputContainer = document.getElementById("dateInputContainer");
    const oddWeekSchedule = document.querySelector(".odd-week-schedule");
    if (this.checked) {
      dateInputContainer.style.maxHeight =
        dateInputContainer.scrollHeight + "px";
      dateInputContainer.style.opacity = "1";
      oddWeekSchedule.style.display = "block";
    } else {
      dateInputContainer.style.maxHeight = "0";
      dateInputContainer.style.opacity = "0";
      oddWeekSchedule.style.display = "none";
    }
  });

  document
    .querySelectorAll(".create-shedule .add-subject")
    .forEach((button) => {
      button.addEventListener("click", function () {
        const weekType = this.id.includes("even") ? "even" : "odd";
        const dayIndex = this.id.split("-").pop();
        addSubject(weekType, dayIndex);
      });
    });

  saveButton.addEventListener("click", function () {
    const scheduleData = collectScheduleData();
    sendScheduleData(scheduleData);
  });
}

function addSubject(weekType, dayIndex) {
  const weekDay = document.querySelector(
    `.create-shedule #${weekType}-week-${dayIndex}`
  );
  if (!weekDay) return;

  const dayContent = weekDay.querySelector(".day-content");
  let table = dayContent.querySelector("table");

  if (!table) {
    table = document.createElement("table");
    dayContent.insertBefore(table, dayContent.firstChild);
  }

  const newRow = table.insertRow();
  const cellsContent = [
    table.rows.length,
    '<textarea class="time-input"></textarea>',
    '<textarea class="subject-input"></textarea>',
    '<textarea class="audience-input"></textarea>',
    `<div class="mini-menu">
      <img src="/static/img/add-menu.svg" alt="" class="open-mini-menu" />
      <div class="opened-menu">
        <img class="arrow" src="/static/img/mini-menu/arrow-0.svg" alt="" />
        <img class="delete" src="/static/img/mini-menu/delete.svg" alt="" />
        <img class="arrow" src="/static/img/mini-menu/arrow-1.svg" alt="" />
      </div>
    </div>`,
  ];

  cellsContent.forEach((content, index) => {
    const cell = newRow.insertCell();
    cell.innerHTML = content;
    cell.className = ["number", "time", "subject", "audience", "mini-menu"][
      index
    ];
  });

  addMiniMenuHandlers(newRow.querySelector(".mini-menu"));
}

function addMiniMenuHandlers(miniMenu) {
  if (!miniMenu) return;

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

function deleteRow(row) {
  const table = row.closest("table");
  row.remove();
  recalculateRowNumbers(table);
}

function recalculateRowNumbers(table) {
  Array.from(table.rows).forEach((row, index) => {
    const numberCell = row.cells[0];
    if (numberCell && numberCell.classList.contains("number")) {
      numberCell.textContent = index + 1;
    }
  });
}

function collectScheduleData() {
  const container = document.querySelector(
    ".create-shedule .schedule-container"
  );
  const groupName = container.querySelector('input[placeholder="Имя"]').value;
  const groupDesc = container.querySelector(
    'input[placeholder="Подпись"]'
  ).value;
  const isAlternating = container.querySelector("#is-emergency").checked;
  const startDate = container.querySelector("#even-start").value;

  const scheduleData = {
    groupName,
    groupDesc,
    isAlternating,
    startDate,
    evenWeek: {},
    oddWeek: {},
  };
  ["even", "odd"].forEach((weekType) => {
    for (let i = 0; i < 7; i++) {
      const dayContent = container.querySelector(
        `#${weekType}-week-${i} .day-content`
      );
      const table = dayContent.querySelector("table");
      if (table) {
        scheduleData[`${weekType}Week`][i] = Array.from(table.rows).map(
          (row) => ({
            time: row.cells[1].querySelector("textarea").value,
            subject: row.cells[2].querySelector("textarea").value,
            audience: row.cells[3].querySelector("textarea").value,
          })
        );
      }
    }
  });

  return scheduleData;
}

function sendScheduleData(scheduleData) {
  console.log("Отправка данных расписания:", scheduleData);
  sendPostRequestAsync("/create/group", JSON.stringify(scheduleData))
    .then((result) => {
      if (result.success) {
        alert("Группа успешно создана!");
        SetGroups();
        hideScheduleContainer();
      } else {
        alert(
          "Ошибка при создании группы: " +
            (result.response.responseText || "Неизвестная ошибка")
        );
      }
    })
    .catch((error) => {
      console.error("Ошибка при отправке данных:", error);
      alert("Произошла ошибка при создании группы");
    });
}

function hideScheduleContainer() {
  document.querySelector(".schedule-container").style.display = "none";
  document.querySelector(".search-group .groups").style.display = "block";
}

function showScheduleContainer() {
  document.querySelector(".schedule-container").style.display = "flex";
  document.querySelector(".search-group .groups").style.display = "none";
}

function sendPostRequestAsync(url, data) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onload = function () {
      if (this.status >= 200 && this.status < 300) {
        resolve({ success: true, response: xhr });
      } else {
        resolve({ success: false, response: xhr });
      }
    };

    xhr.onerror = function () {
      reject(new Error("Network Error"));
    };

    xhr.send(data);
  });
}

function debugLog(message) {
  console.log("DEBUG: " + message);
}

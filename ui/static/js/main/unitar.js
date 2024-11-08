function sendGetRequestAsync(url, urlQuery) {
  return new Promise((resolve, reject) => {
    var Request = false;

    if (window.XMLHttpRequest) {
      Request = new XMLHttpRequest();
    } else if (window.ActiveXObject) {
      try {
        Request = new ActiveXObject("Microsoft.XMLHTTP");
      } catch (CatchException) {
        try {
          Request = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (CatchException2) {
          Request = false;
        }
      }
    }

    if (!Request) {
      reject(new Error("Невозможно создать XMLHttpRequest"));
      return;
    }

    const fullUrl = urlQuery ? `${url}?${urlQuery}` : url;
    console.log("fullUrl    ", fullUrl);

    Request.open("GET", fullUrl, true);

    Request.onreadystatechange = function () {
      if (Request.readyState === 4) {
        if (Request.status === 200) {
          resolve({ success: true, response: Request.responseText });
        } else {
          let response;
          switch (Request.status) {
            case 404:
              response = "404 (Not Found)";
              break;
            case 403:
              response = "403 (Forbidden)";
              break;
            case 500:
              response = "500 (Internal Server Error)";
              break;
            default:
              response = `${Request.status} (${
                Request.statusText || "Unknown Error"
              })`;
          }
          resolve({ success: false, response: response });
        }
      }
    };

    Request.onerror = function () {
      reject(new Error("Network Error"));
    };

    try {
      Request.send();
    } catch (error) {
      reject(error);
    }
  });
}

function sendPostRequestAsync(url, data) {
  return new Promise((resolve, reject) => {
    var Request = false;

    if (window.XMLHttpRequest) {
      Request = new XMLHttpRequest();
    } else if (window.ActiveXObject) {
      try {
        Request = new ActiveXObject("Microsoft.XMLHTTP");
      } catch (CatchException) {
        try {
          Request = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (CatchException2) {
          Request = false;
        }
      }
    }

    if (!Request) {
      reject(new Error("Невозможно создать XMLHttpRequest"));
      return;
    }

    console.log("URL:", url);
    console.log("Data:", data);

    Request.open("POST", url, true);

    Request.setRequestHeader(
      "Content-Type",
      "application/x-www-form-urlencoded"
    );

    Request.onreadystatechange = function () {
      if (Request.readyState === 4) {
        if (Request.status === 200) {
          resolve({ success: true, response: Request });
        } else {
          let response;
          switch (Request.status) {
            case 404:
              response = "404 (Not Found)";
              break;
            case 403:
              response = "403 (Forbidden)";
              break;
            case 500:
              response = "500 (Internal Server Error)";
              break;
            default:
              response = `${Request.status} (${
                Request.statusText || "Unknown Error"
              })`;
          }
          resolve({ success: false, response: Request });
        }
      }
    };

    Request.onerror = function () {
      reject(new Error("Network Error"));
    };

    try {
      Request.send(data);
    } catch (error) {
      reject(error);
    }
  });
}

/**
 * Функция showAlert отображает уведомление для пользователя с заданным текстом и определённым стилем.
 *
 * @param {string} error_text - Текст уведомления, который будет показан пользователю.
 * @param {number} durationSeconds - Продолжительность анимации (в секундах), после которой уведомление исчезнет.
 * @param {boolean} isError - Флаг, указывающий на тип уведомления. Если true, уведомление будет красным и помечено как ошибка. Если false, уведомление будет зелёным и помечено как успешно.
 */
function showAlert(error_text, durationSeconds, isError) {
  function easeInOutQuad(t) {
    return t < 0.5 ? 2 * t * t : 1 - Math.pow(-2 * t + 2, 2) / 2;
  }

  // Устанавливаем стили и текст уведомления в зависимости от типа сообщения
  if (isError) {
    error_alert.alert_close_button.style.color = "#ff0000"; // Красный цвет для ошибок
    error_alert.title.innerHTML = "ОШИБКА";
  } else {
    error_alert.title.innerHTML = "УСПЕШНО";
    error_alert.alert_close_button.style.color = "#1ede00"; // Зелёный цвет для успешных сообщений
  }

  // Устанавливаем текст уведомления
  error_alert.error_description.innerHTML = error_text;

  // Добавляем обработчик клика для кнопки закрытия уведомления
  error_alert.alert_close_button.addEventListener("click", function () {
    moveElement(error_alert.alert_container, "toRight"); // Перемещаем уведомление за пределы экрана при закрытии
  });

  // Перемещаем уведомление в центр экрана
  moveElement(error_alert.alert_container, "toCenter");

  // Запускаем анимацию ширины уведомления и перемещаем его за пределы экрана по окончании анимации
  animateWidth(
    error_alert.loading,
    durationSeconds,
    easeInOutQuad,
    function () {
      moveElement(error_alert.alert_container, "toRight");
    }
  );
}

/**
 *  animateWidth активирует плавное увеличение ширины обьекта с 0 до 100%
 */
function animateWidth(
  element,
  durationSeconds,
  easingFunction = (t) => t,
  callback
) {
  if (!(element instanceof Element)) {
    console.error("Переданный аргумент не является DOM элементом");
    return;
  }

  const fps = 60;
  const totalFrames = durationSeconds * fps;

  function step(timestamp) {
    if (!step.startTime) step.startTime = timestamp;
    const elapsed = timestamp - step.startTime;
    const progress = Math.min(elapsed / (durationSeconds * 1000), 1);

    const easedProgress = easingFunction(progress);
    const currentWidth = easedProgress * 90;

    element.style.width = currentWidth + "%";

    if (progress < 1) {
      requestAnimationFrame(step);
    } else {
      // Анимация завершена, вызываем callback, если он предоставлен
      if (typeof callback === "function") {
        callback();
      }
    }
  }

  requestAnimationFrame(step);
}

/**
 * Функция moveElement управляет позиционированием HTML-элемента, изменяя его CSS-классы в зависимости от указанного направления.
 *
 * @param {HTMLElement} element - HTML-элемент, который требуется переместить.
 * @param {string} direction - Направление перемещения элемента. Может принимать значения:
 *   - "toCenter": Перемещает элемент в центр экрана, добавляя класс "center-screen" и убирая класс "off-screen".
 *   - "toRight": Перемещает элемент в правую часть экрана, добавляя класс "off-screen" и убирая класс "center-screen".
 */
function moveElement(element, direction) {
  if (direction === "toCenter") {
    element.classList.remove("off-screen");
    element.classList.add("center-screen");
  } else if (direction === "toRight") {
    element.classList.remove("center-screen");
    element.classList.add("off-screen");
  }
}

/**
 * Преобразует все строковые значения в объекте (и его вложенных объектах) в соответствующие HTML элементы.
 * Если элемент с указанным id не найден, выводится предупреждение в консоль.
 *
 * @param {Object} objects - Объект, содержащий другие объекты, значения которых нужно преобразовать
 */
function convertIdsToElements(objects) {
  /**
   * Рекурсивно обрабатывает объект, заменяя строковые значения на HTML элементы
   *
   * @param {Object} obj - Объект для обработки
   */
  function processObject(obj) {
    for (let key in obj) {
      if (typeof obj[key] === "string") {
        const id = obj[key]; // Используем значение как id
        const element = document.getElementById(id);
        if (element) {
          obj[key] = element; // Заменяем строку на HTML элемент
        } else {
          console.warn(`Элемент с id "${id}" не найден для ключа "${key}"`);
        }
      } else if (typeof obj[key] === "object" && obj[key] !== null) {
        processObject(obj[key]); // Рекурсивный вызов для вложенных объектов
      }
    }
  }

  // Обрабатываем каждый объект в переданном objects
  for (let objName in objects) {
    processObject(objects[objName]);
  }
}

/**
 * setTg инициализирует Telegram Web App и возвращает объекты tg и initData
 * @param {number} [maxAttempts=3] Максимальное количество попыток инициализации
 * @param {number} [delay=1000] Задержка между попытками в миллисекундах
 * @returns {Object|null} Объект с свойствами tg и initData или null в случае неудачи
 */
function initializeTg(maxAttempts = 3, delay = 1000) {
  function sleep(ms) {
    const start = Date.now();
    while (Date.now() - start < ms) {}
  }

  for (let attempts = 1; attempts <= maxAttempts; attempts++) {
    if (window.Telegram && window.Telegram.WebApp) {
      const tg = window.Telegram.WebApp;
      const initData = tg.initData;

      // Проверяем, что initData не пустой и не равен ""
      if (!initData || initData === "") {
        console.warn(
          `Попытка ${attempts}: initData пуст или равен "". Возможно, приложение запущено не в Telegram.`
        );
      } else {
        // Вызываем tg.ready() для сообщения Telegram, что приложение готово
        tg.ready();

        // Добавим проверку на поддержку основных методов
        if (
          typeof tg.sendData !== "function" ||
          typeof tg.expand !== "function"
        ) {
          console.warn(
            `Попытка ${attempts}: Некоторые ожидаемые методы Telegram Web App отсутствуют.`
          );
        } else {
          console.log(
            `Telegram Web App успешно инициализирован (попытка ${attempts})`
          );
          return { tg, initData };
        }
      }
    } else {
      console.error(
        `Попытка ${attempts}: Telegram Web App не найден. Проверьте подключение библиотеки telegram-web-app.js.`
      );
    }

    if (attempts < maxAttempts) {
      console.log(`Повторная попытка через ${delay}мс...`);
      sleep(delay);
    }
  }

  console.error(
    `Не удалось инициализировать Telegram Web App после ${maxAttempts} попыток.`
  );
  return null;
}

// Глобальные переменные для отслеживания текущего предупреждения
let currentWarning = null; // Хранит текущий элемент предупреждения
let currentWarningTimeout = null; // Хранит таймер для автоматического скрытия предупреждения

/**
 * Показывает предупреждение рядом с целевым элементом.
 * @param {Element} targetElement - Элемент, рядом с которым нужно показать предупреждение.
 * @param {string} message - Текст предупреждения.
 * @param {string} position - Позиция предупреждения относительно целевого элемента (по умолчанию "top").
 * @param {string} percent - Горизонтальное смещение стрелочки предупреждения в процентах (по умолчанию "50%").
 * @param {number} duration - Продолжительность показа предупреждения в миллисекундах (по умолчанию 3000).
 */
function showWarning(
  targetElement,
  message,
  position = "top",
  percent = "50%",
  duration = 3000
) {
  console.log("showWarning called", {
    targetElement,
    message,
    position,
    percent,
    duration,
  });

  // Удаляем предыдущее предупреждение, если оно существует
  if (currentWarning) {
    currentWarning.remove();
  }

  // Очищаем предыдущий таймер, если он существует
  if (currentWarningTimeout) {
    clearTimeout(currentWarningTimeout);
  }

  // Проверяем, является ли targetElement допустимым DOM-элементом
  if (!targetElement || !(targetElement instanceof Element)) {
    console.error("Invalid targetElement");
    return;
  }

  // Создаем уникальный идентификатор для предупреждения
  const warningId = `warning-${
    targetElement.id || Math.random().toString(36).substr(2, 9)
  }`;

  // Создаем элемент предупреждения
  const warning = document.createElement("div");
  warning.id = warningId;
  warning.className = `warning warning-${position}`;
  warning.textContent = message;
  warning.style.visibility = "hidden";
  warning.style.opacity = "0";
  document.body.appendChild(warning);

  currentWarning = warning;

  console.log("Warning element created", warning);

  // Устанавливаем CSS-переменную для позиционирования псевдоэлемента
  document.documentElement.style.setProperty("--pseudo-left", percent);

  // Используем setTimeout, чтобы дать браузеру время на отрисовку предупреждения
  setTimeout(() => {
    // Получаем позицию целевого элемента
    const targetRect = targetElement.getBoundingClientRect();
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const scrollLeft =
      window.pageXOffset || document.documentElement.scrollLeft;

    // Позиционируем предупреждение
    warning.style.left = `${
      targetRect.left + scrollLeft + targetRect.width / 2
    }px`;
    warning.style.top = `${targetRect.top + scrollTop - 10}px`;
    warning.style.transform = "translate(-50%, -100%)";
    warning.style.visibility = "visible";

    console.log("Warning styles set", {
      position: warning.style.position,
      top: warning.style.top,
      left: warning.style.left,
      transform: warning.style.transform,
      zIndex: warning.style.zIndex,
      pseudoLeft: percent,
    });

    // Анимируем появление предупреждения
    requestAnimationFrame(() => {
      warning.style.opacity = "1";
      targetElement.scrollIntoView({ behavior: "smooth", block: "center" });
    });

    console.log("Warning position:", warning.getBoundingClientRect());
    console.log("Target position:", targetRect);

    // Устанавливаем таймер для автоматического скрытия предупреждения
    if (duration !== Infinity) {
      currentWarningTimeout = setTimeout(() => {
        hideWarning(targetElement);
      }, duration);
    }
  }, 0);
}

/**
 * Скрывает текущее предупреждение.
 * @param {Element} targetElement - Элемент, рядом с которым было показано предупреждение.
 */
function hideWarning(targetElement) {
  if (currentWarning) {
    console.log("Hiding warning", currentWarning);
    // Анимируем исчезновение предупреждения
    currentWarning.style.opacity = "0";
    setTimeout(() => {
      currentWarning.remove();
      currentWarning = null;
      console.log("Warning removed");
    }, 300);
  } else {
    console.log("No active warning to hide");
  }

  // Сбрасываем CSS-переменную для позиционирования псевдоэлемента
  document.documentElement.style.removeProperty("--pseudo-left");
}

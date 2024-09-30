let tg, initData;

document.addEventListener("DOMContentLoaded", function () {
  // Иницилизация телеграмм
  const result = initializeTg(5, 2000); // 5 попыток с интервалом в 2 секунды
  if (result) {
    postAndRenderResponse("/main", {
      initData: result.InitData,
    });
  } else {
    alert("Пожалуйста, запустите форму через телеграмм");
  }
});

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

function postAndRenderResponse(url, data) {
  fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => response.text())
    .then((html) => {
      // Заменяем содержимое страницы полученным HTML
      document.open();
      document.write(html);
      document.close();

      // Если нужно обновить URL без перезагрузки страницы
      window.history.pushState({}, "", url);
    })
    .catch((error) => {
      console.error("Ошибка:", error);
      // Здесь можно добавить обработку ошибок, например, показать сообщение пользователю
    });
}

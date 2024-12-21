/**
 * Пытается разобрать текст как JSON с сообщением об ошибке
 * @param text - текст для парсинга
 * @returns {status: string, message: string} или null
 */
function tryParseErrorMessage(text: string): { status: string; message: string } | null {
  try {
    const json = JSON.parse(text);
    if (json.status && json.message) {
      return json;
    }
  } catch (e) {
    return null;
  }
  return null;
}

/**
 * Базовая функция для отправки POST-запросов
 * @param url - URL endpoint для запроса
 * @param data - Данные для отправки
 * @returns Promise с ответом от сервера
 */
async function postData(url = "", data = {}) {
  // Логируем начало запроса для отладки
  console.log("Отправка POST запроса", { url, data });

  // Выполняем fetch запрос с настройками
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json", // Указываем тип контент как JSON
    },
    credentials: "include", // Включаем куки в запрос для аутентификации
    body: JSON.stringify(data), // Преобразуем объект в JSON строку
  });

  // Логируем полученный ответ
  console.log("Ответ получен", response);

  // Получаем тело ответа в виде текста
  const textResponse = await response.text();
  console.log("Ответ в виде текста", textResponse);

  // Проверяем успешность запроса по статусу
  if (!response.ok) {
    // Пытаемся разобрать ошибку как JSON
    const errorJson = tryParseErrorMessage(textResponse);
    if (errorJson) {
      throw new Error(errorJson.message);
    }
    
    console.error("Ошибка в ответе", response.statusText);
    console.error("Статус ответа:", response.status);
    console.error("Тело ответа:", textResponse);
    throw new Error(`Ошибка HTTP: ${response.status}`);
  }

  // Получаем тип контента из заголовков
  const contentType = response.headers.get("content-type");

  // Если тело ответа пустое, возвращаем пустой объект
  if (!textResponse) {
    console.error("Пустой ответ от сервера");
    return {};
  }

  // Обрабатываем JSON ответ
  if (contentType && contentType.includes("application/json")) {
    try {
      const jsonResponse = JSON.parse(textResponse);
      console.log("Ответ JSON", jsonResponse);
      return jsonResponse;
    } catch (error) {
      console.error("Ошибка при разборе JSON", error);
      console.error("Текст ответа, который не удалось разобрать:", textResponse);
      throw new Error("Некорректный JSON");
    }
  } else {
    // Пытаемся распарсить как JSON даже если content-type не указан
    try {
      const jsonResponse = JSON.parse(textResponse);
      console.log("Успешно распарсили ответ как JSON:", jsonResponse);
      return jsonResponse;
    } catch {
      // Если не получилось распарсить как JSON, возвращаем как текст
      console.log("Ответ не является JSON, возвращаем текст");
      return { text: textResponse };
    }
  }
}

// Типы ответов от сервера
interface AuthResponse {
  success: boolean;
  message: string;
  role?: string;
  name?: string;
}

interface EmailValidationResponse {
  success: boolean;
  message: string;
}

interface RegistrationResponse {
  success: boolean;
  message: string;
}

/**
 * Авторизация пользователя
 */
async function login(email: string, password: string): Promise<AuthResponse> {
  try {
    const response = await postData("/api/v1/auth/login", { email, password });
    return {
      success: response.status === "success",
      message: response.message,
      role: response.role,
      name: response.name
    };
  } catch (error) {
    return {
      success: false,
      message: (error as Error).message
    };
  }
}

/**
 * Проверка текущей авторизации
 */
async function checkAuth(): Promise<AuthResponse> {
  try {
    const response = await postData("/api/v1/auth/jwt", {});
    return {
      success: response.status === "success",
      message: response.message,
      role: response.role,
      name: response.name
    };
  } catch (error) {
    return {
      success: false,
      message: (error as Error).message
    };
  }
}

/**
 * Валидация email при регистрации
 */
async function validateEmail(email: string): Promise<EmailValidationResponse> {
  try {
    const response = await postData("/api/v1/registration/validate-email", { email });
    return {
      success: response.status === "success",
      message: response.message
    };
  } catch (error) {
    return {
      success: false,
      message: (error as Error).message
    };
  }
}

/**
 * Начало регистрации (установка пароля)
 */
async function startRegistration(email: string, password: string): Promise<RegistrationResponse> {
  try {
    const response = await postData("/api/v1/registration/start", { email, password });
    return {
      success: response.status === "success",
      message: response.message
    };
  } catch (error) {
    return {
      success: false,
      message: (error as Error).message
    };
  }
}

/**
 * Подтверждение регистрации кодом
 */
async function confirmRegistration(email: string, code: string): Promise<AuthResponse> {
  try {
    const response = await postData("/api/v1/registration/confirm", { email, code });
    return {
      success: response.status === "success",
      message: response.message,
      role: response.role,
      name: response.name
    };
  } catch (error) {
    return {
      success: false,
      message: (error as Error).message
    };
  }
}

/**
 * Выход из системы
 */
async function logout(): Promise<{ success: boolean; message: string }> {
  try {
    const response = await postData("/api/v1/auth/logout", {});
    if (response.status === "success") {
      window.location.reload(); // Перезагрузка страницы после успешного выхода
    }
    return {
      success: response.status === "success",
      message: response.message
    };
  } catch (error) {
    return {
      success: false,
      message: (error as Error).message
    };
  }
}

export {
  postData,
  login,
  checkAuth,
  validateEmail,
  startRegistration,
  confirmRegistration,
  logout
};

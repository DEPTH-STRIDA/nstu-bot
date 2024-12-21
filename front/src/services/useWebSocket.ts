import { useEffect, useRef, useState, useCallback } from "react";

// Хук для работы с WebSocket соединением
export function useWebSocket<T>(url: string) {
  // Состояние для хранения полученных данных от сервера
  const [data, setData] = useState<T | null>(null);
  // Состояние для хранения ошибок
  const [error, setError] = useState<string | null>(null);
  // Состояние подключения (подключены/отключены)
  const [isConnected, setIsConnected] = useState(false);
  
  // Реф для хранения экземпляра WebSocket
  // Используем useRef чтобы сохранить соединение между рендерами
  const ws = useRef<WebSocket | null>(null);
  // Реф для хранения таймера переподключения
  const reconnectTimeout = useRef<number>();

  // Функция для установки соединения
  // useCallback используется для мемоизации функции
  const connect = useCallback(() => {
    // Если соединение уже открыто, ничего не делаем
    if (ws.current?.readyState === WebSocket.OPEN) {
      return;
    }

    console.log("Connecting to WebSocket:", url);
    try {
      // Создаем новое WebSocket соединение
      ws.current = new WebSocket(url);

      // Обработчик успешного подключения
      ws.current.onopen = () => {
        console.log("WebSocket connected");
        setIsConnected(true);
        setError(null);
      };

      // Обработчик входящих сообщений
      ws.current.onmessage = (event) => {
        try {
          // Парсим JSON из полученного сообщения
          const response = JSON.parse(event.data);
          console.log("WebSocket message received:", response);
          
          // Если сервер вернул ошибку
          if (response.status === "error") {
            setError(response.message);
          } else {
            // Если все ок, сохраняем данные
            setData(response);
          }
        } catch (e) {
          console.error("Error parsing WebSocket message:", e);
        }
      };

      // Обработчик ошибок соединения
      ws.current.onerror = (event) => {
        console.error("WebSocket error:", event);
        setError("Connection error");
        setIsConnected(false);
      };

      // Обработчик закрытия соединения
      ws.current.onclose = (event) => {
        console.log("WebSocket closed:", event);
        setIsConnected(false);

        // Если соединение закрылось не по нашей инициативе (код != 1000)
        // пытаемся переподключиться через 3 секунды
        if (event.code !== 1000) {
          reconnectTimeout.current = setTimeout(connect, 3000);
        }
      };
    } catch (e) {
      console.error("Failed to create WebSocket:", e);
      setError("Failed to connect");
    }
  }, [url]);

  // Эффект для управления жизненным циклом соединения
  useEffect(() => {
    // При монтировании компонента устанавливаем соединение
    connect();

    // При размонтировании компонента
    return () => {
      // Очищаем таймер переподключения если он есть
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }
      // Закрываем соединение если оно открыто
      if (ws.current) {
        ws.current.close(1000); // 1000 - нормальное закрытие
      }
    };
  }, [connect]);

  // Функция для отправки сообщений
  const sendMessage = useCallback((message: any) => {
    // Проверяем что соединение существует и открыто
    if (!ws.current || ws.current.readyState !== WebSocket.OPEN) {
      console.error("WebSocket is not connected");
      return;
    }

    try {
      // Преобразуем сообщение в строку и отправляем
      const messageStr = JSON.stringify(message);
      console.log("Sending WebSocket message:", messageStr);
      ws.current.send(messageStr);
    } catch (e) {
      console.error("Error sending message:", e);
      setError("Failed to send message");
    }
  }, []);

  // Возвращаем объект с данными и методами для использования в компонентах
  return { 
    data,           // последние полученные данные
    error,          // текущая ошибка если есть
    isConnected,    // статус подключения
    sendMessage     // функция для отправки сообщений
  };
}

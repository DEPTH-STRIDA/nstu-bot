// Взаимодействие с горизонтальным и вертикальным меню
document.addEventListener("DOMContentLoaded", () => {
  // Слайдер секций
  const sliderWrapper = document.querySelector(".slider-wrapper");
  // Сами секции для свайпа
  const sections = document.querySelectorAll(".slider-section");
  // Горизонтальное меню
  const scrollMenu = document.querySelector(".scroll-menu");
  // Пункты горизонтального меню
  const menuItems = document.querySelectorAll(".scroll-menu-ul li");
  // Вертикальное меню
  const sideMenu = document.querySelector(".sidebar-menu");
  // Пункты вертикального меню
  const sideMenuItems = document.querySelectorAll(".sidebar-menu li");
  // Пункты меню с главной секции. Кнопки
  const mainMenuItems = document.querySelectorAll(".button-card");

  // Начальная секция 1 - главная (0 - поиск аудитории)
  let currentIndex = 1;
  let startX, startY, moveX, moveY, lastMove;
  let isMoving = false;
  // Необходимая длина свайпа для смены секции
  const threshold = window.innerWidth * 0.2;
  let isVerticalScroll = false;
  let isSideMenuOpening = false;

  // Индексы для перехода на соотвествующую секцию при нажатии (из-за другого порядка в вертикальном меню)
  const sideMenuMap = {
    "Поиск аудитории": 0,
    Главная: 1,
    "Поиск группы": 2,
    "Мои группы": 3,
    "Создать группу": 4,
  };

  function setTransform(index) {
    sliderWrapper.style.transform = `translateX(-${index * 100}%)`;
  }

  function updateMenu(index) {
    menuItems.forEach((item, i) => {
      item.classList.toggle("selected", i === index);
    });

    const selectedItem = menuItems[index];
    const menuWidth = scrollMenu.offsetWidth;
    const itemWidth = selectedItem.offsetWidth;
    const scrollLeft = selectedItem.offsetLeft - menuWidth / 2 + itemWidth / 2;
    scrollMenu.scrollTo({ left: scrollLeft, behavior: "smooth" });
  }

  function showSideMenu() {
    sideMenu.style.transform = "translateX(0)";
  }

  function hideSideMenu() {
    sideMenu.style.transform = "translateX(-100%)";
  }

  function handleTouchStart(e) {
    startX = e.type.includes("mouse") ? e.clientX : e.touches[0].clientX;
    startY = e.type.includes("mouse") ? e.clientY : e.touches[0].clientY;
    isMoving = true;
    lastMove = startX;
    isVerticalScroll = false;
    isSideMenuOpening = false;
    sliderWrapper.style.transition = "none";
  }

  function handleTouchMove(e) {
    if (!isMoving) return;
    moveX = e.type.includes("mouse") ? e.clientX : e.touches[0].clientX;
    moveY = e.type.includes("mouse") ? e.clientY : e.touches[0].clientY;
    const diffX = moveX - startX;
    const diffY = moveY - startY;

    if (!isVerticalScroll && Math.abs(diffY) > Math.abs(diffX)) {
      isVerticalScroll = true;
    }

    if (isVerticalScroll) return;

    const movePercent = (diffX / window.innerWidth) * 100;

    if (Math.abs(movePercent) < 1) return; // Ignore very small movements

    if (currentIndex === 0 && diffX > 0) {
      // Show side menu gradually only for "Поиск аудитории" section
      const menuMovePercent = Math.min(movePercent, 100);
      sideMenu.style.transform = `translateX(${menuMovePercent - 100}%)`;
      isSideMenuOpening = true;
    } else {
      sliderWrapper.style.transform = `translateX(calc(-${
        currentIndex * 100
      }% + ${movePercent}%))`;
      isSideMenuOpening = false;
    }

    lastMove = moveX;

    // Prevent default behavior for mouse events to avoid text selection
    if (e.type.includes("mouse")) {
      e.preventDefault();
    }
  }

  function handleTouchEnd(e) {
    if (!isMoving) return;
    isMoving = false;
    const diff = lastMove - startX;

    sliderWrapper.style.transition = "transform 0.3s ease-out";
    sideMenu.style.transition = "transform 0.3s ease-out";

    if (!isVerticalScroll && Math.abs(diff) >= threshold) {
      if (diff > 0 && currentIndex > 0) {
        currentIndex--;
      } else if (diff < 0 && currentIndex < sections.length - 1) {
        currentIndex++;
      }
    }

    if (
      isSideMenuOpening &&
      currentIndex === 0 &&
      diff > 0 &&
      !isVerticalScroll
    ) {
      showSideMenu();
    } else {
      hideSideMenu();
    }
    setTransform(currentIndex);
    updateMenu(currentIndex);
  }

  function goToSection(index) {
    if (index >= 0 && index < sections.length) {
      currentIndex = index;
      setTransform(currentIndex);
      updateMenu(currentIndex);
    }
    hideSideMenu();
  }

  // Initialize
  setTransform(currentIndex);
  updateMenu(currentIndex);
  hideSideMenu();

  // Event listeners
  sliderWrapper.addEventListener("touchstart", handleTouchStart);
  sliderWrapper.addEventListener("touchmove", handleTouchMove);
  sliderWrapper.addEventListener("touchend", handleTouchEnd);

  // Mouse events
  sliderWrapper.addEventListener("mousedown", handleTouchStart);
  sliderWrapper.addEventListener("mousemove", handleTouchMove);
  sliderWrapper.addEventListener("mouseup", handleTouchEnd);
  sliderWrapper.addEventListener("mouseleave", handleTouchEnd);

  menuItems.forEach((item, index) => {
    item.addEventListener("click", () => goToSection(index));
  });

  mainMenuItems.forEach((item, index) => {
    if (index >= 1) {
      index++;
    }
    item.addEventListener("click", () => goToSection(index));
  });

  sideMenuItems.forEach((item) => {
    item.addEventListener("click", (e) => {
      const itemText = item.textContent.trim();
      if (sideMenuMap.hasOwnProperty(itemText)) {
        goToSection(sideMenuMap[itemText]);
      }
      e.stopPropagation();
    });
  });

  document.addEventListener("click", (e) => {
    if (
      !sideMenu.contains(e.target) &&
      !e.target.classList.contains("open-menu-button")
    ) {
      hideSideMenu();
    }
  });

  const openMenuButton = document.querySelector(".open-menu-button");
  if (openMenuButton) {
    openMenuButton.addEventListener("click", showSideMenu);
  }
});

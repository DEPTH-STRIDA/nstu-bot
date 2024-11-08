// Сортировка групп при поиске
document.addEventListener("DOMContentLoaded", function () {
  function setSearchStyleSafely(searchInputId, searchImgId) {
    try {
      setSearchStyle(searchInputId, searchImgId);
      console.log(
        `Слушатели успешно установлены для ${searchInputId} и ${searchImgId}`
      );
    } catch (error) {
      console.error(
        `Ошибка при установке слушателей для ${searchInputId} и ${searchImgId}:`,
        error
      );
    }
  }

  function setSearchStyle(searchInputId, searchImgId) {
    // Получаем ссылку на элементы
    const searchInput = document.getElementById(searchInputId);
    const searchImg = document.getElementById(searchImgId);

    // Проверяем, существуют ли элементы
    if (!searchInput || !searchImg) {
      throw new Error(
        `Элементы не найдены: input (${searchInputId}) или img (${searchImgId})`
      );
    }

    // Добавляем слушатель события на получение фокуса
    searchInput.addEventListener("focus", function () {
      console.log("focus");
      searchInput.style.textAlign = "left";
      searchImg.style.display = "none";
    });

    // Добавляем слушатель события на потерю фокуса
    searchInput.addEventListener("blur", function () {
      console.log("Input потерял фокус");
      if (searchInput.value === "") {
        searchInput.style.textAlign = "center";
        searchImg.style.display = "block";
      }
    });
  }

  setSearchStyleSafely("search-group", "search-group-img");
  setSearchStyleSafely("search-audience", "search-audience-img");
  setSearchStyleSafely("search-yours-group", "search-yours-group-img");
});

document.addEventListener("DOMContentLoaded", function () {
  function searchAndSortGroups(inputElement, groupsContainer) {
    inputElement.addEventListener("input", function () {
      const searchTerm = this.value.toLowerCase();
      const groupCards = Array.from(
        groupsContainer.querySelectorAll(".group-card")
      );

      groupCards.sort((a, b) => {
        const aName = a.querySelector(".group-name").textContent.toLowerCase();
        const bName = b.querySelector(".group-name").textContent.toLowerCase();
        const aDesc = a.querySelector(".desc").textContent.toLowerCase();
        const bDesc = b.querySelector(".desc").textContent.toLowerCase();

        const aMatch = aName.includes(searchTerm) || aDesc.includes(searchTerm);
        const bMatch = bName.includes(searchTerm) || bDesc.includes(searchTerm);

        if (aMatch && !bMatch) return -1;
        if (!aMatch && bMatch) return 1;

        // Если оба соответствуют или оба не соответствуют, сохраняем исходный порядок
        return groupCards.indexOf(a) - groupCards.indexOf(b);
      });

      groupCards.forEach((card) => groupsContainer.appendChild(card));
    });
  }

  // Применяем функцию к секции "search-group"
  const searchGroupInput = document.querySelector(
    '.slider-section.search-group input[name="search-group"]'
  );
  const searchGroupContainer = document.querySelector(
    ".slider-section.search-group .groups"
  );
  if (searchGroupInput && searchGroupContainer) {
    searchAndSortGroups(searchGroupInput, searchGroupContainer);
  }

  // Применяем функцию к секции "yours-group"
  const yoursGroupInput = document.querySelector(
    '.slider-section.yours-group input[name="search-yours-group"]'
  );
  const yoursGroupContainer = document.querySelector(
    ".slider-section.yours-group .groups"
  );
  if (yoursGroupInput && yoursGroupContainer) {
    searchAndSortGroups(yoursGroupInput, yoursGroupContainer);
  }
});

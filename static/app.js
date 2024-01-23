const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  loadMore: (ev) => {
    ev.preventDefault();
    fetch('/loadMore').then((response) => {
      response.json().then((results) => {
        const table = document.getElementById("table-body");
        let existingRows = table.innerHTML; // Get current HTML of the table
        for (let result of results) {
          existingRows += `<tr><td>${result}</td></tr>`; // Append new rows
        }
        table.innerHTML = existingRows;
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    table.innerHTML = rows;
  },
};

const form = document.getElementById("form");
const loadMoreBtn = document.getElementById("load-more");
form.addEventListener("submit", Controller.search);
loadMoreBtn.addEventListener("click", Controller.loadMore);

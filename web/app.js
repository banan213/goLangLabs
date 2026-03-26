function showError(msg) {
  document.getElementById('error').textContent = msg;
}

function showPhoto(url) {
  showError('');
  const container = document.getElementById('photo-container');
  container.innerHTML = '';
  const img = document.createElement('img');
  img.src = url;
  img.alt = 'фото';
  container.appendChild(img);
}

async function getRandom() {
  const res = await fetch('/photo/random');
  if (!res.ok) {
    showError('Ошибка сервера: ' + res.status);
    return;
  }
  const data = await res.json();
  showPhoto(data.urls.full);
}

async function searchPhoto() {
  const q = document.getElementById('query').value.trim();
  if (!q) {
    showError('Поле поиска не должно быть пустым');
    return;
  }
  const res = await fetch('/photo/search?q=' + encodeURIComponent(q));
  if (!res.ok) {
    showError('Ошибка сервера: ' + res.status);
    return;
  }
  const data = await res.json();
  showPhoto(data.urls.full);
}

document.getElementById('btn-random').addEventListener('click', getRandom);
document.getElementById('btn-search').addEventListener('click', searchPhoto);

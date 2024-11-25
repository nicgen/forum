// Waits for the DOM to load.
// document.addEventListener('DOMContentLoaded', function () {

const themeToggle = document.getElementById('theme-toggle');
const root = document.documentElement;
const themes = ['light', 'dark', 'auto'];
let currentThemeIndex = 2; // Start with auto

// function setTheme(theme) {
//     root.classList.remove('light-theme', 'dark-theme');
//     if (theme === 'auto') {
//         root.classList.remove('light-theme', 'dark-theme');
//     } else {
//         root.classList.add(`${theme}-theme`);
//     }
//     localStorage.setItem('theme', theme);
// }

function setTheme(theme) {
  root.classList.remove('light-theme', 'dark-theme');
  themeToggle.classList.remove('light', 'dark', 'auto');
  
  if (theme === 'auto') {
      root.classList.remove('light-theme', 'dark-theme');
  } else {
      root.classList.add(`${theme}-theme`);
  }
  
  themeToggle.classList.add(theme);
  localStorage.setItem('theme', theme);
}

function rotateTheme() {
    currentThemeIndex = (currentThemeIndex + 1) % themes.length;
    setTheme(themes[currentThemeIndex]);
}

themeToggle.addEventListener('click', rotateTheme);

// Set initial theme based on user's preference or system setting
const savedTheme = localStorage.getItem('theme');
if (savedTheme) {
    setTheme(savedTheme);
    currentThemeIndex = themes.indexOf(savedTheme);
} else {
    setTheme('auto');
}

//fonction for the expandable post
function togglePost(element) {
    const fullPost = element.querySelector('.post-full');
    fullPost.classList.toggle('active');
}

function closePost(event) {
    event.stopPropagation();
    const fullPost = event.target.closest('.post-full');
    fullPost.classList.remove('active');
}
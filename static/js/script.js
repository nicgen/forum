// Waits for the DOM to load.
// document.addEventListener('DOMContentLoaded', function () {

const themeToggle = document.getElementById('theme-toggle');
const root = document.documentElement;
const themes = ['light', 'dark', 'auto'];
let currentThemeIndex = 2; // Start with auto

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

// POPUP

function togglePost(button) {
    const post = button.closest('.showhide');
    const fullPost = post.querySelector('.js__toggle__block');
    fullPost.classList.toggle('active');
}

// function closePost(event) {
//     event.stopPropagation(); // Prevent the event from bubbling up
//     const fullPost = event.target.closest('.js__toggle__block');
//     fullPost.classList.remove('active');
// }

function closePost(event) {
    event.stopPropagation(); // Prevent the event from bubbling up
    event.preventDefault(); // Prevent any default action (like form submission)

    const fullPost = event.target.closest('.js__toggle__block');
    fullPost.classList.remove('active'); // Close the popup
}

// function closePost(event) {
//     event.stopPropagation(); // Prevent the event from bubbling up
//     const fullPost = event.target.closest('.js__toggle__block');
//     const form = fullPost.querySelector('form'); // Get the form inside the popup

//     // Check if the form is valid
//     if (form.checkValidity()) {
//         fullPost.classList.remove('active'); // Close the popup if the form is valid
//     } else {
//         // Optionally, you can focus on the first invalid input
//         const firstInvalidInput = form.querySelector(':invalid');
//         if (firstInvalidInput) {
//             firstInvalidInput.focus(); // Focus on the first invalid input
//         }
//     }
// }
// Waits for the DOM to load.
// document.addEventListener('DOMContentLoaded', function () {

const themeToggle = document.getElementById('theme-toggle');
const root = document.documentElement;
const themes = ['light', 'dark', 'auto'];
let currentThemeIndex = 2; // Start with auto

// Function to save the scroll position
function saveScrollPosition() {
    localStorage.setItem('scrollPosition', window.scrollY);
}

// Function to restore the scroll position
function restoreScrollPosition() {
    const scrollPosition = localStorage.getItem('scrollPosition');
    if (scrollPosition) {
        window.scrollTo(0, parseInt(scrollPosition, 10));
    }
}

// Event listener to save the scroll position before the page unloads
window.addEventListener('beforeunload', saveScrollPosition);

// Event listener to restore the scroll position when the page loads
window.addEventListener('load', restoreScrollPosition);

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

function checkImageSize(event) {
    const file = event.target.files[0];
    const form = document.getElementById('imageForm');

    if (file) {
        const sizeInMB = file.size;
        if (sizeInMB > 20971520){
            alert("Image size: > 20Mo");
            form.reset(); // Reset the form
        }
        // return false; // Prevent form submission
    }
}

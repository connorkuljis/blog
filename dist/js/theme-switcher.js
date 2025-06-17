const themeToggle = document.querySelector('.theme-toggle');
const body = document.body;

// Check for saved theme in localStorage
const currentTheme = localStorage.getItem('theme');
if (currentTheme) {
	body.setAttribute('data-theme', currentTheme);
	updateButtonIcon(currentTheme);
} else {
	// Check system preference
	const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
	if (prefersDark) {
		body.setAttribute('data-theme', 'dark');
		updateButtonIcon('dark');
	}
}

themeToggle.addEventListener('click', () => {
	const isDark = body.getAttribute('data-theme') === 'dark';
	body.setAttribute('data-theme', isDark ? 'light' : 'dark');
	localStorage.setItem('theme', isDark ? 'light' : 'dark');
	updateButtonIcon(isDark ? 'light' : 'dark');
});

function updateButtonIcon(theme) {
	themeToggle.textContent = theme === 'dark' ? '🌙' : '☀️';
	themeToggle.setAttribute('aria-label', 
		theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme');
}

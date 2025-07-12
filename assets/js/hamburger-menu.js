document.addEventListener('DOMContentLoaded', function() {
    const hamburger = document.querySelector('.hamburger');
    const categoriesList = document.getElementById('categories-list');

    if (hamburger && categoriesList) {
        hamburger.addEventListener('click', function() {
            const isExpanded = this.getAttribute('aria-expanded') === 'true';
            this.setAttribute('aria-expanded', !isExpanded);
            categoriesList.setAttribute('aria-expanded', !isExpanded);
        });
    }
});
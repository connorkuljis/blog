document.addEventListener('DOMContentLoaded', () => {
	const container = document.querySelector('.ascii-art-wrapper');
	const item = document.querySelector('.ascii-art');

	if (!container || !item) {
		console.error("Marquee container or item not found.");
		return;
	}

	// Duplicate the item to create a seamless loop
	const itemClone = item.cloneNode(true);
	container.appendChild(itemClone);

	let position = 0;
	const speed = 0.5; // Pixels per frame (adjust for speed)

	// Ensure the container has enough width to fit both items initially
	// Or set flex-wrap: nowrap on container and let items define width

	function animateMarquee() {
		// Get the actual width of a single item
		// Use getBoundingClientRect().width for accurate width including padding/border
		const itemWidth = item.getBoundingClientRect().width;

		position -= speed; // Move left

		// If the first item has scrolled completely out of view, reset position
		// This makes the clone appear to be the original continuing
		if (position <= -itemWidth) {
			position = 0; // Reset to the starting point
		}

		item.style.transform = `translateX(${position}px)`;
		itemClone.style.transform = `translateX(${position}px)`;

		requestAnimationFrame(animateMarquee);
	}

	animateMarquee();
});

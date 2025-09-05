class TocComponent extends HTMLElement {
  connectedCallback() {
    // Wait for DOM to be fully loaded
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => this.render());
    } else {
      this.render();
    }
  }

  render() {
    const markdown = document.querySelector('.markdown');
    if (!markdown) return;

    const headings = markdown.querySelectorAll('h2, h3, h4, h5, h6');
    if (headings.length === 0) return;

    // Generate IDs for headings if they don't have them
    headings.forEach((heading, index) => {
      if (!heading.id) {
        const text = heading.textContent.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-');
        heading.id = `heading-${index}-${text}`;
      }
    });

    const tocList = this.buildTocList(headings);
    if (!tocList) return;

    this.innerHTML = `
      <details class="toc" open>
        <summary>Table of Contents</summary>
        <nav>
          <ul>
            ${tocList}
          </ul>
        </nav>
      </details>
    `;

    // Add smooth scroll behavior
    this.addEventListener('click', (e) => {
      if (e.target.tagName === 'A') {
        e.preventDefault();
        const targetId = e.target.getAttribute('href').substring(1);
        const target = document.getElementById(targetId);
        if (target) {
          target.scrollIntoView({ behavior: 'smooth' });
        }
      }
    });
  }

  buildTocList(headings) {
    if (headings.length === 0) return '';

    const stack = [];
    let result = '';

    headings.forEach((heading) => {
      const level = parseInt(heading.tagName.charAt(1));
      const text = heading.textContent.trim();
      const id = heading.id;

      const item = `<li><a href="#${id}">${text}</a></li>`;

      // Pop stack until we find the right parent level
      while (stack.length > 0 && stack[stack.length - 1].level >= level) {
        result += '</ul></li>';
        stack.pop();
      }

      // If we need to go deeper, start a new sublist
      if (stack.length > 0 && stack[stack.length - 1].level < level) {
        result += '<ul>';
      }

      result += item;
      stack.push({ level, text, id });
    });

    // Close remaining open lists
    while (stack.length > 0) {
      result += '</ul></li>';
      stack.pop();
    }

    return result;
  }
}

if (!customElements.get('toc-component')) {
  customElements.define('toc-component', TocComponent);
}
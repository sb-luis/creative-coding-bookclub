document.addEventListener('DOMContentLoaded', () => {
  const sidebarToggleCheckbox = document.getElementById('sidebar-toggle-checkbox');
  const sidebarToggle = document.getElementById('sidebar-toggle');

  if (sidebarToggleCheckbox && sidebarToggle) {
    sidebarToggleCheckbox.addEventListener('change', () => {
      const isOpen = sidebarToggleCheckbox.checked;
      sidebarToggle.setAttribute('aria-expanded', isOpen);
    });

    // Optional: Close sidebar with Escape key
    document.addEventListener('keydown', (event) => {
      if (event.key === 'Escape' && sidebarToggleCheckbox.checked) {
        sidebarToggleCheckbox.checked = false;
        sidebarToggle.setAttribute('aria-expanded', 'false');
      }
    });
  }
});
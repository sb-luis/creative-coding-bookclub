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

// Sign out function
async function signOut() {
  try {
    const response = await fetch('/api/auth/sign-out', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (response.ok) {
      const result = await response.json();
      if (result.success) {
        // Redirect to homepage 
        window.location.href = '/';
      } else {
        console.error('Sign out failed');
        alert('Failed to sign out. Please try again.');
      }
    } else {
      console.error('Sign out request failed:', response.status);
      alert('Failed to sign out. Please try again.');
    }
  } catch (error) {
    console.error('Error during sign out:', error);
    alert('Failed to sign out. Please try again.');
  }
}
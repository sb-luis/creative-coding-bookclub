const sketchListerEl = document.getElementById('sketch-lister');
const sketchViewEl = document.getElementById('sketch-view');
const sketchLinkEl = document.getElementById('sketch-link');
const prevSketchBtn = document.getElementById('prev-sketch');
const nextSketchBtn = document.getElementById('next-sketch');
const cycleViewBtn = document.getElementById('cycle-view');
const playStopBtn = document.getElementById('play-stop-sketch');

let currentSketchIndex = -1;
let sketchElements = [];
let currentViewMode = 'overlay'; // Track current view mode
let isSketchRunning = false; // Track sketch state

// Helper function to get sketch metadata from DOM element
function getSketchData(button) {
  const alias = button.dataset.alias;
  const page = button.dataset.page;

  if (!alias || !page) {
    console.warn(
      `Sketch button for "${button.innerText}" is missing essential data attributes (data-alias or data-page). This sketch cannot be loaded.`
    );
    return null;
  }
  const title = button.dataset.title || page; // Fallback title to page name
  const description = button.dataset.description || ''; // Default to empty string

  return {
    params: { alias, page },
    props: { metadata: { title, description } },
  };
}

function updateNavigationButtons() {
  // Always enable navigation buttons when there are sketches to navigate
  prevSketchBtn.disabled = sketchElements.length <= 1;
  nextSketchBtn.disabled = sketchElements.length <= 1;
}

function updatePlayStopButton() {
  if (isSketchRunning) {
    playStopBtn.innerHTML = 'stop';
    playStopBtn.title = 'Stop sketch (Ctrl+.)';
  } else {
    playStopBtn.innerHTML = 'play';
    playStopBtn.title = 'Run sketch (Ctrl+Enter)';
  }
}

function loadSketchByIndex(index) {
  if (index < 0 || index >= sketchElements.length) {
    console.warn(`Invalid sketch index: ${index}`);
    return;
  }

  currentSketchIndex = index;
  const button = sketchElements[index];
  const sketchData = getSketchData(button);

  if (sketchData) {
    // Reset sketch running state when loading new sketch
    isSketchRunning = false;
    updatePlayStopButton();

    loadSketch(sketchData);
    updateNavigationButtons();

    // Update visual indicator and ARIA attributes in sidebar
    sketchElements.forEach((btn) => {
      btn.classList.remove('ccb-active');
      btn.setAttribute('aria-current', 'false');
      btn.setAttribute('aria-selected', 'false');
    });
    button.classList.add('ccb-active');
    button.setAttribute('aria-current', 'page');
    button.setAttribute('aria-selected', 'true');
  }
}

function loadSketch(s) {
  console.log(`Loading sketch: ${s.params.alias}/${s.params.page}`);

  // Include current view mode as query parameter
  const viewModeParam =
    currentViewMode !== 'overlay' ? `?viewMode=${currentViewMode}` : '';
  const sketchEditorUrl = `/sketches/${s.params.alias}/${s.params.page}/edit${viewModeParam}`;
  const sketchCleanUrl = `/sketches/${s.params.alias}/${s.params.page}`;

  console.log('🎯 Loading sketch with view mode:', currentViewMode);
  console.log('📤 Sketch editor URL:', sketchEditorUrl);

  // Load editor view in iframe for inline editing
  sketchViewEl.setAttribute('src', sketchEditorUrl);
  // Set clean view for "open in new tab" functionality
  sketchLinkEl.setAttribute('href', sketchCleanUrl);

  // Collapse sidebar when loading a sketch
  const sidebarToggleCheckbox = document.getElementById(
    'sidebar-toggle-checkbox'
  );
  if (sidebarToggleCheckbox && sidebarToggleCheckbox.checked) {
    sidebarToggleCheckbox.checked = false;
    const sidebarToggle = document.getElementById('sidebar-toggle');
    if (sidebarToggle) {
      sidebarToggle.setAttribute('aria-expanded', 'false');
    }
  }
}

// Navigation button event listeners
prevSketchBtn.addEventListener('click', () => {
  if (sketchElements.length > 1) {
    const newIndex =
      currentSketchIndex > 0
        ? currentSketchIndex - 1
        : sketchElements.length - 1;
    loadSketchByIndex(newIndex);
  }
});

nextSketchBtn.addEventListener('click', () => {
  if (sketchElements.length > 1) {
    const newIndex =
      currentSketchIndex < sketchElements.length - 1
        ? currentSketchIndex + 1
        : 0;
    loadSketchByIndex(newIndex);
  }
});

// View cycle button event listener
cycleViewBtn.addEventListener('click', () => {
  // Send message to the iframe to cycle view mode
  if (sketchViewEl && sketchViewEl.contentWindow) {
    sketchViewEl.contentWindow.postMessage(
      {
        type: 'cycleViewMode',
      },
      '*'
    );
  }
});

// Play/Stop button event listener
playStopBtn.addEventListener('click', () => {
  if (sketchViewEl && sketchViewEl.contentWindow) {
    if (isSketchRunning) {
      // Stop the sketch (equivalent to Ctrl+.)
      sketchViewEl.contentWindow.postMessage(
        {
          type: 'keyboardShortcut',
          shortcut: 'ctrl+period',
        },
        '*'
      );
    } else {
      // Run the sketch (equivalent to Ctrl+Enter)
      sketchViewEl.contentWindow.postMessage(
        {
          type: 'keyboardShortcut',
          shortcut: 'ctrl+enter',
        },
        '*'
      );
    }
  }
});

// Listen for view mode changes from the iframe
window.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'viewModeChanged') {
    let newViewMode = event.data.viewMode;
    console.log('📥 Raw view mode from iframe:', newViewMode);

    // Clean up view mode (remove quotes if present)
    if (typeof newViewMode === 'string') {
      newViewMode = newViewMode.replace(/^"|"$/g, '');
    }

    currentViewMode = newViewMode;
    console.log('📥 View mode updated to:', currentViewMode);
  } else if (event.data && event.data.type === 'sketchRunning') {
    isSketchRunning = true;
    updatePlayStopButton();
  } else if (event.data && event.data.type === 'sketchStopped') {
    isSketchRunning = false;
    updatePlayStopButton();
  }
});

// Event delegation for sketch lister clicks
sketchListerEl.addEventListener('click', (event) => {
  const button = event.target.closest('button.ccb-link');
  if (button && sketchListerEl.contains(button)) {
    const index = sketchElements.indexOf(button);
    if (index !== -1) {
      loadSketchByIndex(index);
    }
  }
});

// Initialize sketch elements and load random sketch
sketchElements = Array.from(
  sketchListerEl.querySelectorAll('button.ccb-link[data-alias][data-page]')
);

if (sketchElements.length > 0) {
  const randomIndex = Math.floor(Math.random() * sketchElements.length);
  loadSketchByIndex(randomIndex);
} else {
  console.log('No sketches with required data-alias and data-page found.');
}

// Initialize play/stop button
updatePlayStopButton();
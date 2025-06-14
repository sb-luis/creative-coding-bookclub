import { elements, state, VIEW_MODES } from './dom-elements.js';
import { hideConsole, showConsole } from './console-manager.js';
import { updateLineNumbers, highlightCurrentLine } from './line-numbers.js';
import { createAndRunSketch } from './iframe-manager.js';

export function updateVisibility() {
  // Save code editor state before hiding
  if (!elements.codeEditorContainer.classList.contains('hidden')) {
    state.savedScrollTop = elements.codeEditor.scrollTop;
    state.savedSelectionStart = elements.codeEditor.selectionStart;
    state.savedSelectionEnd = elements.codeEditor.selectionEnd;
  }

  // Update visibility based on current view mode
  switch (state.currentViewMode) {
    case VIEW_MODES.CODE:
      elements.codeEditorContainer.classList.remove('hidden');
      elements.sketchIframeContainer.classList.add('hidden');
      hideConsole();
      break;
    case VIEW_MODES.SKETCH:
      elements.codeEditorContainer.classList.add('hidden');
      elements.sketchIframeContainer.classList.remove('hidden');
      hideConsole();
      break;
    case VIEW_MODES.OVERLAY:
      // Show sketch with code editor overlay
      elements.codeEditorContainer.classList.remove('hidden');
      elements.sketchIframeContainer.classList.remove('hidden');
      hideConsole();
      break;
    case VIEW_MODES.DEBUG:
      // Show sketch with console overlay
      elements.codeEditorContainer.classList.add('hidden');
      elements.sketchIframeContainer.classList.remove('hidden');
      showConsole();
      break;
  }

  // Restore code editor previous state if applicable
  if (
    (state.currentViewMode === VIEW_MODES.CODE ||
      state.currentViewMode === VIEW_MODES.OVERLAY) &&
    !elements.codeEditorContainer.classList.contains('just-restored')
  ) {
    setTimeout(() => {
      if (!elements.codeEditorContainer.classList.contains('hidden')) {
        elements.codeEditor.scrollTop = state.savedScrollTop;
        elements.codeEditor.setSelectionRange(
          state.savedSelectionStart,
          state.savedSelectionEnd
        );
        elements.codeEditor.focus();
        updateLineNumbers();
        highlightCurrentLine();
        // Mark as restored to prevent unnecessary re-application
        elements.codeEditorContainer.classList.add('just-restored');
      }
    }, 50);
  }
}

// Cycle through view modes: overlay -> sketch -> debug -> overlay
export function cycleViewMode() {
  elements.codeEditorContainer.classList.remove('just-restored');

  console.log('🔄 View mode before cycling:', state.currentViewMode);

  switch (state.currentViewMode) {
    case VIEW_MODES.OVERLAY:
      state.currentViewMode = VIEW_MODES.SKETCH;
      console.log('🔄 Switched to Sketch view');
      break;
    case VIEW_MODES.SKETCH:
      state.currentViewMode = VIEW_MODES.DEBUG;
      console.log('🔄 Switched to Debug view');
      break;
    case VIEW_MODES.DEBUG:
      state.currentViewMode = VIEW_MODES.OVERLAY;
      console.log('🔄 Switched to Overlay view');
      break;
    default:
      // Fallback to overlay if somehow we're in an invalid state
      state.currentViewMode = VIEW_MODES.OVERLAY;
      console.log('🔄 Reset to Overlay view (fallback)');
      break;
  }

  console.log('🔄 New view mode after cycling:', state.currentViewMode);
  updateVisibility();

  // Create sketch if needed
  if (!state.sketchIframe && elements.codeEditor.value.trim() !== '') {
    setTimeout(createAndRunSketch, 0);
  }
}

export function getCurrentViewMode() {
  return state.currentViewMode;
}

// Set a specific view mode
export function setViewMode(mode) {
  if (Object.values(VIEW_MODES).includes(mode)) {
    elements.codeEditorContainer.classList.remove('just-restored');
    state.currentViewMode = mode;
    updateVisibility();
  }
}

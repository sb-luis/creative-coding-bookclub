import { elements, state } from './dom-elements.js';
import { formatCode } from './code-formatter.js';
import { createAndRunSketch, stopSketch } from './iframe-manager.js';
import { clearConsole } from './console-manager.js';
import { toggleComment } from './code-formatter.js';
import {
  getCurrentViewMode,
  setViewMode,
  cycleViewMode,
} from './view-manager.js';

export function setupKeyboardShortcuts() {
  // Add keyboard shortcuts to the document
  document.addEventListener('keydown', function (event) {
    if (event.ctrlKey && event.key === 'Enter') {
      console.log('⌨️ Ctrl+Enter pressed - running sketch');
      event.preventDefault();
      createAndRunSketch();
    } else if (event.ctrlKey && event.key === '.') {
      event.preventDefault();
      stopSketch();
      clearConsole();
    } else if (event.ctrlKey && event.key === 'f') {
      event.preventDefault();
      formatCode();
    } else if (event.ctrlKey && event.key === '/') {
      event.preventDefault();
      toggleComment();
    } else if (event.ctrlKey && event.key === 's') {
      // Forward Ctrl+S to parent window (sketch manager) if we're in an iframe
      if (window.parent && window.parent !== window) {
        event.preventDefault();
        window.parent.postMessage(
          {
            type: 'keyboardShortcut',
            shortcut: 'ctrl+s',
          },
          '*'
        );
      }
    }
  });

  // VIEW TOGGLING SHORTCUTS
  document.addEventListener('keydown', function (event) {
    // Ctrl+, prioritizes overlay mode first, then toggles to sketch
    if (event.ctrlKey && event.key === ',') {
      event.preventDefault();
      const currentMode = getCurrentViewMode();
      if (currentMode === 'overlay') {
        // If already in overlay, go to sketch-only
        setViewMode('sketch');
      } else {
        setViewMode('overlay');
      }
      // Ctrl+; prioritizes debug mode first, then toggles to overlay
    } else if (
      event.ctrlKey &&
      (event.key === ';' || event.code === 'Semicolon')
    ) {
      event.preventDefault();
      const currentMode = getCurrentViewMode();
      if (currentMode === 'debug') {
        // If already in debug, go to overlay
        setViewMode('overlay');
      } else {
        setViewMode('debug');
      }
    }
  });
}

export function setupMessageHandling() {
  window.addEventListener('message', (event) => {
    if (
      event.source === window.parent &&
      event.data &&
      event.data.type === 'runSketch'
    ) {
      setViewMode('overlay');
      createAndRunSketch();
    } else if (
      event.source === window.parent &&
      event.data &&
      event.data.type === 'stopSketch'
    ) {
      setViewMode('code');
      stopSketch();
    } else if (
      event.source === window.parent &&
      event.data &&
      event.data.type === 'cycleViewMode'
    ) {
      cycleViewMode();
    } else if (
      event.source === window.parent &&
      event.data &&
      event.data.type === 'sketchSaved'
    ) {
      state.isDirty = false;
      console.log('✅ Sketch saved notification received from parent');
    } else if (
      event.source === window.parent &&
      event.data &&
      event.data.type === 'sketchLoaded'
    ) {
      state.isDirty = false;
      console.log('✅ Sketch loaded notification received from parent');
    } else if (
      event.data &&
      event.data.type === 'keyboardShortcut'
    ) {
      // Handle keyboard shortcuts forwarded from iframe
      switch (event.data.shortcut) {
        case 'ctrl+comma':
          const currentMode = getCurrentViewMode();
          if (currentMode === 'sketch') {
            setViewMode('overlay');
          } else {
            setViewMode('sketch');
          }
          break;
        case 'ctrl+semicolon':
          const currentMode2 = getCurrentViewMode();
          if (currentMode2 === 'debug') {
            setViewMode('overlay');
          } else {
            setViewMode('debug');
          }
          break;
        case 'ctrl+enter':
          createAndRunSketch();
          break;
        case 'ctrl+period':
          stopSketch();
          clearConsole();
          break;
        case 'ctrl+s':
          if (window.parent && window.parent !== window) {
            window.parent.postMessage({
              type: 'keyboardShortcut',
              shortcut: 'ctrl+s'
            }, '*');
          }
          break;
      }
    }
  });
}

export function setupParentCommunication() {
  if (window.parent && window.parent !== window) {
    elements.codeEditor.addEventListener('input', function () {
      state.isDirty = true;
      window.parent.postMessage({ type: 'sketchDirty', status: true }, '*');
    });
  }
}

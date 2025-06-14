// Console overlay management
import { elements } from './dom-elements.js';

export function clearConsole() {
  if (elements.consoleOutput) {
    elements.consoleOutput.innerHTML = '';
  }
}

export function showConsole() {
  if (elements.consoleOverlayContainer) {
    elements.consoleOverlayContainer.classList.remove('hidden');
  }
}

export function hideConsole() {
  if (elements.consoleOverlayContainer) {
    elements.consoleOverlayContainer.classList.add('hidden');
  }
}

export function toggleConsole() {
  if (elements.consoleOverlayContainer) {
    if (elements.consoleOverlayContainer.classList.contains('hidden')) {
      showConsole();
    } else {
      hideConsole();
    }
  }
}

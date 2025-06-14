import { elements } from './dom-elements.js';
import { updateLineNumbers, highlightCurrentLine } from './line-numbers.js';

// Functions to update cursor position and file size
export function updateCursorPosition() {
  const text = elements.codeEditor.value;
  const cursorPos = elements.codeEditor.selectionStart;

  // Count lines and columns
  let lineNumber = 1;
  let colNumber = 1;

  for (let i = 0; i < cursorPos; i++) {
    if (text[i] === '\n') {
      lineNumber++;
      colNumber = 1;
    } else {
      colNumber++;
    }
  }

  elements.cursorPositionEl.textContent = `Ln ${lineNumber}, Col ${colNumber}`;
}

export function formatFileSize(bytes) {
  if (bytes < 1024) {
    return bytes + ' bytes';
  } else if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(1) + ' KB';
  } else {
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }
}

export function updateFileSize() {
  const sizeInBytes = new TextEncoder().encode(elements.codeEditor.value).length;
  elements.fileSizeEl.textContent = formatFileSize(sizeInBytes);
}

export function setupStatusTracking() {
  elements.codeEditor.addEventListener('input', function () {
    updateCursorPosition();
    updateFileSize();
    updateLineNumbers();
  });

  elements.codeEditor.addEventListener('click', function () {
    updateCursorPosition();
    highlightCurrentLine();
  });

  elements.codeEditor.addEventListener('keyup', function () {
    updateCursorPosition();
    highlightCurrentLine();
  });

  elements.codeEditor.addEventListener('mouseup', function () {
    updateCursorPosition();
    highlightCurrentLine();
  });

  elements.codeEditor.addEventListener('focus', function () {
    updateCursorPosition();
    updateFileSize();
    updateLineNumbers();
  });
}

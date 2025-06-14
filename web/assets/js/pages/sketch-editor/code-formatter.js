import { elements, state } from './dom-elements.js';
import { updateLineNumbers } from './line-numbers.js';
import { updateFileSize, updateCursorPosition } from './status-tracker.js';

export function formatCode() {
  console.log('🛠️ Formatting code...');

  const code = elements.codeEditor.value;
  if (code.trim() === '') return;

  // Convert tabs to spaces and normalize line endings 
  let normalized = code.split('\r\n').join('\n'); // Normalize line endings
  normalized = normalized.split('\t').join('  '); // Convert tabs to spaces

  // Split into lines and process
  const lines = normalized.split('\n');
  let indentLevel = 0;
  const indentSize = 2;
  const formattedLines = [];

  for (let i = 0; i < lines.length; i++) {
    const originalLine = lines[i];
    const trimmedLine = originalLine.trim();

    // Skip empty lines for now
    if (trimmedLine === '') {
      formattedLines.push('');
      continue;
    }

    // Decrease indent level for closing braces at the beginning of the line
    if (trimmedLine.startsWith('}')) {
      indentLevel = Math.max(0, indentLevel - 1);
    }

    // Apply indentation 
    formattedLines.push(' '.repeat(indentLevel * indentSize) + trimmedLine);

    // Increase indent level for opening braces at the end of the line
    if (trimmedLine.endsWith('{')) {
      indentLevel++;
    }
  }

  // Remove redundant blank lines (more than one consecutive blank line)
  const finalLines = [];
  let consecutiveBlankLines = 0;

  for (const line of formattedLines) {
    if (line.trim() === '') {
      consecutiveBlankLines++;
      if (consecutiveBlankLines <= 1) {
        finalLines.push(line);
      }
    } else {
      consecutiveBlankLines = 0;
      finalLines.push(line);
    }
  }

  const originalScrollTop = elements.codeEditor.scrollTop;
  const originalSelectionStart = elements.codeEditor.selectionStart;
  const originalSelectionEnd = elements.codeEditor.selectionEnd;

  elements.codeEditor.value = finalLines.join('\n');

  // Restore scroll position and selection 
  elements.codeEditor.scrollTop = originalScrollTop;
  elements.codeEditor.setSelectionRange(originalSelectionStart, originalSelectionEnd);

  // Update UI components
  updateLineNumbers();
  updateFileSize();
  updateCursorPosition();
}

export function clearCode() {
  if (
    confirm(
      'Are you sure you want to clear all code? This action cannot be undone.'
    )
  ) {
    elements.codeEditor.value = '';
    elements.codeEditor.focus();
    
    // Update UI components
    updateLineNumbers();
    updateFileSize();
    updateCursorPosition();

    // Stop any running sketch
    if (state.sketchIframe) {
      state.sketchIframe.remove();
      state.sketchIframe = null;
    }
  }
}

export function toggleComment() {
  const start = elements.codeEditor.selectionStart;
  const end = elements.codeEditor.selectionEnd;
  const text = elements.codeEditor.value;

  // If no selection, comment/uncomment the current line
  let lineStart, lineEnd;
  if (start === end) {
    // Find the start and end of the current line
    lineStart = text.lastIndexOf('\n', start - 1) + 1;
    lineEnd = text.indexOf('\n', start);
    if (lineEnd === -1) lineEnd = text.length;
  } else {
    // Find the start of the first selected line and end of the last selected line
    lineStart = text.lastIndexOf('\n', start - 1) + 1;
    lineEnd = text.indexOf('\n', end - 1);
    if (lineEnd === -1) lineEnd = text.length;
  }

  // Extract the lines to work with
  const selectedText = text.substring(lineStart, lineEnd);
  const lines = selectedText.split('\n');

  // Check if all non-empty lines are commented
  const nonEmptyLines = lines.filter((line) => line.trim().length > 0);
  const allCommented =
    nonEmptyLines.length > 0 &&
    nonEmptyLines.every((line) => line.trim().startsWith('//'));

  let newLines;
  if (allCommented) {
    // Uncomment: remove the first occurrence of // and trailing spaces 
    newLines = lines.map((line) => {
      if (line.trim().startsWith('//')) {
        return line.replace(/^(\s*)\/\/\s?/, '$1');
      }
      return line;
    });
  } else {
    // Comment: add // to the beginning of each non-empty line
    newLines = lines.map((line) => {
      if (line.trim().length > 0) {
        const indent = line.match(/^\s*/)[0];
        const content = line.substring(indent.length);
        return indent + '// ' + content;
      }
      return line;
    });
  }

  // Replace the text
  const newSelectedText = newLines.join('\n');
  const newText =
    text.substring(0, lineStart) + newSelectedText + text.substring(lineEnd);

  elements.codeEditor.value = newText;

  // Restore selection
  const selectionOffset = newSelectedText.length - selectedText.length;
  elements.codeEditor.setSelectionRange(start, end + selectionOffset);
  elements.codeEditor.focus();

  // Update UI components
  updateLineNumbers();
  updateFileSize();
  updateCursorPosition();
}

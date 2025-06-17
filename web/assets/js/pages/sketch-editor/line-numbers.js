import { elements, state } from './dom-elements.js';

// Estimate line width for monospace fonts
function lineLength(line, charWidth) {
  return line.length * charWidth;
}

let charWidth = null;

function calculateCharWidth() {
  const testChar = 'M'; 
  const span = document.createElement('span');
  const styles = getComputedStyle(elements.codeEditor);
  
  span.style.font = styles.font;
  span.style.visibility = 'hidden';
  span.style.position = 'absolute';
  span.textContent = testChar;
  
  document.body.appendChild(span);
  const width = span.offsetWidth;
  document.body.removeChild(span);
  
  return width;
}

function getCharWidth() {
  if (charWidth === null) {
    charWidth = calculateCharWidth();
  }
  return charWidth;
}

// Helper function to check if a line is a comment
function isCommentLine(lineText) {
  const trimmed = lineText.trim();
  return trimmed.startsWith('//');
}

// Helper function to check if a line is blank
function isBlankLine(lineText) {
  return lineText.trim() === '';
}

export function updateLineNumbers() {
  const text = elements.codeEditor.value;
  const lines = text.split('\n');
  const lineCount = lines.length;

  // Only rebuild if line count has changed
  const existingLineNumbers = elements.lineNumbersContainer.querySelectorAll('.line-number');
  const needsRebuild = existingLineNumbers.length !== lineCount;

  // Get textarea styles and dimensions for width calculations
  const textareaStyles = getComputedStyle(elements.codeEditor);
  const paddingLeft = parseFloat(textareaStyles.paddingLeft); 
  const paddingRight = parseFloat(textareaStyles.paddingRight);
  const textareaWidth = elements.codeEditor.clientWidth;
  const availableTextWidth = textareaWidth - paddingLeft - paddingRight;

  if (needsRebuild) {
    // Clear the line numbers container
    elements.lineNumbersContainer.innerHTML = '';

    // Add a line number for each logical line in the textarea
    for (let i = 1; i <= lineCount; i++) {
      const lineText = lines[i - 1] || '';

      // Create line number element
      const lineElement = document.createElement('div');
      lineElement.textContent = i;
      lineElement.id = `line-${i}`;
      lineElement.classList.add('line-number');
      
      // Add line type classes
      if (isBlankLine(lineText)) {
        lineElement.classList.add('blank-line');
      } else if (isCommentLine(lineText)) {
        lineElement.classList.add('comment-line');
      }

      const textWidth = lineLength(lineText, getCharWidth());

      // Element should extend to cover line numbers (45px) + text + small padding
      const elementWidth = Math.max(70, 45 + Math.min(textWidth, availableTextWidth) + 5);
      lineElement.style.width = `${elementWidth}px`;

      // Append to container
      elements.lineNumbersContainer.appendChild(lineElement);
    }
  } else {
    // Update line dimensions and classes when line count hasn't changed
    const lineElements = elements.lineNumbersContainer.querySelectorAll('.line-number');
    
    for (let i = 0; i < lineCount; i++) {
      const lineText = lines[i] || '';
      const lineElement = lineElements[i];

      if (lineElement) {
        const textWidth = lineLength(lineText, getCharWidth());

        // Element should extend to cover line numbers (45px) + text + small padding
        const elementWidth = Math.max(70, 45 + Math.min(textWidth, availableTextWidth) + 5);
        lineElement.style.width = `${elementWidth}px`;
        
        // Update line type classes
        lineElement.classList.remove('comment-line', 'blank-line');
        if (isBlankLine(lineText)) {
          lineElement.classList.add('blank-line');
        } else if (isCommentLine(lineText)) {
          lineElement.classList.add('comment-line');
        }
      }
    }
  }

  // Update highlighting for current line
  const lineElements = elements.lineNumbersContainer.querySelectorAll('.line-number');
  
  lineElements.forEach((lineElement, index) => {
    const lineText = lines[index] || '';
    
    // First, update line type classes
    lineElement.classList.remove('comment-line', 'blank-line');
    if (isBlankLine(lineText)) {
      lineElement.classList.add('blank-line');
    } else if (isCommentLine(lineText)) {
      lineElement.classList.add('comment-line');
    }
    
    // Then, update current line highlighting
    if (index + 1 === state.currentHighlightedLine) {
      lineElement.classList.add('current-line');
    } else {
      lineElement.classList.remove('current-line');
    }
  });

  // Sync scroll position
  elements.lineNumbersContainer.scrollTop = elements.codeEditor.scrollTop;
}

export function highlightCurrentLine() {
  const text = elements.codeEditor.value;
  const cursorPos = elements.codeEditor.selectionStart;

  // Count lines to find current logical line number
  let lineNumber = 1;
  for (let i = 0; i < cursorPos; i++) {
    if (text[i] === '\n') {
      lineNumber++;
    }
  }

  // Only update if the line has actually changed
  if (state.currentHighlightedLine !== lineNumber) {
    state.currentHighlightedLine = lineNumber;
    
    // Update highlighting without rebuilding the entire line numbers
    const lineElements = elements.lineNumbersContainer.querySelectorAll('.line-number');
    
    // Get current line text to check if it's a comment
    const lines = text.split('\n');
    
    lineElements.forEach((lineElement, index) => {
      const lineText = lines[index] || '';
      
      // First, update line type classes
      lineElement.classList.remove('comment-line', 'blank-line');
      if (isBlankLine(lineText)) {
        lineElement.classList.add('blank-line');
      } else if (isCommentLine(lineText)) {
        lineElement.classList.add('comment-line');
      }
      
      // Then, update current line highlighting
      if (index + 1 === state.currentHighlightedLine) {
        lineElement.classList.add('current-line');
      } else {
        lineElement.classList.remove('current-line');
      }
    });
  }
}

export function setupLineNumbers() {
  // Sync scrolling between textarea and line numbers
  elements.codeEditor.addEventListener('scroll', function () {
    // Use requestAnimationFrame for smooth scrolling sync
    requestAnimationFrame(() => {
      elements.lineNumbersContainer.scrollTop = elements.codeEditor.scrollTop;
    });
  });

  // Recalculate line numbers when window is resized (affects width calculations)
  let resizeTimeout;
  window.addEventListener('resize', function () {
    // Debounce resize events to avoid excessive recalculation
    clearTimeout(resizeTimeout);
    resizeTimeout = setTimeout(() => {
      updateLineNumbers();
      highlightCurrentLine();
    }, 150);
  });
}

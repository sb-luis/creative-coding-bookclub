import { elements, state } from './dom-elements.js';

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

  // Calculate character width using canvas 
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');
  ctx.font = `${textareaStyles.fontSize} ${textareaStyles.fontFamily}`;

  // Fixed line height (no wrapping)
  const lineHeight = 1.2;

  if (needsRebuild) {
    // Clear the line numbers container
    elements.lineNumbersContainer.innerHTML = '';

    // Add a line number for each logical line in the textarea
    for (let i = 1; i <= lineCount; i++) {
      const lineText = lines[i - 1] || '';

      // Create background element for this line
      const lineBackground = document.createElement('div');
      lineBackground.classList.add('line-background');
      
      // Check line type and add the appropriate class
      if (isBlankLine(lineText)) {
        lineBackground.classList.add('blank-line-bg');
      } else if (isCommentLine(lineText)) {
        lineBackground.classList.add('comment-line-bg');
      }
      
      lineBackground.style.top = `${(i - 1) * lineHeight}rem`;
      lineBackground.style.height = `${lineHeight}rem`;

      // Calculate the width of the text content using canvas
      ctx.font = `${textareaStyles.fontSize} ${textareaStyles.fontFamily}`;
      const textWidth = ctx.measureText(lineText).width;

      // Background should start from 0 and extend to cover line numbers (45px) + text + small padding
      const backgroundWidth = Math.max(70, 45 + Math.min(textWidth, availableTextWidth) + 5);
      lineBackground.style.width = `${backgroundWidth}px`;

      // Create the line number element
      const lineNumberSpan = document.createElement('div');
      lineNumberSpan.textContent = i;
      lineNumberSpan.id = `line-${i}`;
      lineNumberSpan.classList.add('line-number');
      lineNumberSpan.style.top = `${(i - 1) * lineHeight}rem`;
      lineNumberSpan.style.height = `${lineHeight}rem`;
      lineNumberSpan.style.display = 'flex';
      lineNumberSpan.style.alignItems = 'flex-start';
      lineNumberSpan.style.paddingTop = '0';

      // Append background first, then line number
      elements.lineNumbersContainer.appendChild(lineBackground);
      elements.lineNumbersContainer.appendChild(lineNumberSpan);
    }
  } else {
    // Update line background dimensions when line count hasn't changed
    const backgrounds = elements.lineNumbersContainer.querySelectorAll('.line-background');
    const lineNumbers = elements.lineNumbersContainer.querySelectorAll('.line-number');
    
    for (let i = 0; i < lineCount; i++) {
      const lineText = lines[i] || '';
      
      // Calculate the width of the text content using canvas
      ctx.font = `${textareaStyles.fontSize} ${textareaStyles.fontFamily}`;
      const textWidth = ctx.measureText(lineText).width;
      
      // Background should start from 0 and extend to cover line numbers (45px) + text + small padding
      const backgroundWidth = Math.max(70, 45 + Math.min(textWidth, availableTextWidth) + 5);
      
      if (backgrounds[i]) {
        backgrounds[i].style.top = `${i * lineHeight}rem`;
        backgrounds[i].style.height = `${lineHeight}rem`;
        backgrounds[i].style.width = `${backgroundWidth}px`;
        
        // Update line type classes
        backgrounds[i].classList.remove('comment-line-bg', 'blank-line-bg');
        if (isBlankLine(lineText)) {
          backgrounds[i].classList.add('blank-line-bg');
        } else if (isCommentLine(lineText)) {
          backgrounds[i].classList.add('comment-line-bg');
        }
      }

      if (lineNumbers[i]) {
        lineNumbers[i].style.top = `${i * lineHeight}rem`;
        lineNumbers[i].style.height = `${lineHeight}rem`;
        lineNumbers[i].style.display = 'flex';
        lineNumbers[i].style.alignItems = 'flex-start';
        lineNumbers[i].style.paddingTop = '0';
      }
    }
  }

  // Update highlighting for current line
  const backgrounds = elements.lineNumbersContainer.querySelectorAll('.line-background');
  const lineNumbers = elements.lineNumbersContainer.querySelectorAll('.line-number');
  
  backgrounds.forEach((bg, index) => {
    const lineText = lines[index] || '';
    
    // First, update line type classes
    bg.classList.remove('comment-line-bg', 'blank-line-bg');
    if (isBlankLine(lineText)) {
      bg.classList.add('blank-line-bg');
    } else if (isCommentLine(lineText)) {
      bg.classList.add('comment-line-bg');
    }
    
    // Then, update current line highlighting
    if (index + 1 === state.currentHighlightedLine) {
      bg.classList.add('current-line-bg');
    } else {
      bg.classList.remove('current-line-bg');
    }
  });

  lineNumbers.forEach((lineNum, index) => {
    if (index + 1 === state.currentHighlightedLine) {
      lineNum.classList.add('current-line');
    } else {
      lineNum.classList.remove('current-line');
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
    const backgrounds = elements.lineNumbersContainer.querySelectorAll('.line-background');
    const lineNumbers = elements.lineNumbersContainer.querySelectorAll('.line-number');
    
    // Get current line text to check if it's a comment
    const text = elements.codeEditor.value;
    const lines = text.split('\n');
    
    backgrounds.forEach((bg, index) => {
      const lineText = lines[index] || '';
      
      // First, update line type classes
      bg.classList.remove('comment-line-bg', 'blank-line-bg');
      if (isBlankLine(lineText)) {
        bg.classList.add('blank-line-bg');
      } else if (isCommentLine(lineText)) {
        bg.classList.add('comment-line-bg');
      }
      
      // Then, update current line highlighting
      if (index + 1 === state.currentHighlightedLine) {
        bg.classList.add('current-line-bg');
      } else {
        bg.classList.remove('current-line-bg');
      }
    });

    lineNumbers.forEach((lineNum, index) => {
      if (index + 1 === state.currentHighlightedLine) {
        lineNum.classList.add('current-line');
        lineNum.style.fontWeight = 'bold';
      } else {
        lineNum.classList.remove('current-line');
        lineNum.style.fontWeight = 'normal';
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

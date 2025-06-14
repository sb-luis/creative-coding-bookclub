import { elements, state } from './dom-elements.js';
import { fetchOriginalIframeHTML, extractSketchScript } from './iframe-manager.js';
import { updateLineNumbers, highlightCurrentLine } from './line-numbers.js';
import { updateCursorPosition, updateFileSize } from './status-tracker.js';

export async function initializeIframe() {
  const iframe = document.getElementById('sketch-iframe');

  if (iframe) {
    if (iframe.src) {
      const success = await fetchOriginalIframeHTML(iframe.src);
      if (success) {
        console.log('✅ Successfully fetched clean iframe HTML');
      } else {
        console.log('⚠️ Pre-fetch failed, will try again after iframe loads');
      }
    }
    
    iframe.addEventListener('load', function() {
      setTimeout(async () => {
        // Capture the original iframe HTML structure from server (if not already done)
        if (!state?.originalIframeHTML || state.originalIframeHTML.length === 0) {
          console.log('📥 Original HTML not captured yet, fetching now...');
          await fetchOriginalIframeHTML();
        } 
        
        // Extract the sketch script content and populate the code editor
        const sketchCode = await extractSketchScript();
        console.log('📝 Extracted sketch code length:', sketchCode.length);
        
        if (sketchCode && elements.codeEditor.value.trim() === '') {
          elements.codeEditor.value = sketchCode;
        } else if (sketchCode) {
          console.log('⚠️ Code editor already has content, not overwriting');
          console.log('📏 Current editor content length:', elements.codeEditor.value.length);
        } else {
          console.log('❌ No sketch code extracted');
        }
        
        // Initialize UI components after loading content
        updateLineNumbers();
        highlightCurrentLine();
        updateCursorPosition();
        updateFileSize();
        
        console.log('✅ Iframe processing completed');
      }, 100); 
    });
  } else {
    console.error('❌ No iframe element found during initialization');
  }
}

export function initializeUI() {
  elements.statusBar.classList.remove('hidden');
  elements.codeEditorContainer.classList.add('just-restored');
  // Initialize with overlay view 
  elements.codeEditorContainer.classList.remove('hidden');
  elements.sketchIframeContainer.classList.remove('hidden');
  // Ensure console overlay starts hidden
  elements.consoleOverlayContainer.classList.add('hidden');

  // Initial updates
  updateCursorPosition();
  updateFileSize();
  updateLineNumbers();
  highlightCurrentLine();
}

import { elements, state } from './dom-elements.js';
import { fetchOriginalIframeHTML, extractSketchScript } from './iframe-manager.js';
import { updateLineNumbers, highlightCurrentLine } from './line-numbers.js';
import { updateCursorPosition, updateFileSize } from './status-tracker.js';

export async function initializeIframe() {
  const iframe = document.getElementById('sketch-iframe');

  if (iframe) {
    const sketchUrl = iframe.getAttribute('data-sketch-src');
    
    if (sketchUrl) {
      const success = await fetchOriginalIframeHTML(sketchUrl);
      if (success) {
        console.log('✅ Successfully fetched clean iframe HTML');
        
        iframe.src = sketchUrl;
        
        iframe.addEventListener('load', function() {
          setTimeout(async () => {
            // Extract the sketch script content and populate the code editor
            const sketchCode = await extractSketchScript();
            console.log('📝 Extracted sketch code length:', sketchCode?.length || 0);
            
            if (sketchCode && elements.codeEditor.value.trim() === '') {
              elements.codeEditor.value = sketchCode;
            } else if (sketchCode) {
              console.log('⚠️ Code editor already has content, not overwriting');
              console.log('📏 Current editor content length:', elements.codeEditor.value.length);
            } else {
              console.log('❌ No sketch code extracted');
            }
            
            // Remove the iframe src to stop the sketch from running
            iframe.removeAttribute('src');
            iframe.src = 'about:blank';
            
            // Initialize UI components after loading content
            updateLineNumbers();
            highlightCurrentLine();
            updateCursorPosition();
            updateFileSize();
            
            console.log('✅ Iframe processing completed, sketch stopped');
          }, 100); 
        }, { once: true }); 
      } else {
        console.log('⚠️ Failed to fetch iframe HTML');
      }
    }
  } else {
    console.error('❌ No iframe element found during initialization');
  }
}

export function initializeUI() {
  elements.statusBar.classList.remove('hidden');
  elements.codeEditorContainer.classList.add('just-restored');
  // Initialize with only code editor visible 
  elements.codeEditorContainer.classList.remove('hidden');
  // Ensure console overlay starts hidden
  elements.consoleOverlayContainer.classList.add('hidden');

  // Initial updates
  updateCursorPosition();
  updateFileSize();
  updateLineNumbers();
  highlightCurrentLine();
}

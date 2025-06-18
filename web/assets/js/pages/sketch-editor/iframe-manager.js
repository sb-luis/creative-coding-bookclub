import { elements, state } from './dom-elements.js';
import { setViewMode } from './view-manager.js';
import { clearConsole } from './console-manager.js';

// Fetches iframe's HTML from server 
export async function fetchOriginalIframeHTML(iframeUrl) {
  console.log('🌐 fetchOriginalIframeHTML() called with URL:', iframeUrl);

  try {
    console.log('📥 Fetching clean iframe HTML from server...');

    const response = await fetch(iframeUrl);
    if (response.ok) {
      const html = await response.text();
      state.originalIframeHTML = html;
      console.log('✅ Fetched original iframe HTML from server');
      console.log('📏 HTML length:', state.originalIframeHTML.length);
      return true;
    } else {
      console.error(
        '❌ Failed to fetch iframe HTML from server:',
        response.status
      );
      return false;
    }
  } catch (e) {
    console.error('❌ Error fetching iframe HTML from server:', e);
    return false;
  }
}

// Function to extract the sketch script path from the iframe's data attributes
export async function extractSketchScript() {
  const iframe = document.getElementById('sketch-iframe');

  if (iframe && iframe.contentDocument) {
    try {
      const body = iframe.contentDocument.body;
      const sketchJsPath = body.getAttribute('data-sketch-js-path');

      if (sketchJsPath) {
        console.log('📥 Fetching JavaScript content from URL...');

        // Fetch the JavaScript content from the path
        const response = await fetch(sketchJsPath);
        if (response.ok) {
          const content = await response.text();
          console.log('✅ Successfully fetched script content');
          console.log('📏 Content length:', content.length);
          return content;
        } else {
          console.error('❌ Failed to fetch sketch script:', response.status);
        }
      } else {
        console.log('❌ No sketch JS path found in data attributes');
      }
    } catch (e) {
      console.error('❌ Error extracting sketch script:', e);
    }
  } else {
    console.error(
      '❌ Iframe or contentDocument not available for script extraction'
    );
  }
  return '';
}

export async function createAndRunSketch() {
  console.log('🚀 createAndRunSketch() called');

  // Clear the console at the start of every sketch run
  clearConsole();
  console.log('🧹 Console cleared for new sketch run');

  // If we don't have the original HTML captured yet, we can't proceed
  if (!state.originalIframeHTML) {
    console.error('❌ Original iframe HTML not captured yet');

    // Try to capture it now as a last resort
    const iframe = document.getElementById('sketch-iframe');
    const sketchUrl = iframe?.getAttribute('data-sketch-src');
    
    if (sketchUrl) {
      console.log('🚨 Attempting emergency capture of original HTML from:', sketchUrl);
      const success = await fetchOriginalIframeHTML(sketchUrl);
      if (!success) {
        console.error('❌ Emergency capture failed, cannot proceed');
        return;
      }
    } else {
      console.error('❌ No sketch URL available for emergency capture');
      return;
    }
  }

  // Clean up any existing iframes in the container
  if (state.sketchIframe) {
    console.log('🗑️ Removing existing iframe from state');
    state.sketchIframe.remove();
    state.sketchIframe = null;
  }

  // Also remove any other iframes that might be in the container
  const existingIframes =
    elements.sketchIframeContainer.querySelectorAll('iframe');
  if (existingIframes.length > 0) {
    console.log(
      `🧹 Cleaning up ${existingIframes.length} existing iframe(s) from container`
    );
    existingIframes.forEach((iframe) => iframe.remove());
  }

  state.sketchIframe = document.createElement('iframe');
  state.sketchIframe.id = 'sketch-iframe';
  state.sketchIframe.className = 'w-full h-full border-0';
  state.sketchIframe.setAttribute('sandbox', 'allow-scripts allow-same-origin');
  elements.sketchIframeContainer.appendChild(state.sketchIframe);

  const sketchDocument = state.sketchIframe.contentWindow.document;
  const userCode = elements.codeEditor.value;

  // Use the original HTML as base and modify it to load member code
  let iframeHTML = state.originalIframeHTML;

  // Extract sketch path from data attributes in the original HTML
  const sketchJsPathMatch = iframeHTML.match(
    /data-sketch-js-path=['"]([^'"]*)['"]/
  );

  if (sketchJsPathMatch) {
    const sketchJsPath = sketchJsPathMatch[1];
    console.log('🔍 Extracting data from original HTML...');
    console.log('📄 Sketch JS path match:', sketchJsPath);
  }

  // Replace the sketch-source script with the user's code
  const sketchSourceRegex = /<script id="sketch-source">[\s\S]*?<\/script>/;
  const userCodeScript = `<script id="sketch-source">
        // Store user's sketch source code directly in the global variable
        window.SKETCH_SOURCE_CODE = \`${userCode
          .replace(/`/g, '\\`')
          .replace(/\$/g, '\\$')}\`;
    </script>`;

  if (iframeHTML.match(sketchSourceRegex)) {
    iframeHTML = iframeHTML.replace(sketchSourceRegex, userCodeScript);
    console.log('✅ Replaced sketch-source script with user code');
  } else {
    console.error('❌ Could not find sketch-source script section to replace');
    console.log('🔍 Looking for pattern:', sketchSourceRegex);

    // Let's try to find the script in the body section specifically
    const bodyMatch = iframeHTML.match(/<body[\s\S]*<\/body>/);
    if (bodyMatch) {
      console.log(
        '🔍 HTML snippet in body section:',
        bodyMatch[0].substring(0, 500) + '...'
      );
    } else {
      console.log('🔍 Could not find body section');
    }
    return;
  }

  sketchDocument.open();
  sketchDocument.write(iframeHTML);
  sketchDocument.close();
  console.log('✅ Iframe document written and closed');

  // Show the sketch viewport and set overlay view mode
  console.log('📺 Sketch is running, showing sketch viewport');
  elements.sketchIframeContainer.classList.remove('hidden');
  setViewMode('overlay');

  if (window.parent && window.parent !== window) {
    console.log('📨 Sending sketchRunning message to parent');
    window.parent.postMessage({ type: 'sketchRunning' }, '*');
  }

  console.log('🎉 createAndRunSketch() completed successfully');
}

export function stopSketch() {
  // Remove the iframe from state if it exists
  if (state.sketchIframe) {
    state.sketchIframe.remove();
    state.sketchIframe = null;
  }

  // Also clean up any other iframes that might be in the container
  const existingIframes =
    elements.sketchIframeContainer.querySelectorAll('iframe');
  if (existingIframes.length > 0) {
    console.log(
      `🧹 Cleaning up ${existingIframes.length} remaining iframe(s) from container`
    );
    existingIframes.forEach((iframe) => iframe.remove());
  }

  // Hide the sketch viewport and switch to code view 
  console.log('📺 No sketch running, hiding sketch viewport');
  elements.sketchIframeContainer.classList.add('hidden');
  setViewMode('code');

  if (window.parent && window.parent !== window) {
    window.parent.postMessage({ type: 'sketchStopped' }, '*');
  }
}

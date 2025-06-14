/* Sketch Editor Main Module
 */

// Import all modules in dependency order
import './dom-elements.js';
import './console-manager.js';
import { updateVisibility } from './view-manager.js';
import { setupStatusTracking } from './status-tracker.js';
import { setupLineNumbers } from './line-numbers.js';
import {
  setupKeyboardShortcuts,
  setupMessageHandling,
  setupParentCommunication,
} from './keyboard-shortcuts.js';
import { initializeIframe, initializeUI } from './initialization.js';
import { createAndRunSketch } from './iframe-manager.js';

// Initialize the sketch editor application
async function initializeSketchEditor() {
  console.log('🚀 Initializing Sketch Editor...');

  // Setup UI and basic functionality
  initializeUI();

  // Setup event handlers
  setupStatusTracking();
  setupLineNumbers();
  setupKeyboardShortcuts();
  setupMessageHandling();
  setupParentCommunication();

  // Update visibility based on initial state
  updateVisibility();

  // Initialize iframe functionality
  await initializeIframe();

  console.log('✅ Sketch Editor initialized successfully');

  // Re-run the sketch after initialization
  createAndRunSketch();
}

// Start the application
initializeSketchEditor().catch((error) => {
  console.error('❌ Failed to initialize Sketch Editor:', error);
});

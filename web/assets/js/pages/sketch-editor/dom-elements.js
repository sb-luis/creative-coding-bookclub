// DOM element references and constants
export const elements = {
  codeEditorContainer: document.getElementById('code-viewport'),
  sketchIframeContainer: document.getElementById('sketch-viewport'),
  consoleOverlayContainer: document.getElementById('console-overlay'),
  consoleOutput: document.getElementById('console-output'),
  codeEditor: document.getElementById('code-editor'),
  cursorPositionEl: document.getElementById('cursor-position'),
  fileSizeEl: document.getElementById('file-size'),
  lineNumbersContainer: document.getElementById('line-numbers'),
  statusBar: document.getElementById('status-bar'),
};

// View modes: 'code', 'sketch', 'overlay', 'debug'
export const VIEW_MODES = {
  CODE: 'code', // Only code visible
  SKETCH: 'sketch', // Only sketch visible
  OVERLAY: 'overlay', // Both code and sketch visible (code-overlay mode)
  DEBUG: 'debug', // Sketch and console visible (console-overlay mode)
};

// Global state variables
export const state = {
  sketchIframe: null,
  originalIframeHTML: '',
  externalLibs: [],
  savedScrollTop: 0,
  savedSelectionStart: 0,
  savedSelectionEnd: 0,
  currentHighlightedLine: 1,
  isDirty: false, 
  currentViewMode: (() => {
    let serverViewMode = window.INITIAL_VIEW_MODE;

    // Clean up the view mode value (remove quotes if present)
    if (typeof serverViewMode === 'string') {
      serverViewMode = serverViewMode.replace(/^"|"$/g, '');
    }

    // Validate the server-provided view mode
    const validMode = Object.values(VIEW_MODES).includes(serverViewMode)
      ? serverViewMode
      : VIEW_MODES.CODE;

    return validMode;
  })(),
};

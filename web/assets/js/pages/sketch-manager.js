// Update member sketches
// DOM Elements
const sketchIframe = document.getElementById('sketch-viewport');
const saveButton = document.getElementById('save-button');
const newButton = document.getElementById('new-button');
const editMetadataButton = document.getElementById('edit-metadata-button');
const sketchSelector = document.getElementById('sketch-selector');
const deleteButton = document.getElementById('delete-button');
const sketchStatus = document.getElementById('sketch-status');
const saveDialog = document.getElementById('save-dialog');
const saveOverlay = document.getElementById('save-overlay');
const saveCancel = document.getElementById('save-cancel');
const saveConfirm = document.getElementById('save-confirm');
const sketchTitleInput = document.getElementById('sketch-title');
const sketchDescriptionInput = document.getElementById('sketch-description');
const sketchTagsInput = document.getElementById('sketch-tags');

// Top bar control elements
const prevSketchBtn = document.getElementById('prev-sketch');
const nextSketchBtn = document.getElementById('next-sketch');
const cycleViewBtn = document.getElementById('cycle-view');
const playStopBtn = document.getElementById('play-stop-sketch');
const sketchLinkEl = document.getElementById('sketch-link');

// Metadata dialog elements
const metadataDialog = document.getElementById('metadata-dialog');
const metadataOverlay = document.getElementById('metadata-overlay');
const metadataCancel = document.getElementById('metadata-cancel');
const metadataSave = document.getElementById('metadata-save');
const metadataTitleInput = document.getElementById('metadata-title');
const metadataDescriptionInput = document.getElementById(
  'metadata-description'
);
const metadataKeywordsInput = document.getElementById('metadata-keywords');
const metadataTagsInput = document.getElementById('metadata-tags');
const externalLibsContainer = document.getElementById(
  'external-libs-container'
);
const addExternalLibButton = document.getElementById('add-external-lib');

// Character counters
const titleCountSpan = document.getElementById('title-count');
const descriptionCountSpan = document.getElementById('description-count');
const keywordsCountSpan = document.getElementById('keywords-count');

// Global state
let currentSketch = null;
let sketches = [];
let hasUnsavedChanges = false;
let currentSketchIndex = -1;
let currentViewMode = 'overlay'; 
let isSketchRunning = false; 

// Initialize the IDE
document.addEventListener('DOMContentLoaded', async function () {
  console.log('🚀 DOM Content Loaded - Initializing sketch manager...');
  console.log('🔍 DOM Elements check:', {
    sketchIframe,
    saveButton,
    newButton,
    deleteButton,
    sketchSelector,
    sketchStatus,
    saveDialog,
    saveOverlay,
    saveCancel,
    saveConfirm,
    sketchTitleInput,
    sketchDescriptionInput,
    sketchTagsInput,
  });

  await loadSketches();
  setupEventListeners();
  loadEmptySketch();
  updateSketchStatus();
  updateNavigationButtons();
  updatePlayStopButton();

  sketchStatus.classList.remove('hidden');

  // Consolidated message handler for all iframe communications
  window.addEventListener('message', function (event) {
    // Only handle messages from the sketch iframe
    if (event.source === sketchIframe.contentWindow && event.data) {
      console.log('📨 Message received from iframe:', event.data);

      if (event.data.type === 'sketchDirty') {
        hasUnsavedChanges = event.data.status;
        updateSketchStatus();
      } else if (event.data.type === 'sketchRunning') {
        isSketchRunning = true;
        updatePlayStopButton();
      } else if (event.data.type === 'sketchStopped') {
        isSketchRunning = false;
        updatePlayStopButton();
      } else if (event.data.type === 'sketch-changed') {
        hasUnsavedChanges = true;
        updateSketchStatus();
      } else if (event.data.type === 'viewModeChanged') {
        let newViewMode = event.data.viewMode;
        console.log('📥 Raw view mode from iframe:', newViewMode);

        if (typeof newViewMode === 'string') {
          newViewMode = newViewMode.replace(/^"|"$/g, '');
        }

        currentViewMode = newViewMode;
        console.log('📥 View mode updated to:', currentViewMode);
      } else if (
        event.data.type === 'keyboardShortcut' &&
        event.data.shortcut === 'ctrl+s'
      ) {
        console.log('⌨️ Ctrl+S forwarded from iframe - saving sketch');

        // Check if there's anything to save
        if (currentSketch || hasUnsavedChanges) {
          saveSketch();
        } else {
          console.log('⚠️ No sketch to save or no unsaved changes');
        }
      }
    }
  });
});

// Load sketches from backend
async function loadSketches() {
  try {
    // First get the current member
    const memberResponse = await fetch('/api/members/me', {
      credentials: 'include',
    });

    if (!memberResponse.ok) {
      if (memberResponse.status === 401) {
        console.log('User not authenticated, redirecting to sign-in');
        window.location.href = '/sign-in';
        return;
      }
      throw new Error(`HTTP error! status: ${memberResponse.status}`);
    }

    const memberData = await memberResponse.json();
    const memberName = memberData.name;

    // Then load sketches for this member
    const response = await fetch(`/api/sketches/${memberName}`, {
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const sketchData = await response.json();
    sketches = Array.isArray(sketchData) ? sketchData : []; 
    updateSketchSelector();
    console.log(`Loaded ${sketches.length} sketches for ${memberName}`);
  } catch (error) {
    console.error('Error loading sketches:', error);
    sketches = []; // Ensure sketches is empty array on error
    updateSketchSelector();
  }
}

// Update sketch selector dropdown
function updateSketchSelector() {
  sketchSelector.innerHTML = '<option value="">Select a sketch...</option>';

  // Ensure sketches is an array and not null/undefined
  if (Array.isArray(sketches) && sketches.length > 0) {
    sketches.forEach((sketch) => {
      const option = document.createElement('option');
      option.value = sketch.slug; // Using slug instead of id as the identifier
      option.textContent = sketch.title;
      sketchSelector.appendChild(option);
    });
  }
  updateNavigationButtons();
}

// Load empty sketch
function loadEmptySketch() {
  const emptySketchUrl = `/sketches/${window.MEMBER_NAME}/new/edit`;
  sketchIframe.src = emptySketchUrl;
  
  // Notify iframe that a new sketch has been loaded 
  sketchIframe.addEventListener('load', function() {
    setTimeout(() => {
      if (sketchIframe && sketchIframe.contentWindow) {
        sketchIframe.contentWindow.postMessage(
          { type: 'sketchLoaded' },
          '*'
        );
      }
    }, 100);
  }, { once: true });
  
  updateSketchStatus();
}

// Load a specific sketch
async function loadSketch(sketchSlug) {
  try {
    // First get the current member
    const memberResponse = await fetch('/api/members/me', {
      credentials: 'include',
    });

    if (!memberResponse.ok) {
      throw new Error(`HTTP error! status: ${memberResponse.status}`);
    }

    const memberData = await memberResponse.json();
    const memberName = memberData.name;

    // Find the sketch from the already loaded sketches array
    if (!Array.isArray(sketches)) {
      throw new Error('Sketches not loaded properly');
    }
    
    const foundSketch = sketches.find((sketch) => sketch.slug === sketchSlug);

    if (!foundSketch) {
      throw new Error(`Sketch not found: ${sketchSlug}`);
    }

    currentSketch = foundSketch;
    currentSketchIndex = sketches.findIndex((sketch) => sketch.slug === sketchSlug);

    // Load sketch in iframe using the member name and slug
    const sketchUrl = `/sketches/${memberName}/${sketchSlug}/edit`;
    const sketchCleanUrl = `/sketches/${memberName}/${sketchSlug}`;
    
    sketchIframe.src = sketchUrl;
    sketchSelector.value = sketchSlug;
    
    if (sketchLinkEl) {
      sketchLinkEl.href = sketchCleanUrl;
    }

    hasUnsavedChanges = false; // Reset unsaved changes when loading a sketch
    
    // Notify iframe that a sketch has been loaded 
    sketchIframe.addEventListener('load', function() {
      setTimeout(() => {
        if (sketchIframe && sketchIframe.contentWindow) {
          sketchIframe.contentWindow.postMessage(
            { type: 'sketchLoaded' },
            '*'
          );
        }
      }, 100);
    }, { once: true });
    
    updateSketchStatus();
    updateNavigationButtons();
    console.log(`Loaded sketch: ${currentSketch.title}`);
  } catch (error) {
    console.error('Error loading sketch:', error);
  }
}

// Get code from iframe
function getCodeFromIframe() {
  try {
    const iframeContentWindow = sketchIframe.contentWindow;
    if (iframeContentWindow && iframeContentWindow.document) {
      const editor = iframeContentWindow.document.getElementById('code-editor');
      if (editor) {
        return editor.value;
      }
    }
  } catch (error) {
    console.error('Error getting code from iframe:', error);
  }
  return '';
}

// Save current sketch
async function saveSketch() {
  console.log('🔄 Starting saveSketch function...');

  const sourceCode = getCodeFromIframe();

  console.log('📝 Save data:', {
    sourceCodeLength: sourceCode.length,
    currentSketch: currentSketch
      ? {
          id: currentSketch.id,
          slug: currentSketch.slug,
          title: currentSketch.title,
        }
      : null,
  });

  if (!sourceCode) {
    console.error('❌ Validation failed: No source code found');
    alert(
      'No code found to save. Please check that the sketch is loaded properly.'
    );
    return;
  }

  // For new sketches, we only need source code. Metadata is handled in the backend 
  let title, description, tags;

  if (currentSketch) {
    // Existing sketch - use current metadata
    title = currentSketch.title;
    description = currentSketch.description || '';
    tags = currentSketch.tags || [];
  } else {
    // New sketch - 
    title = 'Auto-generated';
    description = 'Auto-generated sketch';
    tags = [];
  }

  if (!sourceCode) {
    console.error('❌ Validation failed: No source code found');
    console.log('🔍 Attempting to debug iframe access...');
    console.log('iframe element:', sketchIframe);
    console.log(
      'iframe src:',
      sketchIframe ? sketchIframe.src : 'iframe not found'
    );

    try {
      const iframeContentWindow = sketchIframe.contentWindow;
      console.log('iframe contentWindow:', iframeContentWindow);
      if (iframeContentWindow && iframeContentWindow.document) {
        console.log('iframe document found');
        const editor =
          iframeContentWindow.document.getElementById('code-editor');
        console.log('code-editor element:', editor);
        if (editor) {
          console.log('editor value:', editor.value);
        } else {
          console.log('❌ code-editor element not found in iframe');
        }
      } else {
        console.log(
          '❌ Cannot access iframe document (likely CORS or not loaded)'
        );
      }
    } catch (e) {
      console.error('❌ Error accessing iframe:', e);
    }

    alert(
      'No code found to save. Please check that the sketch is loaded properly.'
    );
    return;
  }


  try {
    let response;
    let responseData;

    console.log('🔐 Getting current user authentication...');
    // Get the current user's name
    const userResponse = await fetch('/api/members/me', {
      credentials: 'include',
    });

    console.log('👤 User response status:', userResponse.status);

    if (!userResponse.ok) {
      console.error(
        '❌ Failed to get current user:',
        userResponse.status,
        userResponse.statusText
      );
      throw new Error(`HTTP error! status: ${userResponse.status}`);
    }

    const userData = await userResponse.json();
    console.log('✅ Current user data:', userData);
    const memberName = userData.name;

    if (currentSketch && currentSketch.slug) {
      console.log(
        '🔄 Updating existing sketch source code:',
        currentSketch.slug
      );
      const updateUrl = `/api/sketches/${memberName}/${currentSketch.slug}`;
      console.log('📡 Update URL:', updateUrl);

      // Update existing sketch source code only via PUT
      const sourceCodeData = {
        source_code: sourceCode,
      };

      response = await fetch(updateUrl, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(sourceCodeData),
      });

      console.log('📡 Update response status:', response.status);
      console.log(
        '📡 Update response headers:',
        Object.fromEntries(response.headers.entries())
      );

      if (!response.ok) {
        const errorText = await response.text();
        console.error(
          '❌ Update failed:',
          response.status,
          response.statusText
        );
        console.error('❌ Update error body:', errorText);
        throw new Error(
          `HTTP error! status: ${response.status} - ${errorText}`
        );
      }

      responseData = await response.json();
      console.log('✅ Update successful! Response data:', responseData);

      // Update currentSketch with new data from response
      const sketchIndex = Array.isArray(sketches) ? sketches.findIndex(
        (s) => s.slug === currentSketch.slug
      ) : -1;
      console.log('🔍 Found sketch at index:', sketchIndex);
      if (sketchIndex !== -1 && Array.isArray(sketches)) {
        sketches[sketchIndex] = { ...sketches[sketchIndex], ...responseData };
        console.log('✅ Updated sketch in local array');
      }
      currentSketch = { ...currentSketch, ...responseData };
      console.log('✅ Updated currentSketch:', currentSketch);
    } else {
      console.log('🆕 Creating new sketch...');
      const createUrl = `/api/sketches/${memberName}/new`; // Use 'new' as placeholder - backend will ignore this
      console.log('📡 Create URL:', createUrl);

      // Create new sketch - only send source code
      const createData = {
        source_code: sourceCode,
      };

      response = await fetch(createUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(createData),
      });

      console.log('📡 Create response status:', response.status);
      console.log(
        '📡 Create response headers:',
        Object.fromEntries(response.headers.entries())
      );

      if (!response.ok) {
        const errorText = await response.text();
        console.error(
          '❌ Create failed:',
          response.status,
          response.statusText
        );
        console.error('❌ Create error body:', errorText);
        throw new Error(
          `HTTP error! status: ${response.status} - ${errorText}`
        );
      }

      responseData = await response.json();
      console.log('✅ Create successful! Response data:', responseData);

      // Add new sketch to list and set as current
      if (!Array.isArray(sketches)) {
        sketches = [];
      }
      sketches.push(responseData);
      console.log(
        '✅ Added new sketch to local array. Total sketches:',
        sketches.length
      );
      currentSketch = responseData;
      console.log('✅ Set currentSketch:', currentSketch);
    }

    hasUnsavedChanges = false;
    console.log('🔄 Updating UI after successful save...');
    updateSketchSelector();
    sketchSelector.value = currentSketch.slug;
    updateSketchStatus();
    console.log('✅ Save operation completed successfully!');

    // Notify iframe that sketch has been saved
    if (sketchIframe && sketchIframe.contentWindow) {
      sketchIframe.contentWindow.postMessage(
        { type: 'sketchSaved' },
        '*'
      );
    }

    // Show success message to user - use the title from the response
    const sketchTitle = responseData && responseData.title ? responseData.title : (currentSketch && currentSketch.title ? currentSketch.title : 'Sketch');
    alert(`Sketch "${sketchTitle}" saved successfully!`);
  } catch (error) {
    console.error('💥 Error in saveSketch:', error);
    console.error('💥 Error stack:', error.stack);
    let errorMessage = 'Failed to save sketch. Please try again.';
    alert(errorMessage);
  }
}

// Delete current sketch
async function deleteSketch() {
  console.log('🗑️ Starting deleteSketch function...');

  if (!currentSketch) {
    console.error('❌ No current sketch to delete');
    alert('No sketch selected to delete');
    return;
  }

  console.log('📋 Current sketch to delete:', {
    id: currentSketch.id,
    slug: currentSketch.slug,
    title: currentSketch.title,
  });

  if (
    !confirm(
      `Are you sure you want to delete "${currentSketch.title}"? This action cannot be undone.`
    )
  ) {
    console.log('❌ User cancelled delete operation');
    return;
  }

  try {
    console.log('🔐 Getting current user for delete operation...');
    // Get the current user's name
    const userResponse = await fetch('/api/members/me', {
      credentials: 'include',
    });

    if (!userResponse.ok) {
      console.error(
        '❌ Failed to get current user for delete:',
        userResponse.status,
        userResponse.statusText
      );
      throw new Error(`HTTP error! status: ${userResponse.status}`);
    }

    const userData = await userResponse.json();
    const memberName = userData.name;

    const deleteUrl = `/api/sketches/${memberName}/${currentSketch.slug}`;

    const response = await fetch(deleteUrl, {
      method: 'DELETE',
      credentials: 'include',
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('❌ Delete failed:', response.status, response.statusText);
      console.error('❌ Delete error body:', errorText);

      let errorData;
      try {
        errorData = JSON.parse(errorText);
      } catch (e) {
        errorData = { error: errorText };
      }
      throw new Error(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    const responseData = await response.json();
    console.log(
      `✅ Deleted ${currentSketch.title} successfully! Response:`,
      responseData
    );

    // Reset state and reload
    const deletedSketchTitle = currentSketch.title;
    currentSketch = null;
    console.log('🔄 Reloading sketches after delete...');
    await loadSketches();
    loadEmptySketch();
    sketchSelector.value = '';
    updateSketchStatus();

    console.log('✅ Delete operation completed successfully!');
    alert(`Sketch "${deletedSketchTitle}" deleted successfully!`);
  } catch (error) {
    console.error('💥 Error in deleteSketch:', error);
    console.error('💥 Error stack:', error.stack);

    alert('Failed to delete sketch. Please try again.');
  }
}

// Create new sketch
function newSketch() {
  if (hasUnsavedChanges) {
    if (
      !confirm(
        'You have unsaved changes. Are you sure you want to create a new sketch? Your changes will be lost.'
      )
    ) {
      return;
    }
  }

  // Call the internal function that doesn't check for unsaved changes
  newSketchInternal();
}

// Internal function to create new sketch without checking for unsaved changes
function newSketchInternal() {
  currentSketch = null;
  currentSketchIndex = -1;
  sketchSelector.value = '';
  
  if (sketchLinkEl) {
    sketchLinkEl.href = '/empty-iframe';
  }
  
  loadEmptySketch(); // This will load /sketch/new into the iframe

  hasUnsavedChanges = false;
  isSketchRunning = false;
  updateSketchStatus(); // Update button visibility
  updateNavigationButtons();
  updatePlayStopButton();
}

// Show/hide save dialog
function showSaveDialog() {
  console.log('💬 showSaveDialog() called');
  console.log('🔍 Save dialog elements:', {
    saveDialog,
    saveOverlay,
    sketchTitleInput,
  });

  saveDialog.classList.remove('hidden');
  saveOverlay.classList.remove('hidden');
  sketchTitleInput.focus();

  console.log('✅ Save dialog should now be visible');
}

function hideSaveDialog() {
  saveDialog.classList.add('hidden');
  saveOverlay.classList.add('hidden');
}

// Show/hide metadata dialog
function showMetadataDialog() {
  if (!currentSketch) {
    alert('Please select a sketch to edit metadata.');
    return;
  }

  // Populate the metadata form with current sketch data
  metadataTitleInput.value = currentSketch.title || '';
  metadataDescriptionInput.value = currentSketch.description || '';
  metadataKeywordsInput.value = currentSketch.keywords || '';
  metadataTagsInput.value = currentSketch.tags
    ? currentSketch.tags.join(',')
    : '';

  // Update character counters
  updateCharacterCount(metadataTitleInput, titleCountSpan, 100);
  updateCharacterCount(metadataDescriptionInput, descriptionCountSpan, 500);
  updateCharacterCount(metadataKeywordsInput, keywordsCountSpan, 200);

  // Clear and populate external libraries
  externalLibsContainer.innerHTML = '';
  const libs = currentSketch.external_libs || [];

  if (libs.length === 0) {
    // For sketches with no external libraries, add p5.js as default
    addExternalLibraryInput();
    const inputs = externalLibsContainer.querySelectorAll('.external-lib-url');
    inputs[inputs.length - 1].value =
      'https://cdn.jsdelivr.net/npm/p5@1.11.7/lib/p5.min.js';
  } else {
    libs.forEach((lib) => {
      addExternalLibraryInput();
      const inputs =
        externalLibsContainer.querySelectorAll('.external-lib-url');
      inputs[inputs.length - 1].value = lib;
    });
  }

  metadataDialog.classList.remove('hidden');
  metadataOverlay.classList.remove('hidden');
  metadataTitleInput.focus();
}

function hideMetadataDialog() {
  metadataDialog.classList.add('hidden');
  metadataOverlay.classList.add('hidden');
}

// Update sketch metadata via PATCH
async function updateSketchMetadata() {
  if (!currentSketch) {
    alert('No sketch selected to update metadata.');
    return;
  }

  const title = metadataTitleInput.value.trim();
  const description = metadataDescriptionInput.value.trim();
  const keywords = metadataKeywordsInput.value.trim();
  const tagsText = metadataTagsInput.value.trim();

  // Collect external libraries from inputs
  const externalLibInputs =
    externalLibsContainer.querySelectorAll('.external-lib-url');
  const externalLibs = Array.from(externalLibInputs)
    .map((input) => input.value.trim())
    .filter((url) => url.length > 0);

  // Validate all fields
  const titleError = validateMetadataTitle(title);
  if (titleError) {
    alert(titleError);
    metadataTitleInput.focus();
    return;
  }

  const descriptionError = validateMetadataDescription(description);
  if (descriptionError) {
    alert(descriptionError);
    metadataDescriptionInput.focus();
    return;
  }

  const keywordsError = validateMetadataKeywords(keywords);
  if (keywordsError) {
    alert(keywordsError);
    metadataKeywordsInput.focus();
    return;
  }

  const tagsError = validateMetadataTags(tagsText);
  if (tagsError) {
    alert(tagsError);
    metadataTagsInput.focus();
    return;
  }

  const externalLibsError = validateExternalLibraries(externalLibs);
  if (externalLibsError) {
    alert(externalLibsError);
    return;
  }

  const tags = tagsText
    ? tagsText
        .split(',')
        .map((tag) => tag.trim())
        .filter((tag) => tag)
    : [];

  const metadataData = {
    title: title,
    description: description,
    keywords: keywords,
    tags: tags,
    external_libs: externalLibs,
  };

  console.log('📝 Updating metadata:', metadataData);

  try {
    // Get the current user's name
    const userResponse = await fetch('/api/members/me', {
      credentials: 'include',
    });

    if (!userResponse.ok) {
      throw new Error(`HTTP error! status: ${userResponse.status}`);
    }

    const userData = await userResponse.json();
    const memberName = userData.name;

    const updateUrl = `/api/sketches/${memberName}/${currentSketch.slug}`;
    console.log('📡 Metadata update URL:', updateUrl);

    const response = await fetch(updateUrl, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(metadataData),
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('❌ Metadata update failed:', response.status, errorText);
      throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
    }

    const responseData = await response.json();
    console.log('✅ Metadata update successful:', responseData);

    // Update local sketch data
    const sketchIndex = Array.isArray(sketches) ? sketches.findIndex(
      (s) => s.slug === currentSketch.slug
    ) : -1;
    if (sketchIndex !== -1 && Array.isArray(sketches)) {
      sketches[sketchIndex] = { ...sketches[sketchIndex], ...responseData };
    }
    currentSketch = { ...currentSketch, ...responseData };

    // Update UI
    updateSketchSelector();
    sketchSelector.value = currentSketch.slug;
    updateSketchStatus();
    hideMetadataDialog();
    alert(`Metadata for "${title}" updated successfully!`);
  } catch (error) {
    console.error('💥 Error updating metadata:', error);
    alert('Failed to update metadata. Please try again.');
  }
}

// Update sketch status display
function updateSketchStatus() {
  console.log('🔍 updateSketchStatus called:', {
    currentSketch,
    hasCurrentSketch: !!currentSketch,
    hasId: currentSketch?.id,
    hasSlug: currentSketch?.slug,
    hasUnsavedChanges,
  });

  // Only show edit/delete buttons for sketches that have been saved (have both id and slug)
  if (currentSketch && currentSketch.id && currentSketch.slug) {
    sketchStatus.textContent = hasUnsavedChanges ? 'unsaved*' : currentSketch.title || 'untitled';
    // Show edit metadata and delete buttons for existing sketches
    editMetadataButton.classList.remove('hidden');
    deleteButton.classList.remove('hidden');
    console.log('✅ Showing edit/delete buttons for saved sketch');
  } else {
    sketchStatus.textContent = hasUnsavedChanges ? 'unsaved*' : 'untitled';
    // Hide edit metadata and delete buttons for new sketches
    editMetadataButton.classList.add('hidden');
    deleteButton.classList.add('hidden');
    console.log('❌ Hiding edit/delete buttons for new sketch');
  }

  // Show save button when there's content to save (either editing or new sketch with changes)
  if (currentSketch || hasUnsavedChanges) {
    saveButton.classList.remove('hidden');
  } else {
    saveButton.classList.add('hidden');
  }
}

// Update navigation buttons
function updateNavigationButtons() {
  const sketchesLength = Array.isArray(sketches) ? sketches.length : 0;
  prevSketchBtn.disabled = sketchesLength <= 1;
  nextSketchBtn.disabled = sketchesLength <= 1;
}

// Update play/stop button
function updatePlayStopButton() {
  if (isSketchRunning) {
    playStopBtn.innerHTML = 'stop';
    playStopBtn.title = 'Stop sketch (Ctrl+.)';
  } else {
    playStopBtn.innerHTML = 'play';
    playStopBtn.title = 'Run sketch (Ctrl+Enter)';
  }
}

// Load sketch by index for navigation
function loadSketchByIndex(index) {
  const sketchesLength = Array.isArray(sketches) ? sketches.length : 0;
  if (index < 0 || index >= sketchesLength) {
    console.warn(`Invalid sketch index: ${index}`);
    return;
  }

  currentSketchIndex = index;
  const sketch = sketches[index];

  if (sketch) {
    isSketchRunning = false;
    updatePlayStopButton();
    sketchSelector.value = sketch.slug;
    loadSketch(sketch.slug);
    updateNavigationButtons();
  }
}

// Navigate to previous sketch
function navigateToPreviousSketch() {
  const sketchesLength = Array.isArray(sketches) ? sketches.length : 0;
  if (sketchesLength > 1) {
    if (hasUnsavedChanges) {
      if (!confirm('You have unsaved changes. Are you sure you want to switch sketches?')) {
        return;
      }
    }
    
    const newIndex = currentSketchIndex > 0 ? currentSketchIndex - 1 : sketchesLength - 1;
    loadSketchByIndex(newIndex);
  }
}

// Navigate to next sketch
function navigateToNextSketch() {
  const sketchesLength = Array.isArray(sketches) ? sketches.length : 0;
  if (sketchesLength > 1) {
    if (hasUnsavedChanges) {
      if (!confirm('You have unsaved changes. Are you sure you want to switch sketches?')) {
        return;
      }
    }
    
    const newIndex = currentSketchIndex < sketchesLength - 1 ? currentSketchIndex + 1 : 0;
    loadSketchByIndex(newIndex);
  }
}

// Cycle view mode
function cycleViewMode() {
  if (sketchIframe && sketchIframe.contentWindow) {
    sketchIframe.contentWindow.postMessage(
      {
        type: 'cycleViewMode',
      },
      '*'
    );
  }
}

// Toggle play/stop
function togglePlayStop() {
  if (sketchIframe && sketchIframe.contentWindow) {
    if (isSketchRunning) {
      // Stop the sketch 
      sketchIframe.contentWindow.postMessage(
        {
          type: 'keyboardShortcut',
          shortcut: 'ctrl+period',
        },
        '*'
      );
    } else {
      // Run the sketch 
      sketchIframe.contentWindow.postMessage(
        {
          type: 'keyboardShortcut',
          shortcut: 'ctrl+enter',
        },
        '*'
      );
    }
  }
}

// Setup event listeners
function setupEventListeners() {
  console.log('🔧 Setting up event listeners...');
  console.log('🔍 Button elements:', {
    saveButton,
    newButton,
    deleteButton,
    saveCancel,
    saveConfirm,
  });

  // Backend integration buttons
  saveButton.addEventListener('click', saveSketch);
  newButton.addEventListener('click', newSketch);
  editMetadataButton.addEventListener('click', showMetadataDialog);
  deleteButton.addEventListener('click', deleteSketch);

  // Top bar navigation controls
  prevSketchBtn.addEventListener('click', navigateToPreviousSketch);
  nextSketchBtn.addEventListener('click', navigateToNextSketch);
  cycleViewBtn.addEventListener('click', cycleViewMode);
  playStopBtn.addEventListener('click', togglePlayStop);

  // Sketch selector
  sketchSelector.addEventListener('change', function () {
    if (this.value) {
      // User selected a specific sketch
      if (hasUnsavedChanges) {
        if (
          !confirm(
            'You have unsaved changes. Are you sure you want to switch sketches?'
          )
        ) {
          this.value = currentSketch ? currentSketch.slug : '';
          return;
        }
      }
      loadSketch(this.value);
    } else {
      // User selected "Select a sketch..." (empty value) - treat as "New Sketch"
      if (hasUnsavedChanges) {
        if (
          !confirm(
            'You have unsaved changes. Are you sure you want to create a new sketch? Your changes will be lost.'
          )
        ) {
          this.value = currentSketch ? currentSketch.slug : '';
          return;
        }
      }
      // Use internal function to avoid double-prompting
      newSketchInternal();
    }
  });

  // Save dialog events
  saveCancel.addEventListener('click', hideSaveDialog);
  saveConfirm.addEventListener('click', saveSketch);
  saveOverlay.addEventListener('click', hideSaveDialog);

  // Metadata dialog events
  metadataCancel.addEventListener('click', hideMetadataDialog);
  metadataSave.addEventListener('click', updateSketchMetadata);
  metadataOverlay.addEventListener('click', hideMetadataDialog);

  // Character counter updates
  metadataTitleInput.addEventListener('input', function () {
    updateCharacterCount(this, titleCountSpan, 100);
  });

  metadataDescriptionInput.addEventListener('input', function () {
    updateCharacterCount(this, descriptionCountSpan, 500);
  });

  metadataKeywordsInput.addEventListener('input', function () {
    updateCharacterCount(this, keywordsCountSpan, 200);
  });

  // Add external library button
  addExternalLibButton.addEventListener('click', addExternalLibraryInput);

  // Initialize with first external lib input
  addExternalLibraryInput();

  // Setup keyboard shortcuts
  setupKeyboardShortcuts();
}

// Setup keyboard shortcuts for the sketch manager
function setupKeyboardShortcuts() {
  // Handle keyboard shortcuts when focus is outside the iframe
  document.addEventListener('keydown', function (event) {
    console.log('🎹 keydown in sketch manager document:', {
      key: event.key,
      ctrlKey: event.ctrlKey,
      activeElement: document.activeElement?.tagName,
    });

    // Ctrl + S: Save current sketch
    if (event.ctrlKey && event.key === 's') {
      event.preventDefault();
      console.log('⌨️ Ctrl+S pressed - saving sketch from manager');

      // Check if there's anything to save
      if (currentSketch || hasUnsavedChanges) {
        saveSketch();
      } else {
        console.log('⚠️ No sketch to save or no unsaved changes');
      }
    }
  });
}

// Character counter update function
function updateCharacterCount(input, counterSpan, maxLength) {
  const currentLength = input.value.length;
  counterSpan.textContent = currentLength;
}

// Add external library input field
function addExternalLibraryInput() {
  const inputs = externalLibsContainer.querySelectorAll('.external-lib-input');
  if (inputs.length >= 5) {
    alert('Maximum 5 external libraries allowed');
    return;
  }

  const inputDiv = document.createElement('div');
  inputDiv.className = 'external-lib-input mb-2 flex items-center';

  const input = document.createElement('input');
  input.type = 'url';
  input.className =
    'external-lib-url flex-1 p-2 border border-base-300 rounded bg-base-100';
  input.placeholder = 'https://cdn.jsdelivr.net/npm/package@version/file.js';

  const removeBtn = document.createElement('button');
  removeBtn.type = 'button';
  removeBtn.className =
    'ml-2 px-2 py-1 text-xs text-red-600 hover:text-red-800';
  removeBtn.textContent = '✕';
  removeBtn.onclick = function () {
    inputDiv.remove();
    updateAddButtonVisibility();
  };

  inputDiv.appendChild(input);
  if (inputs.length > 0) {
    // Don't add remove button to first input
    inputDiv.appendChild(removeBtn);
  }

  externalLibsContainer.appendChild(inputDiv);
  updateAddButtonVisibility();
}

// Update add button visibility
function updateAddButtonVisibility() {
  const inputs = externalLibsContainer.querySelectorAll('.external-lib-input');
  addExternalLibButton.style.display = inputs.length >= 5 ? 'none' : 'inline';
}

// Validation functions
function validateMetadataTitle(title) {
  if (!title || title.trim().length === 0) {
    return 'Title is required';
  }
  if (title.length > 100) {
    return 'Title must be 100 characters or less';
  }
  // Allow alphanumeric, spaces, and common punctuation
  if (!/^[a-zA-Z0-9\s\-_.,:;!?()]+$/.test(title)) {
    return 'Title contains invalid characters';
  }
  return null;
}

function validateMetadataDescription(description) {
  if (description && description.length > 500) {
    return 'Description must be 500 characters or less';
  }
  if (description && !/^[a-zA-Z0-9\s\-_.,:;!?()]*$/.test(description)) {
    return 'Description contains invalid characters';
  }
  return null;
}

function validateMetadataKeywords(keywords) {
  if (keywords && keywords.length > 200) {
    return 'Keywords must be 200 characters or less';
  }
  if (keywords && !/^[a-zA-Z0-9\s,\-_]*$/.test(keywords)) {
    return 'Keywords contains invalid characters';
  }
  return null;
}

function validateMetadataTags(tagsString) {
  if (!tagsString) return null;

  const tags = tagsString
    .split(',')
    .map((tag) => tag.trim())
    .filter((tag) => tag.length > 0);

  if (tags.length > 10) {
    return 'Maximum 10 tags allowed';
  }

  for (const tag of tags) {
    if (!/^[a-zA-Z0-9\-_]+$/.test(tag)) {
      return `Tag "${tag}" contains invalid characters (only alphanumeric, hyphens, and underscores allowed)`;
    }
  }

  return null;
}

function validateExternalLibraries(urls) {
  if (urls.length > 5) {
    return 'Maximum 5 external libraries allowed';
  }

  for (const url of urls) {
    if (!url || url.trim().length === 0) {
      return 'External library URLs cannot be empty';
    }

    try {
      new URL(url);
    } catch {
      return `Invalid URL format: ${url}`;
    }

    if (!url.startsWith('https://')) {
      return `External library URLs must use HTTPS: ${url}`;
    }
  }

  return null;
}

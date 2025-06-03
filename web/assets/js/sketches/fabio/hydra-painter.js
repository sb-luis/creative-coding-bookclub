/*
Simple p5.js sketch that uses a Hydra brush texture for drawing.
Canvas: 1000 x 1000 px, white background.
Drawing occurs with the Hydra brush texture when the mouse is pressed.
*/

let hydraCanvas, hydra, pg, canvasLayer; // Declare Hydra and p5 graphics canvas

function setup() {
  createCanvas(windowWidth, windowHeight); // Set up main canvas
  background(255); // Set background to white
  noSmooth(); // Disable smoothing for better brush effect
  noStroke(); // No outline for drawings
  noCursor(); // Hide the default mouse cursor

  // Create an offscreen Hydra canvas
  hydraCanvas = document.createElement("canvas");
  hydraCanvas.width = windowWidth; // Set width
  hydraCanvas.height = windowHeight; // Set height
  document.body.appendChild(hydraCanvas); // Add to document
  hydraCanvas.style.display = "none"; // Hide Hydra canvas from view

  // Load Hydra library dynamically
  let script = document.createElement("script");
  script.src = "https://unpkg.com/hydra-synth";
  script.onload = () => {
    hydra = new Hydra({ canvas: hydraCanvas, detectAudio: false }); // Initialize Hydra
    initializeHydra(); // Start Hydra texture
  };
  document.head.appendChild(script);

  pg = createGraphics(500, 500); // Create a p5 graphics buffer for the brush
  canvasLayer = createGraphics(windowWidth, windowHeight); // Create an additional layer to store drawn strokes
}

function initializeHydra() {

  // Define a Hydra pattern for the brush
  osc(10, 0.1, 0.3) // Generate an oscillating pattern
    .kaleid(5) // Apply a kaleidoscope effect
    .color(0.5, 0.3) // Adjust color channels
    .colorama(0.04) // Apply colorama effect
    .rotate(0.009, () => Math.sin(time) * -0.001) // Slight rotation animation
    .modulateRotate(src(o0), () => Math.sin(time) * 0.003) // Modulate rotation
    .modulate(src(o0), 0.9) // Modulate with existing output
    .mask(shape(100, 0.5, 0.01)) // Apply a circular mask (100px diameter)
    .out(); // Output Hydra pattern
}

function draw() {
  clear(); // Clear the main canvas
  background(255); // Maintain white background

  if (hydra) {
    pg.clear(); // Clear graphics buffer
    pg.drawingContext.drawImage(hydraCanvas, 0, 0, pg.width, pg.height); // Draw Hydra texture
  }

  // Display the brush as the cursor
  image(pg, mouseX - 50, mouseY - 50, 300, 300);

  if (mouseIsPressed) {
    // Draw the Hydra brush texture on the persistent canvas layer when the mouse is pressed
    canvasLayer.image(pg, mouseX - 50, mouseY - 50, 300, 300);
  }

  // Render the persistent layer with drawn strokes
  image(canvasLayer, 0, 0);
}
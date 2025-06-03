class Walker {
  constructor() {
    this.x = width / 2;
    this.y = height / 2;
  }
  show() {
    stroke(255);
    point(this.x, this.y);
  }
  step() {
    let xstep = randomGaussian(0, 1);
    let ystep = randomGaussian(0, 1);
    this.x += xstep;
    this.y += ystep;
    
    if (random(0, 1) > 0.9) {
      this.x = (this.x < mouseX) ? this.x + 1 : this.x - 1
      this.y = (this.x < mouseY) ? this.y + 1 : this.y - 1
    }
  }
}

let walker;
let x;
let y;
let t = 0;

function setup() {
  createCanvas(600, 600);
  
  pixelDensity(1);
  loadPixels();
  // Start xoff at 0.
  let xoff = 0.0;

  for (let x = 0; x < width; x++) {
    // For every xoff, start yoff at 0.
    let yoff = 0.0;

    for (let y = 0; y < height; y++) {
      // Use xoff and yoff for noise().
      let bright = map(noise(xoff, yoff), 0, 1, 0, 255);
      // Use x and y for the pixel position.
      let index = (x + y * width) * 4;
      // Set the red, green, blue, and alpha values.
      pixels[index] = bright;
      pixels[index + 1] = bright;
      pixels[index + 2] = 255;
      pixels[index + 3] = 255;
      // Increment yoff.
      yoff += 0.01;
    }
    // Increment xoff.
    xoff += 0.01;
    noiseDetail(10);
  }
  updatePixels();
  
  walker = new Walker();
}


function draw() {
  walker.step();
  walker.show();
    
  x = randomGaussian(mouseX, 10);
  y = randomGaussian(mouseY, 10);
  noStroke();
  fill(255, 10);
  circle(x, y, 3);
  
  let n_t = noise(t);
  let n_tp1 = noise(t + 0.5);
  let m_t = map(n_t, 0, 1, 200, 400);
  let m_tp1 = map(n_tp1, 0, 1, 200, 400);
  
  noFill()
  stroke(
    map(n_t, 0, 1, 0, 255), 
    map(n_tp1, 0, 1, 0, 255), 
    255);
  ellipse(m_t, m_tp1, 600 + Math.sin(t), 600 + Math.cos(t));
  t += 0.01;
}


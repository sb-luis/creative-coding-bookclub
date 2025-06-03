let walker;

function setup() {
  createCanvas(windowWidth, windowHeight); // creating canvas of the size of the window
  walker = new Walker(); // creating an instance/object of class Walker
  background(25);

  // Draw text (two lines)
  fill(255);
  textSize(30);
  textAlign(CENTER, CENTER);
  textFont("monospace");

  // First line
  text("The Nature of", width / 2, height / 2);
  
  // Second line
  text("Coding", width / 2, height / 2 + 40);

}

function draw() {
  walker.step();
  walker.show();

}

class Walker {
  constructor() {
    this.tx = 0;
    this.ty = 10000;
  }

  step() {
    // x- and y-position mapped from noise
    this.x = map(noise(this.tx), 0, 1, 0, width);
    this.y = map(noise(this.ty), 0, 1, 0, height);

    // Move forward through “time.”
    this.tx += 0.01;
    this.ty += 0.01;
  }

  show() {
    strokeWeight(2);
    fill(randomGaussian(100, 200), randomGaussian(100, 200), 149);
    stroke(25);
    circle(this.x, this.y, 100);
  }
}

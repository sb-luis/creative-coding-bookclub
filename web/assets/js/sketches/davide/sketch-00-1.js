class Creature {
  constructor() {
    this.speed = 4;
    this.fill = 'white';
    this.strokeWeight = 1;
    this.stroke = 'red';

    this.noiseOffset = random(1000);

    // Initialize at a random point on the canvas border
    let edge = floor(random(4));
    if (edge === 0) {
      // Top edge
      this.x = random(width);
      this.y = 0;
    } else if (edge === 1) {
      // Right edge
      this.x = width;
      this.y = random(height);
    } else if (edge === 2) {
      // Bottom edge
      this.x = random(width);
      this.y = height;
    } else {
      // Left edge
      this.x = 0;
      this.y = random(height);
    }

    // Random initial direction
    this.angle = random(TWO_PI);
  }

  move() {
    // Use Perlin noise to determine the direction
    this.angle += map(noise(this.noiseOffset), 0, 1, -0.1, 0.1);
    this.noiseOffset += 0.01;

    this.x += cos(this.angle) * this.speed;
    this.y += sin(this.angle) * this.speed;

    // Change direction when hitting the border
    if (this.x <= 0 || this.x >= width) {
      this.angle = PI - this.angle;
    }
    if (this.y <= 0 || this.y >= height) {
      this.angle = -this.angle;
    }

    // Bound the creature inside the canvas
    this.x = constrain(this.x, 0, width);
    this.y = constrain(this.y, 0, height);
  }

  display() {
    fill(this.fill);
    strokeWeight(this.strokeWeight);
    stroke(this.stroke);
    ellipse(this.x, this.y, 20, 20);
  }

  update() {
    this.move();
    this.display();
  }
}

const creaturesCount = 3;
let creaturesArray = [];

const message = Math.random() < 0.5 ? 'Hello World' : 'Let Code Die';

function setup() {
  createCanvas(500, 500);

  background(250);

  noFill();
  stroke(230);
  rect(0, 0, width, height);

  fill('red');
  textSize(25);
  textAlign(CENTER);
  textStyle(BOLDITALIC);
  textFont('Times New Roman');
  text(message, width / 2, height / 2);

  for (let i = 0; i < creaturesCount; i++) {
    creaturesArray.push((creature = new Creature()));
  }
}

function draw() {
  creaturesArray.forEach((creature) => {
    creature.update();
  });
}

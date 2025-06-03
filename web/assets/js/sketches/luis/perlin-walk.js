class Walker {
  constructor(x, y) {
    this.x = x;
    this.y = x;
    this.xOffset = 0;
    this.yOffset = 1000;
  }

  step() {
    this.x = map(noise(this.xOffset), 0, 1, 0, width);
    this.y = map(noise(this.yOffset), 0, 1, 0, height);

    this.xOffset += 0.01;
    this.yOffset += 0.01;
  }

  render() {
    stroke(0);
    circle(this.x, this.y, 5);
  }
}

let walker = null;

function setup() {
  createCanvas(windowWidth, windowHeight);
  walker = new Walker(windowWidth / 2, windowHeight / 2);
  background(255);
}

function draw() {
  walker.step();
  walker.render();
}

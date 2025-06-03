class CCBPerlin {
  constructor() {
    this.x = windowWidth / Math.ceil(random(10));
    this.y = windowHeight / Math.ceil(random(10));
    this.xoff = Math.ceil(random(10000));
    this.yoff = Math.ceil(random(10000));
  }

  mutatePerlin() {
    this.x = map(noise(this.xoff), 0, 2, 0, windowWidth);
    this.y = map(noise(this.yoff), 0, 1, 0, windowHeight);

    this.xoff += 0.01;
    this.yoff += 0.01;
  }

  visualise() {
    const r = Math.ceil(random(4));
    if (r === 1) {
      fill("hotpink");
      stroke("blue");
    } else if (r === 2) {
      fill("coral");
      stroke("darkgreen");
    } else if (r === 3) {
      fill("pink");
      stroke("indigo");
    } else {
      fill("lightgreen");
      stroke("crimson");
    }
    strokeWeight(20);
    textFont("Courier New", 44);
    text("Creative Coding Bookclub", this.x, this.y);
    noStroke();
  }
}

class CCBGauss extends CCBPerlin {
  mutateGauss() {
    this.x = randomGaussian(windowWidth / 2, windowWidth / 2);
    this.y = randomGaussian(windowHeight, windowHeight);
  }
}

let count = 0;

function setup() {
  createCanvas(windowWidth, windowHeight);
  background(255, 229, 180);
  ccb = new CCBPerlin();
  ccb2 = new CCBGauss();
}

function draw() {
  ccb.mutatePerlin();
  ccb.visualise();
  ccb2.mutateGauss();
  ccb2.visualise();
  count++;
  if (count === 1000) {
    erase();
    rect(0, 0, windowWidth, windowHeight);
    noErase();
    background(255, 229, 180);
    count = 0;
  }
}

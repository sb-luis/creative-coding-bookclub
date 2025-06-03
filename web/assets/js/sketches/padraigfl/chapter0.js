const defaultHsla = {
  h: 0,
  s: 0,
  l: 0,
  a: 100,
};

class P5Node {
  constructor(params) {
    this.fill = { ...defaultHsla, ...params.fill };
    this.stroke = { ...defaultHsla, ...params.stroke };
    this.strokeWeight = params.strokeWeight;
    this.height = params.height;
    this.width = params.width;
    this.x = params.x;
    this.y = params.y;
  }
  getHsla = (param) => `hsla(${param.h},${param.s}%,${param.l}%,${param.a})`;
  setFill = (color) => {
    this.fill = { ...this.fill, ...color };
  };
  getFill = () => this.getHsla(this.fill);
  getStroke = () => this.getHsla(this.stroke);
}


// Instantiate P5 sketch
new p5((p5Instance) => {
  sketch0(p5Instance);
});

function sketch0(p5, options) {
  const canvasSize = {
    width: options?.canvasWidth || 400,
    height: options?.canvasHeight || 400,
  };

  // a crude class for trailing behind the original walker
  class Trailer extends P5Node {
    constructor(params) {
      super(params);
      this.prevStates = [];
      this.delay = params.delay;
      this.parent = params.parent;
    }
    move = () => {
      const { x, y, w, h } = this.prevStates.shift() || {};
      // console.log(x, y);
      this.x = x || this.x;
      this.y = y || this.y;
      // this.width = w || this.width;
      // this.height = h || this.height;
    };
    show = () => {
      p5.fill(this.getFill());
      p5.strokeWeight(this.strokeWeight);
      p5.stroke(this.getStroke());
      p5.ellipse(this.x, this.y, this.width, this.height);
    };
    update = () => {
      if (this.prevStates.length > this.delay) {
        this.move();
        this.show();
      }
      // console.log(this.prevStates.length);
      if (this.parent) {
        this.prevStates.push({
          x: this.parent.x,
          y: this.parent.y,
          w: this.parent.width - 2,
          h: this.parent.height - 2,
        });
      }
    };
  }

  class Walker extends P5Node {
    trailers = [
      new Trailer({
        stroke: { l: 20, a: 0.5 },
        fill: { l: 0 },
        strokeWeight: 2,
        width: 38,
        height: 38,
        x: 150,
        y: 150,
        delay: 10,
        parent: this,
      }),
      new Trailer({
        stroke: { l: 0 },
        fill: { l: 0, a: 0 },
        strokeWeight: 2,
        width: 38,
        height: 38,
        x: 150,
        y: 150,
        delay: 100,
        parent: this,
      }),
    ];
    move = () => {
      let choice = p5.floor(p5.random(4));

      let steps = p5.map(
        choice % 2 === 0
          ? p5.noise(1 - this.x / this.width)
          : p5.noise(1 - this.y / this.height),
        0,
        1,
        0,
        30
      );
      if (choice === 0) {
        this.x += steps;
      } else if (choice === 1) {
        this.y += steps;
      } else if (choice === 2) {
        this.x -= steps;
      } else {
        this.y -= steps;
      }
      // choice = p5.floor(p5.random(2));
      // if (choice === 0 && this.height < 120) {
      //   this.height++;
      //   this.width++;
      // } else if (choice === 1 && this.height > 10) {
      //   this.height--;
      //   this.width--;
      // }
      if (this.x < 0 || this.x > canvasSize.width) {
        const rangeStart = p5.floor(canvasSize.width / 3);
        this.x = p5.floor(p5.random(rangeStart, rangeStart * 2));
      }
      if (this.y < 0 || this.y > canvasSize.height) {
        const rangeStart = p5.floor(canvasSize.height / 3);
        this.y = p5.floor(p5.random(rangeStart, rangeStart * 2));
      }
    };
    colorChange = () => {
      let choice = p5.floor(p5.random(10));
      if (choice > 7) {
        this.fill.h = ++this.fill.h % 360;
      }
      if (choice < 2) {
        this.stroke.h = ((this.stroke.h + 360 - 2) ^ (choice + 1)) % 360;
      }
    };
    show = () => {
      p5.fill(this.getFill());
      p5.strokeWeight(this.strokeWeight);
      p5.stroke(this.getStroke());
      p5.ellipse(this.x, this.y, this.width, this.height);
    };
    update = () => {
      this.move();
      this.colorChange();
      this.show();
      this.trailers.forEach((trailer) => trailer.update());
    };
  }

  const w = new Walker({
    fill: { h: p5.floor(p5.random(0, 360)), s: 60, l: 50, a: 1 },
    stroke: { h: p5.floor(p5.random(0, 360)), s: 60, l: 50, a: 1 },
    strokeWeight: 5,
    width: 40,
    height: 40,
    x: 300,
    y: 300,
  });

  p5.setup = () => {
    p5.createCanvas(canvasSize.width, canvasSize.height);
    p5.background('black');
  };
  p5.draw = () => {
    w.update();
  };
}

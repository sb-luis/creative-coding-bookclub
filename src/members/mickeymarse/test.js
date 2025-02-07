export default function (p5, options) {
  const canvas = options?.canvas ?? undefined;

  function drawSketch() {
    p5.background("deeppink");
    p5.fill("lime");
    p5.strokeWeight(4);
    p5.textSize(56);
    p5.text("test test soundcheck", p5.windowWidth / 3, p5.windowHeight / 2);
  }

  p5.setup = () => {
    p5.createCanvas(p5.windowWidth, p5.windowHeight, canvas);
    drawSketch();
  };
}

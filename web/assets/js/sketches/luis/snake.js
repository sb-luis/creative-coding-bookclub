// Instantiate P5 sketch
new p5((p5Instance) => {
  sketch0(p5Instance);
});

function sketch0(p5, options) {
  const getHeight = options?.getHeight ?? (() => p5.windowHeight)
  const canvas = options?.canvas ?? undefined

  function drawSketch() {
    p5.background('lightgrey')
    p5.noFill()

    const stepSize = 5
    const circleSize = p5.constrain(p5.height * 0.25, 30, 55)
    const steps = (p5.windowWidth + circleSize) / stepSize
    const yOffset = p5.height * 0.4
    const waveAmplitude1 = p5.random(60, 100)
    const waveAmplitude2 = p5.random(20, 30)
    const waveAmplitude3 = p5.random(5, 10)
    const wavePeriodOffset1 = p5.random(0, 10000)
    const wavePeriodOffset2 = p5.random(0, 10000)
    const wavePeriodOffset3 = p5.random(0, 10000)

    const makeWave = (amplitude, period) => {
      return amplitude * p5.sin(period)
    }

    for (let i = 0; i < steps; i++) {
      const x = stepSize * i
      const wave1 = makeWave(waveAmplitude1, (p5.TWO_PI * i + wavePeriodOffset1) / 300)
      const wave2 = makeWave(waveAmplitude2, (p5.TWO_PI * i + wavePeriodOffset2) / 58)
      const wave3 = makeWave(waveAmplitude3, (p5.TWO_PI * i + wavePeriodOffset3) / 23)

      p5.circle(x - circleSize * 0.5, yOffset + wave1 + wave2 + wave3, circleSize)
    }
  }

  // https://p5js.org/reference/p5/setup/
  p5.setup = () => {
    p5.createCanvas(p5.windowWidth, getHeight(), canvas)
    drawSketch()
  }

  // https://p5js.org/reference/p5/windowResized/
  p5.windowResized = () => {
    if (p5.windowWidth !== p5.width) {
      p5.resizeCanvas(p5.windowWidth, getHeight())
      drawSketch()
    }
  }
}

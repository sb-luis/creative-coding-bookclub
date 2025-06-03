export default function (p5, options) {
  const getHeight = options?.getHeight ?? (() => p5.windowHeight);
  const canvas = options?.canvas ?? undefined;


  // Helper function to determine stroke color based on theme
  function getStrokeColor() {
    const themeAttribute = document.documentElement.getAttribute('data-theme');
    let isDark;

    if (themeAttribute === 'dark') {
      isDark = true;
    } else if (themeAttribute === 'light') {
      isDark = false;
    } else { // System theme (attribute not set or not 'light'/'dark')
      isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    return isDark ? 'white' : 'black';
  }

  function drawSketch() {
    let currentStrokeColor = getStrokeColor();
    getStrokeColor(); // Get latest stroke color
    p5.clear(); // Clear the canvas instead of setting a background color
    p5.stroke(currentStrokeColor); // Use theme color for circles
    p5.noFill(); // Keep circles as outlines

    const stepSize = 5;
    const circleSize = p5.constrain(p5.height * 0.25, 30, 55);
    const steps = (p5.windowWidth + circleSize) / stepSize;
    const yOffset = p5.height * 0.4;
    const waveAmplitude1 = p5.random(60, 100);
    const waveAmplitude2 = p5.random(20, 30);
    const waveAmplitude3 = p5.random(5, 10);
    const wavePeriodOffset1 = p5.random(0, 10000);
    const wavePeriodOffset2 = p5.random(0, 10000);
    const wavePeriodOffset3 = p5.random(0, 10000);

    const makeWave = (amplitude, period) => {
      return amplitude * p5.sin(period);
    };

    for (let i = 0; i < steps; i++) {
      const x = stepSize * i;
      const wave1 = makeWave(
        waveAmplitude1,
        (p5.TWO_PI * i + wavePeriodOffset1) / 300
      );
      const wave2 = makeWave(
        waveAmplitude2,
        (p5.TWO_PI * i + wavePeriodOffset2) / 58
      );
      const wave3 = makeWave(
        waveAmplitude3,
        (p5.TWO_PI * i + wavePeriodOffset3) / 23
      );

      p5.circle(
        x - circleSize * 0.5,
        yOffset + wave1 + wave2 + wave3,
        circleSize
      );
    }
  }

  p5.setup = () => {
    p5.createCanvas(p5.windowWidth, getHeight(), canvas);
    drawSketch(); // Initial draw
  };

  p5.windowResized = () => {
    if (p5.windowWidth !== p5.width) {
      p5.resizeCanvas(p5.windowWidth, getHeight());
      drawSketch(); // Redraw on resize
    }
  };
}

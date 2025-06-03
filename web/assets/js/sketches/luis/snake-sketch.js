<script>
  // Import p5 from NPM
  import { default as P5Sketch } from 'p5';
  import sketch from './snake.js'

  // Instantiate P5 sketch
  new P5Sketch((p5) => {
    sketch(p5);
  });
</script>

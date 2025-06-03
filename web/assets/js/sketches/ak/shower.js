console.log('hi from ak sketch' )

console.log(p5)

function setup() {
  console.log('setting up sketch')
  createCanvas(400, 400);
  background("#e2f1f3");
  
  for (let y = 40; y < 800; y += 40) {
    
    stroke("#d4e9e9");
    line(0, y, 800, y);
    
  }
  for (let x = 40; x < 800; x += 40) {
    stroke("#d4e9e8");
    line(x, 0, x, 800);
    
  }
 
  textFont('Courier new');
  fill('palegreen');
  textSize(35);
  textLeading()
  let fontItalic
  text('shower \nthoughts', 200, 300);
}

function draw() {
  let a = random(0,400)
  let b = random(0,400)
  let c = random(10,35)
  fill(255,255,255,random(0,100))
  
  circle(a,b,c)
  fill(0,0,0,random(0,10))
  stroke(255)
}


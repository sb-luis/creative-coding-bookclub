// Returns a random float between min (inclusive) and max (exclusive)
export function getRandomFloat(min, max) {
  return Math.random() * (max - min) + min;
}

// Returns a random integer between min (inclusive) and max (inclusive)
export function getRandomInt(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

// Returns a random element from an array
export function getRandomElement(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

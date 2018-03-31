export default class CellSelector {
  constructor(width, length) {
    this.width = width;
    this.length = length;
    this.cache = [];
    this.all = [];
    this.initCache();
    this.select = this.select.bind(this);
  }

  select(x, y) {
    if (x < 0 || x >= this.width) {
      return undefined;
    }
    if (y < 0 || y >= this.length) {
      return undefined;
    }
    return this.cache[y][x];
  }

  initCache() {
    for (let y = 0; y < this.length; y += 1) {
      this.cache[y] = [];
      for (let x = 0; x < this.width; x += 1) {
        this.cache[y][x] =
          document.querySelector(`.cell[data-x="${x}"][data-y="${y}"]`);
        this.all.push(this.cache[y][x]);
      }
    }
  }
}

export const CELL_WIDTH = 14;
export const LINE_WIDTH = 2;

export default class BoardRenderer {
  constructor(domElement, initialBoard) {
    this.initialBoard = initialBoard;
    this.domElement = domElement;
    this.ctx = domElement.getContext('2d');
    this.height = (initialBoard.length * CELL_WIDTH) + LINE_WIDTH;
    this.width = (initialBoard.width * CELL_WIDTH) + LINE_WIDTH;
  }

  setDomDimensions() {
    this.domElement.style.height = this.height;
    this.domElement.height = this.height;
    this.domElement.style.width = this.height;
    this.domElement.width = this.width;
  }

  static getRect(x, y) {
    return [
      LINE_WIDTH + (x * CELL_WIDTH),
      LINE_WIDTH + (y * CELL_WIDTH),
      CELL_WIDTH - LINE_WIDTH,
      CELL_WIDTH - LINE_WIDTH,
    ];
  }

  fillRect(x, y, color) {
    this.ctx.fillStyle = color;
    this.ctx.fillRect(...this.constructor.getRect(x, y));
    this.ctx.stroke();
  }

  drawBoard() {
    this.setDomDimensions();

    this.ctx.fillStyle = 'black';
    this.ctx.fillRect(0, 0, this.width, this.height);
    this.ctx.stroke();

    this.ctx.fillStyle = 'white';
    this.ctx.lineWidth = LINE_WIDTH / 2;
    for (let x = 0; x < this.initialBoard.width; x += 1) {
      for (let y = 0; y < this.initialBoard.length; y += 1) {
        this.ctx.fillRect(...this.constructor.getRect(x, y));
      }
    }
    this.ctx.stroke();
  }
}

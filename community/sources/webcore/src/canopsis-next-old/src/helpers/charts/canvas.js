/**
 * Promisified function for convert canvas to blob
 *
 * @param {HTMLCanvasElement} canvas
 * @param {string} [type]
 * @param {number} [quality]
 * @return {Promise<Blob>}
 */
export const canvasToBlob = (canvas, type, quality) => new Promise(
  resolve => canvas.toBlob(resolve, type, quality),
);

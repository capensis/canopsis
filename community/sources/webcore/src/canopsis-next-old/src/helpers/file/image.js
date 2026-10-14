/**
 * Get image properties by src
 *
 * @param {string} src
 * @returns {Promise<{ width: number, height: number }>}
 */
export const getImageProperties = src => new Promise((resolve) => {
  const img = document.createElement('img');

  img.onload = () => {
    const { width, height } = img;

    resolve({ width, height });
  };

  img.src = src;
});

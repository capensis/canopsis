/**
 * Replaces global ResizeObserver with a Jest mock that reports a fixed content rect
 * when `observe` runs (e.g. for components that layout from ResizeObserver callbacks).
 *
 * @param {number} width
 * @param {number} [height=28]
 */
export const installResizeObserver = (width, height = 28) => {
  global.ResizeObserver = jest.fn().mockImplementation(callback => ({
    observe: jest.fn(() => {
      callback([{ contentRect: { width, height } }]);
    }),
    unobserve: jest.fn(),
    disconnect: jest.fn(),
  }));
};

export const uninstallResizeObserver = () => {
  delete global.ResizeObserver;
};

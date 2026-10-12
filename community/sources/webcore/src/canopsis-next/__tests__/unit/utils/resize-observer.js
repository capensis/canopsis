/**
 * Replaces global ResizeObserver with a Jest mock that reports a fixed content rect
 * when `observe` runs (e.g. for components that layout from ResizeObserver callbacks).
 *
 * @param {number} width
 * @param {number} [height=28]
 */
export const installResizeObserver = (width, height = 28) => {
  global.ResizeObserver = vi.fn().mockImplementation(callback => ({
    observe: vi.fn(() => {
      callback([{ contentRect: { width, height } }]);
    }),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
  }));
};

export const uninstallResizeObserver = () => {
  delete global.ResizeObserver;
};

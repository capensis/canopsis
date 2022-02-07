/**
 * Mock for date now. Clear yourself after all tests.
 *
 * @param {number} nowTimestamp
 */
export const mockDateNow = (nowTimestamp) => {
  let dateNowSpy;

  beforeAll(() => {
    dateNowSpy = jest.spyOn(Date, 'now').mockReturnValue(nowTimestamp);
  });

  afterAll(() => {
    dateNowSpy.mockRestore();
  });
};

/**
 * Mock for requestAnimationFrame. Clear yourself after all tests.
 */
export const mockRequestAnimationFrame = () => {
  let requestAnimationFrameSpy = null;

  beforeEach(() => {
    requestAnimationFrameSpy = jest.spyOn(window, 'requestAnimationFrame')
      .mockImplementation(() => {});
  });

  afterEach(() => {
    requestAnimationFrameSpy.mockRestore();
  });
};

/**
 * Mock for date. Clear yourself after all tests.
 *
 * @param {number | Date} nowTimestamp
 */
export const mockDateGetTime = (nowTimestamp) => {
  let dateSpy;

  beforeAll(() => {
    dateSpy = jest
      .spyOn(Date.prototype, 'getTime')
      .mockReturnValue(nowTimestamp);
  });

  afterAll(() => {
    dateSpy.mockRestore();
  });
};

/**
 * Mock for the modals. Clear yourself after all tests.
 */
export const mockModals = () => {
  const modals = {
    show: jest.fn(),
    hide: jest.fn(),
    minimize: jest.fn(),
    maximize: jest.fn(),
  };

  afterEach(() => {
    modals.show.mockReset();
    modals.hide.mockReset();
    modals.minimize.mockReset();
    modals.maximize.mockReset();
  });

  return modals;
};

/**
 * Mock for the popups. Clear yourself after all tests.
 */
export const mockPopups = () => {
  const popups = {
    error: jest.fn(),
  };

  afterEach(() => {
    popups.error.mockReset();
  });

  return popups;
};

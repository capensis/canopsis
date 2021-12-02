/**
 * Stub for date now. Clear yourself after all tests.
 *
 * @param {number} nowTimestamp
 */
export const stubDateNow = (nowTimestamp) => {
  let dateNowSpy;

  beforeAll(() => {
    dateNowSpy = jest.spyOn(Date, 'now').mockReturnValue(nowTimestamp);
  });

  afterAll(() => {
    dateNowSpy.mockRestore();
  });
};

/**
 * Stub for requestAnimationFrame. Clear yourself after all tests.
 */
export const stubRequestAnimationFrame = () => {
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
 * Stub for date. Clear yourself after all tests.
 *
 * @param {number | Date} nowTimestamp
 */
export const stubDateGetTime = (nowTimestamp) => {
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

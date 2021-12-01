/**
 * Stub for date now. Clear yourself after all tests.
 *
 * @param {number} nowTimestamp
 */
export const stubDateNow = (nowTimestamp) => {
  let dateNowSpy;

  beforeAll(() => {
    dateNowSpy = jest.spyOn(Date, 'now').mockImplementation(() => nowTimestamp);
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

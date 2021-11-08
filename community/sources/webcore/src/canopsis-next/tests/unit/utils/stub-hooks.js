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

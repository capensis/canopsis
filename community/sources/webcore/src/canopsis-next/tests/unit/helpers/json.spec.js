import { stringifyJson } from '@/helpers/json';

describe('stringifyJson', () => {
  const defaultIndents = 4;
  const defaultValue = '{}';
  const validJson = { key: 'value' };
  const validJsonString = '{ "key": "value" }';
  const invalidJsonString = '{ key: "value" }';

  it('Valid json with default arguments', () => {
    expect(stringifyJson(validJson)).toBe(JSON.stringify(validJson, null, defaultIndents));
  });

  it('Valid json string with default arguments', () => {
    expect(stringifyJson(validJsonString)).toBe(JSON.stringify(JSON.parse(validJsonString), null, defaultIndents));
  });

  it('Undefined json with default arguments', () => {
    expect(stringifyJson(undefined)).toBe(defaultValue);
  });

  it('Valid json with special indents', () => {
    const indents = 2;

    expect(stringifyJson(validJson, indents)).toBe(JSON.stringify(validJson, null, indents));
  });

  it('Undefined json with special defaultValue', () => {
    expect(stringifyJson(undefined, defaultIndents, validJsonString)).toBe(validJsonString);
  });

  it('Invalid json with default arguments', () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

    expect(stringifyJson(invalidJsonString)).toBe(defaultValue);
    expect(consoleErrorSpy).toBeCalledTimes(1);

    consoleErrorSpy.mockRestore();
  });

  it('Invalid json with special default value', () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

    expect(stringifyJson(invalidJsonString, defaultIndents, validJsonString)).toBe(validJsonString);
    expect(consoleErrorSpy).toBeCalledTimes(1);

    consoleErrorSpy.mockRestore();
  });
});

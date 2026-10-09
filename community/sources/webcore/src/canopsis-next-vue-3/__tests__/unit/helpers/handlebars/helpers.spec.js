import Handlebars from 'handlebars';

import { DATETIME_FORMATS } from '@/constants';

import { nowHelper, requestHelper } from '@/helpers/handlebars/helpers';
import { uid } from '@/helpers/uid';

vi.mock('@/helpers/uid');

describe('nowHelper', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-10-07T12:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  test('Applies format without realtime', () => {
    expect(nowHelper({ hash: { format: '[formatted]' } })).toBe('formatted');
  });

  test('Applies timezone without realtime', () => {
    expect(nowHelper({ hash: { timezone: 'America/New_York' } })).toBe('07/10/2026 08:00:00');
  });

  test('Applies offset without realtime', () => {
    expect(nowHelper({ hash: { offset: '+15m' } })).toMatch(/:15:00$/);
  });

  test('Enables realtime without other options', () => {
    const result = nowHelper({ hash: { realtime: true } });

    expect(result).toBeInstanceOf(Handlebars.SafeString);
    expect(result.toString()).toContain(`format="${DATETIME_FORMATS.long}"`);
    expect(result.toString()).toMatch(/timezone="[^"]+"/);
    expect(result.toString()).not.toContain('offset=');
  });

  test('Renders a realtime component with normalized formatting attributes', () => {
    const result = nowHelper({
      hash: {
        format: 'time',
        timezone: 'Europe/Paris',
        realtime: true,
        offset: '+15m',
      },
    });

    expect(result).toBeInstanceOf(Handlebars.SafeString);
    expect(result.toString()).toBe(
      '<c-now format="time" timezone="Europe/Paris" offset="+15m"></c-now>',
    );
  });
});

describe('requestHelper', () => {
  const mockContext = { test: 'context' };
  const mockFn = vi.fn();

  beforeEach(() => {
    global.window = global.window || {};
    global.window._handlebarsRequestHelper = undefined;
    vi.clearAllMocks();
    uid.mockReturnValue('test-uid-123');
  });

  afterEach(() => {
    if (global.window) {
      delete global.window._handlebarsRequestHelper;
    }
  });

  test('Throws error when url is missing', async () => {
    const options = {
      hash: {},
      fn: mockFn,
    };

    await expect(requestHelper.call(mockContext, options)).rejects.toThrow(
      'helper {{request}}: \'url\' is required',
    );
  });

  test('Returns empty string when options.fn is not a function', async () => {
    const options = {
      hash: {
        url: 'https://api.example.com/data',
      },
      fn: null,
    };

    const result = await requestHelper.call(mockContext, options);

    expect(result).toBe('');
    expect(global.window._handlebarsRequestHelper).toBeUndefined();
  });

  test('Creates helper data in window._handlebarsRequestHelper', async () => {
    const url = 'https://api.example.com/data';
    const options = {
      hash: {
        url,
        method: 'post',
        headers: '{"Content-Type": "application/json"}',
      },
      fn: mockFn,
    };

    await requestHelper.call(mockContext, options);

    expect(uid).toHaveBeenCalledWith('request-helper');
    expect(global.window._handlebarsRequestHelper).toBeDefined();
    expect(global.window._handlebarsRequestHelper['test-uid-123']).toBeDefined();
    expect(global.window._handlebarsRequestHelper['test-uid-123'].context).toBe(mockContext);
    expect(global.window._handlebarsRequestHelper['test-uid-123'].options).toBe(options);
  });

  test('Initializes window._handlebarsRequestHelper if it does not exist', async () => {
    const options = {
      hash: {
        url: 'https://api.example.com/data',
      },
      fn: mockFn,
    };

    expect(global.window._handlebarsRequestHelper).toBeUndefined();

    await requestHelper.call(mockContext, options);

    expect(global.window._handlebarsRequestHelper).toBeDefined();
  });

  test('Preserves existing window._handlebarsRequestHelper data', async () => {
    const existingId = 'existing-id';
    const existingData = { test: 'data' };

    global.window._handlebarsRequestHelper = {
      [existingId]: existingData,
    };

    const options = {
      hash: {
        url: 'https://api.example.com/data',
      },
      fn: mockFn,
    };

    await requestHelper.call(mockContext, options);

    expect(global.window._handlebarsRequestHelper[existingId]).toBe(existingData);
    expect(global.window._handlebarsRequestHelper['test-uid-123']).toBeDefined();
  });

  test('Returns SafeString with correct component markup', async () => {
    const options = {
      hash: {
        url: 'https://api.example.com/data',
      },
      fn: mockFn,
    };

    const result = await requestHelper.call(mockContext, options);

    expect(result).toBeInstanceOf(Handlebars.SafeString);
    expect(result.toString()).toBe(
      '<c-request-helper key="test-uid-123" helper-id="test-uid-123"></c-request-helper>',
    );
  });

  test('Generates unique id for each helper call', async () => {
    const options = {
      hash: {
        url: 'https://api.example.com/data',
      },
      fn: mockFn,
    };

    uid.mockReturnValueOnce('first-id');
    const firstResult = await requestHelper.call(mockContext, options);

    uid.mockReturnValueOnce('second-id');
    const secondResult = await requestHelper.call(mockContext, options);

    expect(firstResult.toString()).toContain('first-id');
    expect(secondResult.toString()).toContain('second-id');
    expect(global.window._handlebarsRequestHelper['first-id']).toBeDefined();
    expect(global.window._handlebarsRequestHelper['second-id']).toBeDefined();
  });

  test('Stores all hash options in helper data', async () => {
    const options = {
      hash: {
        url: 'https://api.example.com/data',
        method: 'post',
        headers: '{"Content-Type": "application/json"}',
        data: '{"key": "value"}',
        path: 'data.result',
        variable: 'myVar',
        username: 'user',
        password: 'pass',
      },
      fn: mockFn,
    };

    await requestHelper.call(mockContext, options);

    expect(global.window._handlebarsRequestHelper['test-uid-123'].options.hash).toEqual(
      options.hash,
    );
  });
});

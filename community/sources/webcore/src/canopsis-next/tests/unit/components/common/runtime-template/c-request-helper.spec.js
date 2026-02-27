import Faker from 'faker';
import axios from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { RESPONSE_STATUSES } from '@/constants';

import CRequestHelper from '@/components/common/runtime-template/c-request-helper.vue';

const stubs = {
  'c-runtime-template': true,
};

const snapshotStubs = {
  'c-runtime-template': true,
};

describe('c-request-helper', () => {
  const axiosMockAdapter = new AxiosMockAdapter(axios);
  const helperId = Faker.datatype.uuid();
  const mockFn = jest.fn(context => Promise.resolve(`<div>${JSON.stringify(context)}</div>`));

  beforeEach(() => {
    global.window = global.window || {
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
    };
    global.window._handlebarsRequestHelper = global.window._handlebarsRequestHelper || {};
    global.window._handlebarsRequestHelper[helperId] = {
      options: {
        fn: mockFn,
        hash: {
          method: 'get',
          url: 'https://api.example.com/data',
        },
      },
    };
    axiosMockAdapter.reset();
    mockFn.mockClear();
  });

  afterEach(() => {
    /**
     * Ensure window._handlebarsRequestHelper exists before component cleanup tries to delete from it
     */
    if (global.window && !global.window._handlebarsRequestHelper) {
      global.window._handlebarsRequestHelper = {};
    }
    if (global.window && global.window._handlebarsRequestHelper) {
      delete global.window._handlebarsRequestHelper[helperId];
    }
    jest.clearAllMocks();
  });

  const factory = generateShallowRenderer(CRequestHelper, { stubs });
  const snapshotFactory = generateRenderer(CRequestHelper, { stubs: snapshotStubs });

  test('Does not send request when helperData is not available', async () => {
    delete global.window._handlebarsRequestHelper[helperId];

    factory({
      propsData: { helperId },
    });

    await flushPromises();

    expect(axiosMockAdapter.history.get).toHaveLength(0);
    expect(mockFn).not.toHaveBeenCalled();
  });

  test('Sends request successfully and updates template', async () => {
    const responseData = { data: { name: 'test', value: 123 } };

    axiosMockAdapter.onGet('https://api.example.com/data').reply(200, responseData);

    const wrapper = factory({
      propsData: { helperId },
    });

    await flushPromises();

    expect(axiosMockAdapter.history.get).toHaveLength(1);
    expect(mockFn).toHaveBeenCalledWith(responseData);
    expect(wrapper.vm.template).toContain('<fragment>');
  });

  test('Sends request with custom method and options', async () => {
    const postData = { key: 'value' };
    const responseData = { success: true };

    global.window._handlebarsRequestHelper[helperId].options.hash = {
      method: 'post',
      url: 'https://api.example.com/create',
      headers: JSON.stringify({ 'Content-Type': 'application/json' }),
      data: JSON.stringify(postData),
      username: 'user',
      password: 'pass',
    };

    axiosMockAdapter.onPost('https://api.example.com/create').reply(200, responseData);

    factory({
      propsData: { helperId },
    });

    await flushPromises();

    expect(axiosMockAdapter.history.post).toHaveLength(1);
    expect(axiosMockAdapter.history.post[0].data).toBe(JSON.stringify(postData));
    expect(axiosMockAdapter.history.post[0].auth).toEqual({ username: 'user', password: 'pass' });
  });

  test('Handles unauthorized error', async () => {
    axiosMockAdapter.onGet('https://api.example.com/data').reply(RESPONSE_STATUSES.unauthorized);

    const wrapper = factory({
      propsData: { helperId },
    });

    await flushPromises();

    expect(wrapper.vm.template).toContain('Unauthorized');
  });

  test('Handles timeout error', async () => {
    axiosMockAdapter.onGet('https://api.example.com/data').reply(RESPONSE_STATUSES.timeout);

    const wrapper = factory({
      propsData: { helperId },
    });

    await flushPromises();

    expect(wrapper.vm.template).toContain('Request timeout');
  });

  test('Handles other errors', async () => {
    axiosMockAdapter.onGet('https://api.example.com/data').reply(500);

    const wrapper = factory({
      propsData: { helperId },
    });

    await flushPromises();

    expect(wrapper.vm.template).toContain('Error while fetching data');
  });

  test('Cancels request on component unmount', async () => {
    axiosMockAdapter.onGet('https://api.example.com/data').reply(200, { data: 'test' });

    const wrapper = factory({
      propsData: { helperId },
    });

    wrapper.destroy();

    await flushPromises();

    expect(global.window._handlebarsRequestHelper[helperId]).toBeUndefined();
  });

  test('Renders `c-request-helper` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: { helperId },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-request-helper` after successful request', async () => {
    axiosMockAdapter.onGet('https://api.example.com/data').reply(200, { data: 'test' });

    const wrapper = snapshotFactory({
      propsData: { helperId },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});

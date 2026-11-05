import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import axios from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';

import { RESPONSE_STATUSES, VIEW_PERIODIC_REFRESH_KEY } from '@/constants';

import Observer from '@/services/observer';

import CRequestHelper from '@/components/common/runtime-template/c-request-helper.vue';

const stubs = {
  'c-runtime-template': true,
};

describe('c-request-helper', () => {
  const axiosMockAdapter = new AxiosMockAdapter(axios);
  const periodicRefresh = new Observer();
  const helperId = Faker.datatype.uuid();
  const mockFn = vi.fn(context => Promise.resolve(`<div>${JSON.stringify(context)}</div>`));
  let consoleErrorSpy;

  beforeEach(() => {
    consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    global.window = global.window || {
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
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
    consoleErrorSpy.mockRestore();
    vi.clearAllMocks();
  });

  const factory = generateShallowRenderer(CRequestHelper, {
    stubs,
    parentComponent: {
      provide: {
        [VIEW_PERIODIC_REFRESH_KEY]: periodicRefresh,
      },
    },
  });

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
    expect(wrapper.vm.template).toContain('<div class="all-inherit">');
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

  test('Fetches fresh data on each view refresh with an unchanged template', async () => {
    const orchestration = { flux: 'NOVALYS-RESEAU', status: 'on' };

    Object.assign(global.window._handlebarsRequestHelper[helperId].options.hash, {
      path: 'orchestrationList',
      variable: 'o',
    });

    axiosMockAdapter.onGet('https://api.example.com/data').reply(() => [
      200,
      { orchestrationList: [{ ...orchestration }] },
    ]);

    const wrapper = factory({ propsData: { helperId } });

    await flushPromises();

    expect(axiosMockAdapter.history.get).toHaveLength(1);
    expect(mockFn).toHaveBeenLastCalledWith({ o: [{ ...orchestration }] });
    expect(wrapper.vm.template).toContain('on');

    orchestration.status = 'off';

    await periodicRefresh.notify();
    await flushPromises();

    expect(axiosMockAdapter.history.get).toHaveLength(2);
    expect(mockFn).toHaveBeenLastCalledWith({ o: [{ ...orchestration }] });
    expect(wrapper.vm.template).toContain('off');

    await periodicRefresh.notify();

    expect(axiosMockAdapter.history.get).toHaveLength(3);
  });

  test('Stops refreshing after component unmount', async () => {
    axiosMockAdapter.onGet('https://api.example.com/data').reply(200, { data: 'test' });

    const wrapper = factory({ propsData: { helperId } });

    await flushPromises();

    expect(periodicRefresh.childrenHandlers).toHaveLength(1);

    wrapper.destroy();

    await periodicRefresh.notify();

    expect(periodicRefresh.childrenHandlers).toHaveLength(0);
    expect(axiosMockAdapter.history.get).toHaveLength(1);
    expect(global.window._handlebarsRequestHelper[helperId]).toBeUndefined();
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
});

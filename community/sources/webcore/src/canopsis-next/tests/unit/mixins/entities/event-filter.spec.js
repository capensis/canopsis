import Faker from 'faker';
import AxiosMockAdapter from 'axios-mock-adapter';

import { generateShallowRenderer } from '@unit/utils/vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import store from '@/store';

import { entitiesEventFilterMixin } from '@/mixins/entities/event-filter';

describe('Entities event filter mixin', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  const totalCount = 15;
  const eventFilters = Array.from({ length: totalCount }).map(() => ({
    _id: Faker.datatype.string(),
    name: Faker.datatype.string(),
  }));
  const meta = {
    total_count: totalCount,
  };

  const factory = generateShallowRenderer({
    render() {},
    mixins: [entitiesEventFilterMixin],
  }, { store });

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  test('Initial values applied', () => {
    const wrapper = factory();

    expect(wrapper.vm.eventFiltersMeta).toEqual({});
    expect(wrapper.vm.eventFiltersPending).toBe(false);
    expect(wrapper.vm.eventFilters).toEqual([]);
  });

  test('Event filter rules list fetched', async () => {
    axiosMockAdapter
      .onGet(API_ROUTES.eventFilter.rules)
      .reply(200, { data: eventFilters, meta });

    const wrapper = factory();

    const waitFetch = wrapper.vm.fetchEventFiltersList();

    expect(wrapper.vm.eventFiltersPending).toBe(true);

    await waitFetch;

    expect(wrapper.vm.eventFiltersPending).toBe(false);

    expect(wrapper.vm.eventFiltersMeta).toEqual(meta);
    expect(wrapper.vm.eventFilters).toEqual(eventFilters);
  });

  test('Event filter rules list fetched with params', async () => {
    const params = {
      page: Faker.datatype.string(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.eventFilter.rules, { params })
      .reply(200, { data: eventFilters, meta });

    const wrapper = factory();

    await wrapper.vm.fetchEventFiltersList({ params });

    expect(wrapper.vm.eventFiltersMeta).toEqual(meta);
    expect(wrapper.vm.eventFilters).toEqual(eventFilters);
  });

  test('Event filter rules list fetched with error', async () => {
    const error = Faker.datatype.string();
    const params = {
      page: Faker.datatype.string(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.eventFilter.rules, { params })
      .replyOnce(200, { data: eventFilters, meta });

    const wrapper = factory();

    await wrapper.vm.fetchEventFiltersList({ params });

    axiosMockAdapter
      .onGet(API_ROUTES.eventFilter.rules, { params })
      .reply(400, error);

    const originalError = console.error;
    console.error = jest.fn();

    try {
      await wrapper.vm.fetchEventFiltersList({ params });
    } catch (err) {
      expect(err).toEqual(error);
    }

    expect(console.error).toBeCalledWith(error);

    expect(wrapper.vm.eventFiltersPending).toBe(false);
    expect(wrapper.vm.eventFiltersMeta).toEqual(meta);
    expect(wrapper.vm.eventFilters).toEqual(eventFilters);

    console.error = originalError;
  });

  test('Event filter rules list fetched with previous params', async () => {
    const params = {
      page: Faker.datatype.string(),
    };
    const reversedEventFilters = eventFilters.slice().reverse();
    axiosMockAdapter
      .onGet(API_ROUTES.eventFilter.rules, { params })
      .replyOnce(200, { data: eventFilters, meta });

    const wrapper = factory();

    await wrapper.vm.fetchEventFiltersList({ params });

    axiosMockAdapter
      .onGet(API_ROUTES.eventFilter.rules, { params })
      .replyOnce(200, { data: reversedEventFilters, meta });

    await wrapper.vm.refreshEventFiltersList();

    expect(wrapper.vm.eventFiltersMeta).toEqual(meta);
    expect(wrapper.vm.eventFilters).toEqual(reversedEventFilters);
  });

  test('Event filter rule created', async () => {
    const [eventFilterRule] = eventFilters;
    axiosMockAdapter
      .onPost(API_ROUTES.eventFilter.rules)
      .replyOnce(200);

    const wrapper = factory();

    await wrapper.vm.createEventFilter({ data: eventFilterRule });

    const [postRequestData] = axiosMockAdapter.history.post;

    expect(JSON.parse(postRequestData.data)).toEqual(eventFilterRule);
  });

  test('Event filter rule updated', async () => {
    const [eventFilterRule] = eventFilters;
    axiosMockAdapter
      .onPut(`${API_ROUTES.eventFilter.rules}/${encodeURIComponent(eventFilterRule._id)}`)
      .replyOnce(200);

    const wrapper = factory();

    const updatedRule = {
      ...eventFilterRule,
      param: Faker.datatype.string(),
    };

    await wrapper.vm.updateEventFilter({ data: updatedRule, id: eventFilterRule._id });

    const [putRequestData] = axiosMockAdapter.history.put;

    expect(JSON.parse(putRequestData.data)).toEqual(updatedRule);
  });

  test('Event filter rule removed', async () => {
    const [eventFilterRule] = eventFilters;
    axiosMockAdapter
      .onDelete(`${API_ROUTES.eventFilter.rules}/${encodeURIComponent(eventFilterRule._id)}`)
      .replyOnce(200);

    const wrapper = factory();

    await wrapper.vm.removeEventFilter({ id: eventFilterRule._id });

    expect(axiosMockAdapter.history.delete).toHaveLength(1);
  });
});

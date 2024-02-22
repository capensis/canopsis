import Faker from 'faker';
import AxiosMockAdapter from 'axios-mock-adapter';

import { generateShallowRenderer } from '@unit/utils/vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import store from '@/store';

import { entitiesPatternsMixin } from '@/mixins/entities/pattern';

describe('Entities pattern mixin', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  const totalCount = 15;
  const patterns = Array.from({ length: totalCount }).map(() => ({
    _id: Faker.datatype.string(),
    name: Faker.datatype.string(),
  }));
  const meta = {
    total_count: totalCount,
  };

  const factory = generateShallowRenderer({
    render() {},
    mixins: [entitiesPatternsMixin],
  }, { store });

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  test('Initial values applied', () => {
    const wrapper = factory();

    expect(wrapper.vm.patternsMeta).toEqual({});
    expect(wrapper.vm.patternsPending).toBe(false);
    expect(wrapper.vm.patterns).toEqual([]);
  });

  test('Patterns list fetched', async () => {
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list)
      .reply(200, { data: patterns, meta });

    const wrapper = factory();

    const waitFetch = wrapper.vm.fetchPatternsList();

    expect(wrapper.vm.patternsPending).toBe(true);

    await waitFetch;

    expect(wrapper.vm.patternsPending).toBe(false);

    expect(wrapper.vm.patternsMeta).toEqual(meta);
    expect(wrapper.vm.patterns).toEqual(patterns);
  });

  test('Patterns list fetched with params', async () => {
    const params = {
      page: Faker.datatype.string(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, { params })
      .reply(200, { data: patterns, meta });

    const wrapper = factory();

    await wrapper.vm.fetchPatternsList({ params });

    expect(wrapper.vm.patternsMeta).toEqual(meta);
    expect(wrapper.vm.patterns).toEqual(patterns);
  });

  test('Patterns list fetched with error', async () => {
    const error = Faker.lorem.word();
    const params = {
      page: Faker.datatype.string(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, { params })
      .replyOnce(200, { data: patterns, meta });

    const wrapper = factory();

    await wrapper.vm.fetchPatternsList({ params });

    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, { params })
      .reply(400, error);

    const originalError = console.error;
    console.error = jest.fn();

    try {
      await wrapper.vm.fetchPatternsList({ params });
    } catch (err) {
      expect(err).toEqual(error);
    }

    expect(console.error).toBeCalledWith(error);

    expect(wrapper.vm.patternsPending).toBe(false);
    expect(wrapper.vm.patternsMeta).toEqual(meta);
    expect(wrapper.vm.patterns).toEqual(patterns);

    console.error = originalError;
  });

  test('Patterns list fetched with previous params', async () => {
    const params = {
      page: Faker.datatype.string(),
    };
    const reversedPatterns = patterns.slice().reverse();
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, { params })
      .replyOnce(200, { data: patterns, meta });

    const wrapper = factory();

    await wrapper.vm.fetchPatternsList({ params });

    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, { params })
      .replyOnce(200, { data: reversedPatterns, meta });

    await wrapper.vm.fetchPatternsListWithPreviousParams();

    expect(wrapper.vm.patternsMeta).toEqual(meta);
    expect(wrapper.vm.patterns).toEqual(reversedPatterns);
  });
});

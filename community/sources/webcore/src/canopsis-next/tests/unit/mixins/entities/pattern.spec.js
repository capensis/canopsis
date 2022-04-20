import Faker from 'faker';
import AxiosMockAdapter from 'axios-mock-adapter';

import { shallowMount, createVueInstance } from '@unit/utils/vue';

import { entitiesPatternsMixin } from '@/mixins/entities/pattern';
import store from '@/store';
import request from '@/services/request';
import { API_ROUTES } from '@/config';

const localVue = createVueInstance();

const factory = () => shallowMount({
  render() {},
  mixins: [entitiesPatternsMixin],
}, {
  localVue,
  store,
});

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
      .onGet(API_ROUTES.patterns)
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
      .onGet(API_ROUTES.patterns, { params })
      .reply(200, { data: patterns, meta });

    const wrapper = factory();

    await wrapper.vm.fetchPatternsList({ params });

    expect(wrapper.vm.patternsMeta).toEqual(meta);
    expect(wrapper.vm.patterns).toEqual(patterns);
  });
});

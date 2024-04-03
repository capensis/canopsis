import Faker from 'faker';
import AxiosMockAdapter from 'axios-mock-adapter';

import { generateShallowRenderer } from '@unit/utils/vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import store from '@/store';

import { entitiesCorporatePatternsMixin } from '@/mixins/entities/pattern/corporate';

describe('Entities corporate pattern mixin', () => {
  const axiosMockAdapter = new AxiosMockAdapter(request);

  const totalCount = 15;
  const corporatePatterns = Array.from({ length: totalCount }).map(() => ({
    _id: Faker.datatype.string(),
    name: Faker.datatype.string(),
    is_corporate: true,
  }));
  const meta = {
    total_count: totalCount,
  };

  const factory = generateShallowRenderer({
    render() {},
    mixins: [entitiesCorporatePatternsMixin],
  }, { store });

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  test('Initial values applied', () => {
    const wrapper = factory();

    expect(wrapper.vm.corporatePatternsMeta).toEqual({});
    expect(wrapper.vm.corporatePatternsPending).toBe(false);
    expect(wrapper.vm.corporatePatterns).toEqual([]);
  });

  test('Corporate patterns list fetched', async () => {
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, { patterns: { corporate: true } })
      .reply(200, { data: corporatePatterns, meta });

    const wrapper = factory();

    const waitFetch = wrapper.vm.fetchCorporatePatternsList();

    expect(wrapper.vm.corporatePatternsPending).toBe(true);

    await waitFetch;

    expect(wrapper.vm.corporatePatternsPending).toBe(false);

    expect(wrapper.vm.corporatePatternsMeta).toEqual(meta);
    expect(wrapper.vm.corporatePatterns).toEqual(corporatePatterns);
  });

  test('Corporate patterns list fetched with params', async () => {
    const params = {
      page: Faker.datatype.string(),
    };
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, {
        params: { corporate: true, ...params },
      })
      .reply(200, { data: corporatePatterns, meta });

    const wrapper = factory();

    await wrapper.vm.fetchCorporatePatternsList({ params });

    expect(wrapper.vm.corporatePatternsMeta).toEqual(meta);
    expect(wrapper.vm.corporatePatterns).toEqual(corporatePatterns);
  });

  test('Corporate patterns list fetched with error', async () => {
    const error = Faker.datatype.string();
    const params = {
      page: Faker.datatype.string(),
    };
    const options = {
      params: { corporate: true, ...params },
    };
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, options)
      .replyOnce(200, { data: corporatePatterns, meta });

    const wrapper = factory();

    await wrapper.vm.fetchCorporatePatternsList({ params });

    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, options)
      .reply(400, error);

    const originalError = console.error;
    console.error = jest.fn();

    try {
      await wrapper.vm.fetchCorporatePatternsList({ params });
    } catch (err) {
      expect(err).toEqual(error);
    }

    expect(console.error).toBeCalledWith(error);

    expect(wrapper.vm.corporatePatternsPending).toBe(false);
    expect(wrapper.vm.corporatePatternsMeta).toEqual(meta);
    expect(wrapper.vm.corporatePatterns).toEqual(corporatePatterns);

    console.error = originalError;
  });

  test('Corporate patterns list fetched with previous params', async () => {
    const params = {
      page: Faker.datatype.string(),
    };
    const reversedPatterns = corporatePatterns.slice().reverse();
    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, {
        params: { corporate: true, ...params },
      })
      .replyOnce(200, { data: corporatePatterns, meta });

    const wrapper = factory();

    await wrapper.vm.fetchCorporatePatternsList({ params });

    axiosMockAdapter
      .onGet(API_ROUTES.pattern.list, {
        params: { corporate: true, ...params },
      })
      .replyOnce(200, { data: reversedPatterns, meta });

    await wrapper.vm.fetchCorporatePatternsListWithPreviousParams();

    expect(wrapper.vm.corporatePatternsMeta).toEqual(meta);
    expect(wrapper.vm.corporatePatterns).toEqual(reversedPatterns);
  });
});

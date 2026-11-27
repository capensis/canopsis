import { flushPromises, shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { nextTick } from 'vue';

import { API_ROUTES } from '@/config';
import {
  ENTITY_INFO_PROPERTY_TYPES,
  ENTITY_PATTERN_FIELDS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PBEHAVIOR_PATTERN_FIELDS,
} from '@/constants';

import request from '@/services/request';

import { useEntityInfoPropertyStore } from '@/stores/modules/entity-info-property';

import { formToAdvancedSearch } from '@/helpers/search/advanced-search';

import CAvailabilityAdvancedSearch from '@/components/common/search/c-availability-advanced-search.vue';
import CEntityAdvancedSearch from '@/components/common/search/c-entity-advanced-search.vue';

vi.unmock('@/stores/modules/service');
vi.unmock('@/stores/modules/dynamic-info');
vi.unmock('@/stores/modules/entity-info-property');

const advancedSearchStub = {
  props: ['attributes', 'pending'],
  template: '<div />',
};

describe('c-entity-advanced-search', () => {
  let pinia;
  let wrapper;
  let requestGet;

  const factory = (component = CEntityAdvancedSearch) => {
    wrapper = shallowMount(component, {
      global: {
        plugins: [pinia],
        stubs: { 'c-advanced-search': advancedSearchStub },
      },
    });

    return wrapper.findComponent(advancedSearchStub);
  };

  beforeEach(() => {
    pinia = createPinia();
    setActivePinia(pinia);
    requestGet = vi.spyOn(request, 'get').mockResolvedValue({ data: [] });
  });

  afterEach(() => {
    wrapper?.unmount();
    vi.restoreAllMocks();
  });

  it('fetches dictionary keys and populates Infos and Component infos', async () => {
    let resolveKeys;
    const keysResponse = new Promise((resolve) => {
      resolveKeys = resolve;
    });
    requestGet.mockImplementation(url => url === API_ROUTES.entityInfosDictionaryKeys
      ? keysResponse
      : Promise.resolve({ data: [] }));

    const search = factory();

    await nextTick();

    expect(requestGet).toHaveBeenCalledWith('/internal/entity-infos-dictionary/keys', {
      params: { paginate: false },
    });
    expect(search.props('pending')).toBe(true);

    resolveKeys({ data: [{ value: 'environment', type: ENTITY_INFO_PROPERTY_TYPES.string }] });

    await flushPromises();

    expect(search.props('pending')).toBe(false);

    for (const field of [ENTITY_PATTERN_FIELDS.infos, ENTITY_PATTERN_FIELDS.componentInfos]) {
      const attribute = search.props('attributes').find(item => item.value === field);

      expect(attribute.disabled).toBe(false);
      expect(attribute.items).toEqual([
        expect.objectContaining({
          value: 'environment',
          definedType: PATTERN_FIELD_TYPES.string,
          items: [
            expect.objectContaining({ value: 'name' }),
            expect.objectContaining({ value: 'value' }),
          ],
        }),
      ]);
    }
  });

  it('passes aliases loaded by the widget to the search reactively', async () => {
    const search = factory();

    await flushPromises();

    expect(search.props('attributes').some(item => item.alias)).toBe(false);

    const property = {
      _id: 'environment',
      name: 'environment',
      alias: 'ENV',
      type: ENTITY_INFO_PROPERTY_TYPES.string,
    };
    useEntityInfoPropertyStore().items = [property];

    await nextTick();

    const alias = search.props('attributes').find(item => item.value === 'ENV');

    expect(alias).toEqual(expect.objectContaining({
      alias: true,
      text: 'ENV',
      definedType: PATTERN_FIELD_TYPES.string,
      original: property,
      disabled: false,
    }));
    expect(formToAdvancedSearch([{
      attribute: alias.value,
      alias: alias.alias,
      operator: PATTERN_OPERATORS.equal,
      fieldType: PATTERN_FIELD_TYPES.string,
      value: 'production',
    }]).search_pattern).toEqual([[{
      alias: 'ENV',
      cond: { type: 'eq', value: 'production' },
      field_type: PATTERN_FIELD_TYPES.string,
    }]]);
  });

  it('keeps the search pending while alias properties are loading', async () => {
    const properties = useEntityInfoPropertyStore();
    properties.pending = true;

    const search = factory();

    await flushPromises();

    expect(search.props('pending')).toBe(true);

    properties.pending = false;

    await nextTick();

    expect(search.props('pending')).toBe(false);
  });

  it('includes the four Pbehavior fields in their own group', async () => {
    const search = factory();

    await flushPromises();

    const attributes = search.props('attributes');
    const fields = [
      PBEHAVIOR_PATTERN_FIELDS.name,
      PBEHAVIOR_PATTERN_FIELDS.reason,
      PBEHAVIOR_PATTERN_FIELDS.type,
      PBEHAVIOR_PATTERN_FIELDS.canonicalType,
    ];
    const firstFieldIndex = attributes.findIndex(item => item.value === fields[0]);

    expect(firstFieldIndex).toBeGreaterThan(0);
    expect(attributes[firstFieldIndex - 1]).toHaveProperty('header');
    expect(attributes.slice(firstFieldIndex).map(item => item.value)).toEqual(fields);

    for (const field of fields) {
      expect(attributes.find(item => item.value === field)).toEqual(expect.objectContaining({
        disabled: false,
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      }));
      expect(formToAdvancedSearch([{
        attribute: field,
        operator: PATTERN_OPERATORS.equal,
        value: 'maintenance',
      }]).search_pattern).toEqual([[{
        field,
        cond: { type: 'eq', value: 'maintenance' },
      }]]);
    }
  });

  it('uses custom_id for the Availability ID filter', async () => {
    const search = factory(CAvailabilityAdvancedSearch);

    await flushPromises();

    const attributes = search.props('attributes');
    const id = attributes.find(item => item.value === ENTITY_PATTERN_FIELDS.customId);

    expect(id).toBeDefined();
    expect(attributes.some(item => item.value === ENTITY_PATTERN_FIELDS.id)).toBe(false);
    expect(formToAdvancedSearch([{
      attribute: id.value,
      operator: PATTERN_OPERATORS.equal,
      value: 'resource/component',
    }]).search_pattern).toEqual([[{
      field: 'custom_id',
      cond: { type: 'eq', value: 'resource/component' },
    }]]);
  });

  it.each([
    [PATTERN_OPERATORS.enabled, true],
    [PATTERN_OPERATORS.disabled, false],
  ])('supports the Availability %s filter', async (operator, value) => {
    const search = factory(CAvailabilityAdvancedSearch);

    await flushPromises();

    const enabled = search.props('attributes').find(item => item.value === ENTITY_PATTERN_FIELDS.enabled);

    expect(enabled).toEqual(expect.objectContaining({
      disabled: false,
      operators: [PATTERN_OPERATORS.enabled, PATTERN_OPERATORS.disabled],
    }));
    expect(formToAdvancedSearch([{ attribute: enabled.value, operator }]).search_pattern).toEqual([[{
      field: 'enabled',
      cond: { type: 'eq', value },
    }]]);
  });

  it('preserves the other search fields for Availability', async () => {
    const entitySearch = factory();

    await flushPromises();

    const entityFields = entitySearch.props('attributes').map(item => item.value);

    wrapper.unmount();

    const availabilitySearch = factory(CAvailabilityAdvancedSearch);

    await flushPromises();

    const availabilityFields = availabilitySearch.props('attributes').map(item => item.value);

    expect(availabilityFields).toEqual(expect.arrayContaining(entityFields.map(field =>
      field === ENTITY_PATTERN_FIELDS.id ? ENTITY_PATTERN_FIELDS.customId : field)));
    expect(availabilityFields).toHaveLength(entityFields.length + 1);
    expect(entityFields).toContain(ENTITY_PATTERN_FIELDS.id);
    expect(entityFields).not.toContain(ENTITY_PATTERN_FIELDS.enabled);
  });
});

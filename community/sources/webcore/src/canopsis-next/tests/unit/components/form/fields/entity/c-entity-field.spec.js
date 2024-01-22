import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';
import { BASIC_ENTITY_TYPES } from '@/constants';

import { PAGINATION_LIMIT } from '@/config';
import CEntityField from '@/components/forms/fields/entity/c-entity-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field';
import CLazySearchField from '@/components/forms/fields/c-lazy-search-field';

const stubs = {
  'c-lazy-search-field': CLazySearchField,
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-lazy-search-field': CLazySearchField,
  'c-select-field': CSelectField,
};

const selectAutocomplete = wrapper => wrapper.find('.c-select-field');

const selectInput = wrapper => wrapper.find('input');

describe('c-entity-field', () => {
  const items = [
    {
      _id: 'id',
      value: 'value',
      text: 'Text',
      type: 'type1',
    },
    {
      _id: 'id 2',
      value: 'value 2',
      text: 'Text 2',
      type: 'type2',
    },
    {
      _id: 'id 3',
      value: 'value 3',
      text: 'Text 3',
      type: 'type3',
    },
  ];
  const fetchListWithoutStore = jest.fn().mockReturnValue({
    data: items,
    meta: {
      page_count: items.length,
    },
  });

  const store = createMockedStoreModules([
    {
      name: 'entity',
      actions: {
        fetchListWithoutStore,
      },
    },
  ]);

  const factory = generateShallowRenderer(CEntityField, {
    stubs,
    attachTo: document.body,
  });
  const snapshotFactory = generateRenderer(CEntityField, {
    stubs: snapshotStubs,
    attachTo: document.body,
  });

  afterEach(() => {
    fetchListWithoutStore.mockClear();
  });

  it('Entities fetched after focus', async () => {
    const wrapper = factory({
      store,
      propsData: {
        itemText: 'text',
        itemValue: 'value',
      },
    });

    const autocompleteElement = selectAutocomplete(wrapper);

    autocompleteElement.trigger('focus');

    await flushPromises();

    expect(fetchListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: PAGINATION_LIMIT,
          page: 1,
          search: null,
          type: Object.values(BASIC_ENTITY_TYPES),
        },
      },
      undefined,
    );
  });

  it('Value changed after trigger the input', async () => {
    const wrapper = factory({
      store,
      propsData: {
        itemText: 'text',
        itemValue: 'value',
      },
    });

    const autocompleteElement = selectAutocomplete(wrapper);

    autocompleteElement.trigger('focus');

    await flushPromises();

    autocompleteElement.setValue(items[0].value);

    expect(wrapper).toEmit('input', items[0].value);
  });

  it('Renders `c-entity-field` with default props', async () => {
    const wrapper = snapshotFactory({
      store,
    });

    await selectInput(wrapper).trigger('focus');

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-entity-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        value: items[2].text,
        search: items[1].text,
        items,
        label: 'Custom label',
        name: 'customName',
        itemText: 'text',
        itemValue: 'value',
        disabled: true,
        loading: true,
      },
    });

    selectInput(wrapper).trigger('focus');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-entity-field` with array value', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        value: items.map(({ text }) => text),
        items,
        itemText: 'text',
        itemValue: 'value',
      },
    });

    selectInput(wrapper).trigger('focus');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});

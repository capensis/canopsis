import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';
import { MAX_LIMIT, PATTERN_CUSTOM_ITEM_VALUE, PATTERN_TYPES } from '@/constants';

import CPatternField from '@/components/forms/fields/pattern/c-pattern-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field.vue';

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-select-field': CSelectField,
};

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('c-pattern-field', () => {
  const fetchPatternsListWithoutStore = jest.fn().mockReturnValue({
    data: [],
  });
  const patternModule = {
    name: 'pattern',
    actions: {
      fetchListWithoutStore: fetchPatternsListWithoutStore,
    },
  };
  const store = createMockedStoreModules([
    patternModule,
  ]);
  const patterns = [
    { title: 'Pattern 1', _id: 'pattern-value-1' },
    { title: 'Pattern 2', _id: 'pattern-value-2' },
    { title: 'Pattern 3', _id: 'pattern-value-3' },
  ];

  const factory = generateShallowRenderer(CPatternField, { stubs });
  const snapshotFactory = generateRenderer(CPatternField, { stubs: snapshotStubs });

  afterEach(() => {
    fetchPatternsListWithoutStore.mockClear();
  });

  test('Patterns fetched after mount', async () => {
    factory({ store });

    await flushPromises();

    expect(fetchPatternsListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        params: { limit: MAX_LIMIT },
      },
      undefined,
    );
  });

  test('Pbehavior patterns fetched after mount', async () => {
    factory({
      store,
      propsData: {
        type: PATTERN_TYPES.pbehavior,
      },
    });

    await flushPromises();

    expect(fetchPatternsListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        params: { limit: MAX_LIMIT, type: PATTERN_TYPES.pbehavior },
      },
      undefined,
    );
  });

  test('Value changed after trigger the select', () => {
    const wrapper = factory({ store });
    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input', PATTERN_CUSTOM_ITEM_VALUE);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(PATTERN_CUSTOM_ITEM_VALUE);
  });

  test('Renders `c-pattern-field` with default props', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-pattern-field` with custom props', () => {
    fetchPatternsListWithoutStore.mockResolvedValueOnce({
      data: patterns,
    });
    const wrapper = snapshotFactory({
      propsData: {
        value: patterns[0],
        label: 'Custom label',
        name: 'customName',
        disabled: true,
        required: true,
        returnObject: true,
        type: PATTERN_TYPES.alarm,
      },
      store: createMockedStoreModules([
        patternModule,
      ]),
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});

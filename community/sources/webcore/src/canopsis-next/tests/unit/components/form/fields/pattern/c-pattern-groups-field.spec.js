import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { PATTERN_FIELD_TYPES, QUICK_RANGES, TIME_UNITS } from '@/constants';

import CPatternGroupsField from '@/components/forms/fields/pattern/c-pattern-groups-field.vue';

const localVue = createVueInstance();

const stubs = {
  'c-pattern-group-field': true,
  'c-pattern-operator-chip': true,
};

const factory = (options = {}) => shallowMount(CPatternGroupsField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternGroupsField, {
  localVue,
  stubs,

  ...options,
});

const selectAddButton = wrapper => wrapper.find('v-btn-stub');
const selectPatternGroupsField = wrapper => wrapper.findAll('c-pattern-group-field-stub');
const selectPatternGroupFieldByIndex = (wrapper, index) => selectPatternGroupsField(wrapper)
  .at(index);

describe('c-pattern-groups-field', () => {
  const groups = [
    { rules: [], key: 'key 1' },
    { rules: [], key: 'key 2' },
    { rules: [], key: 'key 3' },
  ];

  test('Group removed after trigger remove event on the pattern group field', () => {
    const wrapper = factory({
      propsData: {
        groups,
        attributes: [],
      },
    });

    const secondGroup = selectPatternGroupFieldByIndex(wrapper, 1);

    secondGroup.vm.$emit('remove');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      groups[0],
      groups[2],
    ]);
  });

  test('Group updated after trigger update event on the pattern group field', () => {
    const wrapper = factory({
      propsData: {
        groups,
        attributes: [],
      },
    });

    const lastGroup = selectPatternGroupFieldByIndex(wrapper, 2);

    const updatedGroup = {
      rules: [{}],
      key: 'new key',
    };

    lastGroup.vm.$emit('input', updatedGroup);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      groups[0],
      groups[1],
      updatedGroup,
    ]);
  });

  test('Group added after click on the add button', () => {
    const attributeItem = {
      value: 'test',
    };
    const wrapper = factory({
      propsData: {
        groups,
        attributes: [
          attributeItem,
        ],
      },
    });

    const addButton = selectAddButton(wrapper);

    addButton.vm.$emit('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      ...groups,
      {
        key: expect.any(String),
        rules: [
          {
            attribute: attributeItem.value,
            dictionary: '',
            field: '',
            fieldType: PATTERN_FIELD_TYPES.string,
            operator: '',
            value: '',
            range: {
              type: QUICK_RANGES.last1Hour.value,
              from: 0,
              to: 0,
            },
            key: expect.any(String),
            duration: {
              unit: TIME_UNITS.second,
              value: 1,
            },
          },
        ],
      },
    ]);
  });

  test('Renders `c-pattern-groups-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        groups: [],
        attributes: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-pattern-groups-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        groups,
        attributes: [
          { text: 'Attribute text', value: 'attribute value' },
        ],
        required: true,
        disabled: true,
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-pattern-groups-field` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        groups: [],
        attributes: [],
        required: true,
      },
    });

    const validator = wrapper.getValidator();
    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });
});

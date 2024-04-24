import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { PATTERN_FIELD_TYPES, QUICK_RANGES, TIME_UNITS } from '@/constants';

import PatternGroupsField from '@/components/forms/fields/pattern/pattern-groups-field.vue';

const stubs = {
  'pattern-group-field': true,
  'c-pattern-operator-chip': true,
  'c-btn-with-error': true,
};

const selectAddButton = wrapper => wrapper.find('c-btn-with-error-stub');
const selectPatternGroupsField = wrapper => wrapper.findAll('pattern-group-field-stub');
const selectPatternGroupFieldByIndex = (wrapper, index) => selectPatternGroupsField(wrapper)
  .at(index);

describe('pattern-groups-field', () => {
  const factory = generateShallowRenderer(PatternGroupsField, { stubs });
  const snapshotFactory = generateRenderer(PatternGroupsField, { stubs });

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

    secondGroup.triggerCustomEvent('remove');

    expect(wrapper).toEmitInput([
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

    lastGroup.triggerCustomEvent('input', updatedGroup);

    expect(wrapper).toEmitInput([
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

    selectAddButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmitInput([
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

  test('Renders `pattern-groups-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        groups: [],
        attributes: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-groups-field` with custom props', () => {
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

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-groups-field` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        groups: [],
        attributes: [],
        required: true,
      },
    });

    const validator = wrapper.getValidator();
    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});

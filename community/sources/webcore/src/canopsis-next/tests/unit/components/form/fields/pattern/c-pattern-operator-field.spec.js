import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { PATTERN_OPERATORS } from '@/constants';

import CPatternOperatorField from '@/components/forms/fields/pattern/c-pattern-operator-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field';

const localVue = createVueInstance();

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-select-field': CSelectField,
};

const factory = (options = {}) => shallowMount(CPatternOperatorField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternOperatorField, {
  localVue,
  attachTo: document.body,
  stubs: snapshotStubs,

  ...options,
});

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('c-pattern-operator-field', () => {
  test('Value changed after trigger the select', () => {
    const wrapper = factory({
      propsData: {
        value: PATTERN_OPERATORS.notExist,
      },
    });
    const selectField = selectSelectField(wrapper);

    selectField.setValue(PATTERN_OPERATORS.notExist);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(PATTERN_OPERATORS.notExist);
  });

  test('Renders `c-pattern-operator-field` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  test('Renders `c-pattern-operator-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: PATTERN_OPERATORS.exist,
        label: 'Custom label',
        name: 'customName',
        disabled: true,
        required: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

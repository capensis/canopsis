import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { FILTER_OPERATORS } from '@/constants';

import CPatternOperatorField from '@/components/forms/fields/filter/c-pattern-operator-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CPatternOperatorField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternOperatorField, {
  localVue,

  ...options,
});

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-pattern-operator-field', () => {
  test('Value changed after trigger the select', () => {
    const wrapper = factory({
      propsData: {
        value: FILTER_OPERATORS.notExist,
      },
    });
    const selectField = selectSelectField(wrapper);

    selectField.setValue(FILTER_OPERATORS.notExist);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(FILTER_OPERATORS.notExist);
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
        value: FILTER_OPERATORS.exist,
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

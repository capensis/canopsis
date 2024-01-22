import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { PATTERN_OPERATORS } from '@/constants';

import PatternOperatorField from '@/components/forms/fields/pattern/pattern-operator-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field';

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-select-field': CSelectField,
};

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('pattern-operator-field', () => {
  const factory = generateShallowRenderer(PatternOperatorField, { stubs });
  const snapshotFactory = generateRenderer(PatternOperatorField, {
    attachTo: document.body,
    stubs: snapshotStubs,
  });

  test('Value changed after trigger the select', () => {
    const wrapper = factory({
      propsData: {
        value: PATTERN_OPERATORS.notExist,
      },
    });

    selectSelectField(wrapper).setValue(PATTERN_OPERATORS.notExist);

    expect(wrapper).toEmit('input', PATTERN_OPERATORS.notExist);
  });

  test('Renders `pattern-operator-field` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  test('Renders `pattern-operator-field` with custom props', () => {
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

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { ALARM_PATTERN_FIELDS } from '@/constants';

import PatternAttributeField from '@/components/forms/fields/pattern/pattern-attribute-field.vue';

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-select-field': true,
};

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('pattern-attribute-field', () => {
  const factory = generateShallowRenderer(PatternAttributeField, { stubs });
  const snapshotFactory = generateRenderer(PatternAttributeField, { stubs: snapshotStubs });

  it('Value changed after trigger the input', () => {
    const value = {
      value: ALARM_PATTERN_FIELDS.ack,
      text: 'Text',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const selectElement = selectSelectField(wrapper);

    selectElement.triggerCustomEvent('input', value);

    expect(wrapper).toEmit('input', value);
  });

  it('Renders `pattern-attribute-field` with default props', () => {
    const value = {
      value: ALARM_PATTERN_FIELDS.component,
      text: 'Component',
    };
    const wrapper = snapshotFactory({
      propsData: {
        value,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `pattern-attribute-field` with custom props', () => {
    const value = {
      value: ALARM_PATTERN_FIELDS.component,
      text: 'Component',
    };
    const wrapper = snapshotFactory({
      propsData: {
        value,
        label: 'Custom label',
        items: [
          {
            value: ALARM_PATTERN_FIELDS.ackAt,
            text: 'Ack at',
          },
        ],
        name: 'custom_filter_attribute_name',
        disabled: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

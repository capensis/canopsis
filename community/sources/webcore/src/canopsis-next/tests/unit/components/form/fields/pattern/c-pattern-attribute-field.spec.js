import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ALARM_PATTERN_FIELDS } from '@/constants';

import CPatternAttributeField from '@/components/forms/fields/pattern/c-pattern-attribute-field.vue';

const localVue = createVueInstance();

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-select-field': true,
};

const factory = (options = {}) => shallowMount(CPatternAttributeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternAttributeField, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('c-pattern-attribute-field', () => {
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

    selectElement.vm.$emit('input', value);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(value);
  });

  it('Renders `c-pattern-attribute-field` with default props', () => {
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

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-pattern-attribute-field` with custom props', () => {
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

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

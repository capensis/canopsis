import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ALARM_FILTER_FIELDS } from '@/constants';

import CPatternAttributeField from '@/components/forms/fields/pattern/c-pattern-attribute-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CPatternAttributeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternAttributeField, {
  localVue,

  ...options,
});

describe('c-pattern-attribute-field', () => {
  it('Value changed after trigger the input', () => {
    const value = {
      value: ALARM_FILTER_FIELDS.ack,
      text: 'Text',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.vm.$emit('input', value);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(value);
  });

  it('Renders `c-pattern-attribute-field` with default props', () => {
    const value = {
      value: ALARM_FILTER_FIELDS.component,
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
      value: ALARM_FILTER_FIELDS.component,
      text: 'Component',
    };
    const wrapper = snapshotFactory({
      propsData: {
        value,
        label: 'Custom label',
        items: [
          {
            value: ALARM_FILTER_FIELDS.ackAt,
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

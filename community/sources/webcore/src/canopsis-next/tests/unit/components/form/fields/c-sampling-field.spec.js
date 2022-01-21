import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { SAMPLINGS } from '@/constants';
import CSamplingField from '@/components/forms/fields/c-sampling-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CSamplingField, {
  localVue,
  stubs,
  ...options,
});

const snapshotFactory = (options = {}) => mount(CSamplingField, {
  localVue,
  ...options,
});

describe('c-sampling-field', () => {
  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: SAMPLINGS.day,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.setValue(SAMPLINGS.hour);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(SAMPLINGS.hour);
  });

  it('Renders `c-sampling-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: SAMPLINGS.day,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-sampling-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: SAMPLINGS.hour,
        label: 'Custom label',
        name: 'customName',
        disabled: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

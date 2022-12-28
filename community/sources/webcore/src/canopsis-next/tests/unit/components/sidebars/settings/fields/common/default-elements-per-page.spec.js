import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import DefaultElementsPerPage from '@/components/sidebars/settings/fields/common/default-elements-per-page.vue';

const localVue = createVueInstance();

const stubs = {
  'widget-settings-item': true,
  'v-select': createSelectInputStub('v-select'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const factory = (options = {}) => shallowMount(DefaultElementsPerPage, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(DefaultElementsPerPage, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectElementsPerPageField = wrapper => wrapper.find('select.v-select');

describe('default-elements-per-page', () => {
  it('Value changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: 20,
      },
    });

    const elementsPerPageField = selectElementsPerPageField(wrapper);

    const newValue = 10;

    elementsPerPageField.vm.$emit('input', newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(newValue);
  });

  it('Renders `default-sort-column` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `default-sort-column` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 50,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

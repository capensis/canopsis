import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import DefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';

const stubs = {
  'widget-settings-item': true,
  'v-select': createSelectInputStub('v-select'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectElementsPerPageField = wrapper => wrapper.find('select.v-select');

describe('default-elements-per-page', () => {
  const factory = generateShallowRenderer(DefaultElementsPerPage, { stubs });
  const snapshotFactory = generateRenderer(DefaultElementsPerPage, { stubs: snapshotStubs });

  it('Value changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: 20,
      },
    });

    const newValue = 10;

    selectElementsPerPageField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  it('Renders `default-sort-column` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `default-sort-column` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 50,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import DefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';

const stubs = {
  'widget-settings-item': true,
  'v-select': createSelectInputStub('v-select'),
};

const selectElementsPerPageField = wrapper => wrapper.find('select.v-select');

describe('default-elements-per-page', () => {
  const factory = generateShallowRenderer(DefaultElementsPerPage, { stubs });

  it('Value changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: 20,
      },
    });

    const newValue = 10;

    selectElementsPerPageField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});

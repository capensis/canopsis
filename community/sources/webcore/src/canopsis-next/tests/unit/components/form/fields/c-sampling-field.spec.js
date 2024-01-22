import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { SAMPLINGS } from '@/constants';
import CSamplingField from '@/components/forms/fields/c-sampling-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('c-sampling-field', () => {
  const factory = generateShallowRenderer(CSamplingField, { stubs });
  const snapshotFactory = generateRenderer(CSamplingField);

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: SAMPLINGS.day,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.setValue(SAMPLINGS.hour);

    expect(wrapper).toEmit('input', SAMPLINGS.hour);
  });

  it('Renders `c-sampling-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: SAMPLINGS.day,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});

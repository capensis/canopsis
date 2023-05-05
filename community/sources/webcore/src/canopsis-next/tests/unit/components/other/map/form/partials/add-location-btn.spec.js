import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import AddLocationBtn from '@/components/other/map/form/partials/add-location-btn.vue';

const factory = generateShallowRenderer(AddLocationBtn, {

});

const snapshotFactory = generateRenderer(AddLocationBtn, {

});

const selectBtnToggleNode = wrapper => wrapper.vm.$children[0];

describe('add-location-btn', () => {
  test('Value updated after click button', () => {
    const wrapper = factory({
      propsData: {
        value: true,
      },
    });

    const btnToggleNode = selectBtnToggleNode(wrapper);

    btnToggleNode.$emit('change');

    expect(wrapper).toEmit('input', false);
  });

  test('Renders `add-location-btn` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: false,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});

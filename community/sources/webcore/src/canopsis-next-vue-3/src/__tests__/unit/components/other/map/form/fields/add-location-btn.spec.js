import { generateShallowRenderer } from '@unit/utils/vue';

import AddLocationBtn from '@/components/other/map/form/fields/add-location-btn.vue';

const selectBtnToggleNode = wrapper => wrapper.vm.$children[0];

describe('add-location-btn', () => {
  const factory = generateShallowRenderer(AddLocationBtn);

  test('Value updated after click button', () => {
    const wrapper = factory({
      propsData: {
        value: true,
      },
    });

    const btnToggleNode = selectBtnToggleNode(wrapper);

    btnToggleNode.$emit('change');

    expect(wrapper).toEmitInput(false);
  });
});

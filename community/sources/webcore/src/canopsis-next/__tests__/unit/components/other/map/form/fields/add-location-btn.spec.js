import { generateShallowRenderer } from '@unit/utils/vue';

import AddLocationBtn from '@/components/other/map/form/fields/add-location-btn.vue';

const stubs = {
  'v-btn-toggle': {
    name: 'VBtnToggle',
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: '<div><slot /></div>',
  },
  'v-btn': true,
  'v-icon': true,
};

const selectBtnToggle = wrapper => wrapper.findComponent({ name: 'VBtnToggle' });

describe('add-location-btn', () => {
  const factory = generateShallowRenderer(AddLocationBtn, { stubs });

  test('Value updated after click button', async () => {
    const wrapper = factory({
      props: {
        value: true,
      },
    });

    expect(selectBtnToggle(wrapper).props('modelValue')).toBe(true);

    selectBtnToggle(wrapper).vm.$emit('update:modelValue', undefined);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('input')).toEqual([[false]]);
  });
});

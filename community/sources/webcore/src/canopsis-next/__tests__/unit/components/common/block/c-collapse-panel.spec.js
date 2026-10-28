import { generateShallowRenderer } from '@unit/utils/vue';

import CCollapsePanel from '@/components/common/block/c-collapse-panel.vue';

const VExpansionPanelsStub = {
  props: ['modelValue'],
  emits: ['update:modelValue'],
  template: '<div><button class="toggle" @click="$emit(\'update:modelValue\', modelValue === 0 ? undefined : 0)" /><slot /></div>',
};

const VIconStub = {
  props: ['icon'],
  template: '<i class="icon">{{ icon }}</i>',
};

describe('c-collapse-panel', () => {
  const factory = generateShallowRenderer(CCollapsePanel, {
    stubs: {
      'v-expansion-panels': VExpansionPanelsStub,
      'v-expansion-panel': { template: '<div><slot /></div>' },
      'v-expansion-panel-title': { template: '<div><slot /><slot name="actions" /></div>' },
      'v-expansion-panel-text': { template: '<div><slot /></div>' },
      'v-card': { template: '<div><slot /></div>' },
      'v-card-text': { template: '<div><slot /></div>' },
      'v-icon': VIconStub,
    },
  });

  it('updates the arrow when the accordion state changes', async () => {
    const wrapper = factory();

    expect(wrapper.find('.icon').text()).toBe('expand_more');

    await wrapper.find('.toggle').trigger('click');

    expect(wrapper.find('.icon').text()).toBe('expand_less');

    await wrapper.find('.toggle').trigger('click');

    expect(wrapper.find('.icon').text()).toBe('expand_more');
  });
});

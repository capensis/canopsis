import { generateShallowRenderer } from '@unit/utils/vue';

import MermaidPointMarker from '@/components/other/map/partials/mermaid-point-marker.vue';

const stubs = {
  'point-icon': {
    name: 'PointIcon',
    emits: ['click', 'contextmenu', 'dblclick', 'mousedown'],
    template: '<div />',
  },
};

const selectPointIcon = wrapper => wrapper.findComponent({ name: 'PointIcon' });

describe('mermaid-point-marker', () => {
  const factory = generateShallowRenderer(MermaidPointMarker, { stubs });

  test('Listeners applied to point icon', async () => {
    const wrapper = factory({
      props: {
        x: 1,
        y: 2,
      },
    });
    const event = new Event('click');
    const pointIcon = selectPointIcon(wrapper);

    pointIcon.vm.$emit('click', event);
    pointIcon.vm.$emit('contextmenu', event);
    pointIcon.vm.$emit('dblclick', event);
    pointIcon.vm.$emit('mousedown', event);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('click')).toEqual([[event]]);
    expect(wrapper.emitted('contextmenu')).toEqual([[event]]);
    expect(wrapper.emitted('dblclick')).toEqual([[event]]);
    expect(wrapper.emitted('mousedown')).toEqual([[event]]);
  });
});

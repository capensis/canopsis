import { generateShallowRenderer } from '@unit/utils/vue';

import MermaidPointMarker from '@/components/other/map/partials/mermaid-point-marker.vue';

const stubs = {
  'point-icon': true,
};

const selectPointIconNode = wrapper => wrapper.vm.$children[0];

describe('mermaid-point-marker', () => {
  const factory = generateShallowRenderer(MermaidPointMarker, { stubs });

  test('Listeners applied to point icon', () => {
    const click = vi.fn();
    const wrapper = factory({
      propsData: {
        x: 1,
        y: 2,
      },
      listeners: {
        click,
      },
    });

    const pointIconNode = selectPointIconNode(wrapper);

    pointIconNode.$emit('click');

    expect(click).toBeCalled();
  });
});

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import MermaidPointMarker from '@/components/other/map/partials/mermaid-point-marker.vue';
import { COLOR_INDICATOR_TYPES } from '@/constants';

const stubs = {
  'point-icon': true,
};

const selectPointIconNode = wrapper => wrapper.vm.$children[0];

describe('mermaid-point-marker', () => {
  const factory = generateShallowRenderer(MermaidPointMarker, { stubs });
  const snapshotFactory = generateRenderer(MermaidPointMarker, { stubs });

  test('Listeners applied to point icon', () => {
    const click = jest.fn();
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

  test('Renders `mermaid-point-marker` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        x: 1,
        y: 2,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mermaid-point-marker` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        x: 1,
        y: 21,
        entity: {
          state: 1,
        },
        size: 15,
        colorIndicator: COLOR_INDICATOR_TYPES.state,
        pbehaviorEnabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

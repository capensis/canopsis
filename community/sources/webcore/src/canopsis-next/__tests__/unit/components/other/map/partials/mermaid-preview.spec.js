import { generateShallowRenderer } from '@unit/utils/vue';

import { mermaidPointToForm } from '@/helpers/entities/map/form';

import MermaidPreview from '@/components/other/map/partials/mermaid-preview.vue';

const stubs = {
  panzoom: true,
  'c-zoom-overlay': true,
  'mermaid-code-preview': true,
  'mermaid-points-preview': true,
};

const selectMermaidPointsPreview = wrapper => wrapper.find('mermaid-points-preview-stub');

describe('mermaid-preview', () => {
  const map = {
    name: 'Map',
    parameters: {
      code: 'code',
      points: [mermaidPointToForm({ x: 1, y: 2 }), mermaidPointToForm({ x: 10, y: 20 })],
    },
  };

  const factory = generateShallowRenderer(MermaidPreview, { stubs });

  test('Show map emitted after trigger mermaid points preview', async () => {
    const wrapper = factory({
      propsData: {
        map,
      },
    });

    const mermaidPointsPreview = selectMermaidPointsPreview(wrapper);

    const linkedMap = { _id: 'map' };
    mermaidPointsPreview.triggerCustomEvent('show:map', linkedMap);

    expect(wrapper).toEmit('show:map', linkedMap);
  });

  test('Show alarms emitted after trigger mermaid points preview', async () => {
    const wrapper = factory({
      propsData: {
        map,
      },
    });

    const mermaidPointsPreview = selectMermaidPointsPreview(wrapper);

    const linkedPoint = mermaidPointToForm({ x: 1, y: 2 });
    mermaidPointsPreview.triggerCustomEvent('show:alarms', linkedPoint);

    expect(wrapper).toEmit('show:alarms', linkedPoint);
  });

  test('Panzoom reset after map prop updated', async () => {
    const reset = vi.fn();
    const wrapper = factory({
      propsData: {
        map,
      },
      stubs: {
        ...stubs,
        panzoom: {
          template: '<div />',
          methods: {
            reset,
          },
        },
      },
    });

    await wrapper.setProps({
      map: {
        ...map,
      },
    });

    expect(reset).toBeCalled();
  });
});

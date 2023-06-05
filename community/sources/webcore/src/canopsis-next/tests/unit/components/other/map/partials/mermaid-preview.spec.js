import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';

import MermaidPreview from '@/components/other/map/partials/mermaid-preview.vue';
import { COLOR_INDICATOR_TYPES } from '@/constants';
import { mermaidPointToForm } from '@/helpers/forms/map';

const localVue = createVueInstance();

const stubs = {
  panzoom: true,
  'c-zoom-overlay': true,
  'mermaid-code-preview': true,
  'mermaid-points-preview': true,
};

const factory = (options = {}) => shallowMount(MermaidPreview, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(MermaidPreview, {
  localVue,
  stubs,

  ...options,
});

const selectMermaidPointsPreview = wrapper => wrapper.find('mermaid-points-preview-stub');

describe('mermaid-preview', () => {
  const map = {
    name: 'Map',
    parameters: {
      code: 'code',
      points: [
        mermaidPointToForm({ x: 1, y: 2 }),
        mermaidPointToForm({ x: 10, y: 20 }),
      ],
    },
  };

  test('Show map emitted after trigger mermaid points preview', async () => {
    const wrapper = factory({
      propsData: {
        map,
      },
    });

    const mermaidPointsPreview = selectMermaidPointsPreview(wrapper);

    const linkedMap = { _id: 'map' };
    mermaidPointsPreview.vm.$emit('show:map', linkedMap);

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
    mermaidPointsPreview.vm.$emit('show:alarms', linkedPoint);

    expect(wrapper).toEmit('show:alarms', linkedPoint);
  });

  test('Panzoom reset after map prop updated', async () => {
    const reset = jest.fn();
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

  test('Renders `mermaid-preview` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        map,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `mermaid-preview` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        map,
        popupTemplate: 'template',
        popupActions: true,
        colorIndicator: COLOR_INDICATOR_TYPES.state,
        pbehaviorEnabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';

import Flowchart from '@/components/common/flowchart/flowchart.vue';
import { circleShapeToForm, rectShapeToForm } from '@/helpers/flowchart/shapes';

const localVue = createVueInstance();

const stubs = {
  'flowchart-sidebar': true,
  'flowchart-editor': true,
  'flowchart-properties': true,
};

const factory = (options = {}) => shallowMount(Flowchart, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Flowchart, {
  localVue,
  stubs,

  ...options,
});

const selectFlowchartSidebar = wrapper => wrapper.find('flowchart-sidebar-stub');
const selectFlowchartEditor = wrapper => wrapper.find('flowchart-editor-stub');
const selectFlowchartProperties = wrapper => wrapper.find('flowchart-properties-stub');

describe('flowchart', () => {
  test('Shapes added after trigger flowchart sidebar', () => {
    const wrapper = factory({
      propsData: {
        shapes: {},
      },
    });

    const newShapes = {
      rect: rectShapeToForm({ _id: 'rect' }),
    };
    const flowchartSidebar = selectFlowchartSidebar(wrapper);
    flowchartSidebar.vm.$emit('input', newShapes);

    expect(wrapper).toEmit('input', newShapes);
  });

  test('Shapes added after trigger flowchart editor', () => {
    const wrapper = factory({
      propsData: {
        shapes: {},
      },
    });

    const newShapes = {
      rect: rectShapeToForm({ _id: 'rect' }),
    };
    const flowchartEditor = selectFlowchartEditor(wrapper);
    flowchartEditor.vm.$emit('input', newShapes);

    expect(wrapper).toEmit('input', newShapes);
  });

  test('Shapes added after trigger flowchart properties', async () => {
    const shapes = {
      rect: rectShapeToForm({ _id: 'rect' }),
    };
    const wrapper = factory({
      propsData: {
        shapes,
      },
    });

    const flowchartSidebar = selectFlowchartSidebar(wrapper);
    await flowchartSidebar.vm.$emit('update:selected', [shapes.rect]);

    const newShapes = {
      circle: circleShapeToForm({ _id: 'circle' }),
    };
    const flowchartProperties = selectFlowchartProperties(wrapper);
    flowchartProperties.vm.$emit('input', newShapes);

    expect(wrapper).toEmit('input', newShapes);
  });

  test('Points updated after trigger flowchart properties', async () => {
    const shapes = {
      rect: rectShapeToForm({ _id: 'rect' }),
    };
    const wrapper = factory({
      propsData: {
        shapes,
      },
    });

    const flowchartSidebar = selectFlowchartSidebar(wrapper);
    await flowchartSidebar.vm.$emit('update:selected', [shapes.rect]);

    const newShapes = {
      circle: circleShapeToForm({ _id: 'circle' }),
    };
    const flowchartProperties = selectFlowchartProperties(wrapper);
    flowchartProperties.vm.$emit('input', newShapes);

    expect(wrapper).toEmit('input', newShapes);
  });

  test('Renders `flowchart` with form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `flowchart` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes: {
          rect: rectShapeToForm({ _id: 'rect' }),
          circle: circleShapeToForm({ _id: 'circle' }),
        },
        cursorStyle: 'pointer',
        backgroundColor: '#000',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `flowchart` with readonly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes: {
          rect: rectShapeToForm({ x: 1, y: 2, _id: 'rect' }),
          circle: circleShapeToForm({ x: 1, y: 2, _id: 'circle' }),
        },
        readonly: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});

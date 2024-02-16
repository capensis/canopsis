import Faker from 'faker';
import { keyBy } from 'lodash';
import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { MERMAID_THEMES, SHAPES } from '@/constants';
import { flowchartPointToForm } from '@/helpers/entities/map/form';

import FlowchartEditor from '@/components/other/map/form/fields/flowchart-editor.vue';
import { shapeToForm } from '@/helpers/flowchart/shapes';

const stubs = {
  flowchart: {
    props: ['shapes'],
    template: `
    <div
      class="flowchart"
      :shapes="shapes"
    >
      <slot
        name="sidebar-prepend"
        :data="shapes"
      />
      <slot
        name="layers"
        :data="shapes"
      />
    </div>`,
  },
  'add-location-btn': true,
  'flowchart-points-editor': true,
};

const selectFlowchart = wrapper => wrapper.find('.flowchart');
const selectAddLocationBtn = wrapper => wrapper.find('add-location-btn-stub');
const selectFlowchartPointsEditor = wrapper => wrapper.find('flowchart-points-editor-stub');

describe('flowchart-editor', () => {
  const rectShape = shapeToForm({ type: SHAPES.rect, _id: 'rect' });
  const lineShape = shapeToForm({ type: SHAPES.rect, _id: 'line' });

  const initialForm = {
    code: '',
    theme: MERMAID_THEMES.base,
    points: [],
    shapes: {},
  };

  const factory = generateShallowRenderer(FlowchartEditor, { stubs });
  const snapshotFactory = generateRenderer(FlowchartEditor, {
    stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('Background changed after trigger flowchart', () => {
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const newColor = Faker.internet.color();

    const flowchartCodeEditor = selectFlowchart(wrapper);

    flowchartCodeEditor.vm.$emit('update:backgroundColor', newColor);

    expect(wrapper).toEmit('input', {
      ...initialForm,
      background_color: newColor,
    });
  });

  test('Shapes changed after trigger flowchart', () => {
    const shapes = keyBy([rectShape, lineShape], '_id');

    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const flowchart = selectFlowchart(wrapper);

    flowchart.vm.$emit('input', shapes);

    expect(wrapper).toEmit('input', {
      ...initialForm,
      shapes,
    });
  });

  test('Linked point removed when shape removed', () => {
    const shapes = keyBy([rectShape, lineShape], '_id');

    const firstPoint = { _id: 'point-1', x: 12, y: 32 };
    const secondPoint = {
      _id: 'point-2',
      shape: rectShape._id,
      entity: {},
    };
    const points = [firstPoint, secondPoint];

    const wrapper = factory({
      propsData: {
        form: {
          ...initialForm,
          points,
          shapes,
        },
      },
    });

    const newShapes = { [lineShape._id]: lineShape };

    const flowchart = selectFlowchart(wrapper);

    flowchart.vm.$emit('input', newShapes);

    expect(wrapper).toEmit('input', {
      ...initialForm,
      points: [firstPoint],
      shapes: newShapes,
    });
  });

  test('Points changed after trigger points editor', () => {
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const newPoints = [flowchartPointToForm()];

    const flowchartPointsEditor = selectFlowchartPointsEditor(wrapper);

    flowchartPointsEditor.vm.$emit('input', newPoints);

    expect(wrapper).toEmit('input', {
      ...initialForm,
      points: newPoints,
    });
  });

  test('Add on click mode enabled after trigger add location btn', async () => {
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const addLocationBtn = selectAddLocationBtn(wrapper);

    await addLocationBtn.vm.$emit('input', true);

    const flowchartPointsEditor = selectFlowchartPointsEditor(wrapper);

    expect(flowchartPointsEditor.vm.addOnClick).toBeTruthy();
  });

  test('Form re-validated after change form with error', async () => {
    const wrapper = factory({
      propsData: {
        name: 're-validate-name',
        form: initialForm,
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(validator.errors.items).toHaveLength(1);

    wrapper.setProps({
      form: {
        ...initialForm,
        code: 'flowchart-code',
        points: [flowchartPointToForm()],
      },
    });

    await flushPromises();

    expect(validator.errors.items).toHaveLength(0);
  });

  test('Renders `flowchart-editor` with form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `flowchart-editor` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
        name: 'custom_name',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `flowchart-editor` with validation errors ', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
        name: 'custom_name',
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});

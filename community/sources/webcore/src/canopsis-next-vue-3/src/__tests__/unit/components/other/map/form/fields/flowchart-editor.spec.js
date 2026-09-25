import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { keyBy } from 'lodash';

import { MERMAID_THEMES, SHAPES } from '@/constants';

import { flowchartPointToForm } from '@/helpers/entities/map/form';
import { shapeToForm } from '@/helpers/flowchart/shapes';

import FlowchartEditor from '@/components/other/map/form/fields/flowchart-editor.vue';

const stubs = {
  flowchart: {
    name: 'Flowchart',
    props: ['shapes'],
    emits: ['update:shapes', 'update:backgroundColor'],
    template: `
    <div
      class="flowchart"
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
  'add-location-btn': {
    name: 'AddLocationBtn',
    props: ['value'],
    emits: ['input'],
    template: '<div />',
  },
  'flowchart-points-editor': {
    name: 'FlowchartPointsEditor',
    props: ['points', 'shapes', 'addOnClick'],
    emits: ['input'],
    template: '<div />',
  },
  'map-editor-required-field': {
    name: 'MapEditorRequiredField',
    props: ['valid', 'message', 'name'],
    template: '<div><slot :has-error="!valid" :messages="[message]" /></div>',
  },
  'v-divider': true,
  'v-messages': true,
};

const selectFlowchart = wrapper => wrapper.findComponent({ name: 'Flowchart' });
const selectAddLocationBtn = wrapper => wrapper.findComponent({ name: 'AddLocationBtn' });
const selectFlowchartPointsEditor = wrapper =>
  wrapper.findComponent({ name: 'FlowchartPointsEditor' });
const selectRequiredField = wrapper => wrapper.findComponent({ name: 'MapEditorRequiredField' });

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

  test('Background changed after trigger flowchart', async () => {
    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    const newColor = Faker.internet.color();

    selectFlowchart(wrapper).vm.$emit('update:backgroundColor', newColor);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...initialForm,
      background_color: newColor,
    }]]);
  });

  test('Shapes changed after trigger flowchart', async () => {
    const shapes = keyBy([rectShape, lineShape], '_id');

    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    selectFlowchart(wrapper).vm.$emit('update:shapes', shapes);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...initialForm,
      shapes,
    }]]);
  });

  test('Linked point removed when shape removed', async () => {
    const shapes = keyBy([rectShape, lineShape], '_id');

    const firstPoint = { _id: 'point-1', x: 12, y: 32 };
    const secondPoint = {
      _id: 'point-2',
      shape: rectShape._id,
      entity: {},
    };
    const points = [firstPoint, secondPoint];

    const wrapper = factory({
      props: {
        form: {
          ...initialForm,
          points,
          shapes,
        },
      },
    });

    const newShapes = { [lineShape._id]: lineShape };

    selectFlowchart(wrapper).vm.$emit('update:shapes', newShapes);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...initialForm,
      points: [firstPoint],
      shapes: newShapes,
    }]]);
  });

  test('Points changed after trigger points editor', async () => {
    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    const newPoints = [flowchartPointToForm()];

    selectFlowchartPointsEditor(wrapper).vm.$emit('input', newPoints);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...initialForm,
      points: newPoints,
    }]]);
  });

  test('Add on click mode enabled after trigger add location btn', async () => {
    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    selectAddLocationBtn(wrapper).vm.$emit('input', true);
    await wrapper.vm.$nextTick();

    expect(selectFlowchartPointsEditor(wrapper).props('addOnClick')).toBe(true);
  });

  test('Form re-validated after change form with error', async () => {
    const wrapper = factory({
      props: {
        name: 're-validate-name',
        form: initialForm,
      },
    });

    expect(selectRequiredField(wrapper).props('valid')).toBe(false);

    await wrapper.setProps({
      form: {
        ...initialForm,
        code: 'flowchart-code',
        points: [flowchartPointToForm()],
      },
    });

    await flushPromises();

    expect(selectRequiredField(wrapper).props('valid')).toBe(true);
  });
});

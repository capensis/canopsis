import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES, ENTITIES_STATES, PBEHAVIOR_TYPE_TYPES, SHAPES } from '@/constants';

import { shapeToForm } from '@/helpers/flowchart/shapes';
import { getImpactStateColor } from '@/helpers/entities/entity/color';

import FlowchartPreview from '@/components/other/map/partials/flowchart-preview.vue';

const stubs = {
  panzoom: true,
  flowchart: {
    props: ['shapes'],
    template: `
    <div
      class="flowchart"
      :shapes="shapes"
    >
      <slot
        name="layers"
        :data="shapes"
      />
      <slot/>
    </div>`,
  },
  'flowchart-points-preview': true,
  'c-help-icon': true,
};

const selectFlowchart = wrapper => wrapper.find('div.flowchart');

describe('flowchart-preview', () => {
  const rectShape = shapeToForm({ type: SHAPES.rect, _id: 'rect' });
  const lineShape = shapeToForm({ type: SHAPES.rect, _id: 'line' });
  const circleShape = shapeToForm({ type: SHAPES.circle, _id: 'circle' });
  const shapes = [rectShape, lineShape, circleShape];

  const firstPoint = { _id: 'point-1', x: 12, y: 32 };
  const secondPoint = {
    _id: 'point-2',
    x: 12,
    y: 32,
    entity: {
      pbehavior_info: {
        canonical_type: PBEHAVIOR_TYPE_TYPES.active,
      },
    },
  };
  const thirdPoint = {
    _id: 'point-3',
    shape: lineShape._id,
    entity: {
      pbehavior_info: {
        canonical_type: PBEHAVIOR_TYPE_TYPES.pause,
      },
    },
  };
  const fourthPoint = {
    _id: 'point-4',
    shape: circleShape._id,
    entity: {
      impact_state: 2,
      state: {
        val: ENTITIES_STATES.ok,
      },
      pbehavior_info: {
        canonical_type: PBEHAVIOR_TYPE_TYPES.active,
      },
    },
  };
  const points = [firstPoint, secondPoint, thirdPoint, fourthPoint];

  const map = {
    name: 'Map',
    parameters: {
      shapes,
      code: 'code',
      background_color: 'rgba(0, 231, 123, 0.9)',
      points,
    },
  };

  const factory = generateShallowRenderer(FlowchartPreview, { stubs });
  const snapshotFactory = generateRenderer(FlowchartPreview, { stubs });

  test('Shapes color changed by point with color indicator impact state', async () => {
    const wrapper = factory({
      propsData: {
        map,
        popupActions: true,
        colorIndicator: COLOR_INDICATOR_TYPES.state,
        pbehaviorEnabled: true,
      },
    });

    const flowchart = selectFlowchart(wrapper);

    const okColorDarken = '#000000';

    expect(flowchart.vm.shapes).toEqual({
      [lineShape._id]: lineShape,
      [rectShape._id]: rectShape,
      [circleShape._id]: {
        ...circleShape,
        properties: {
          ...circleShape.properties,
          fill: '',
          stroke: okColorDarken,
        },
        textProperties: {
          ...circleShape.textProperties,
          color: okColorDarken,
        },
      },
    });
  });

  test('Shapes color changed by point with color indicator state', async () => {
    const wrapper = factory({
      propsData: {
        map,
        popupActions: true,
        colorIndicator: COLOR_INDICATOR_TYPES.impactState,
        pbehaviorEnabled: true,
      },
    });

    const flowchart = selectFlowchart(wrapper);

    const colorDarken = '#89951a';

    expect(flowchart.vm.shapes).toEqual({
      [lineShape._id]: lineShape,
      [rectShape._id]: rectShape,
      [circleShape._id]: {
        ...circleShape,
        properties: {
          ...circleShape.properties,
          fill: getImpactStateColor(2),
          stroke: colorDarken,
        },
        textProperties: {
          ...circleShape.textProperties,
          color: colorDarken,
        },
      },
    });
  });

  test('Renders `flowchart-preview` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        map,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `flowchart-preview` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        map,
        popupTemplate: 'template',
        popupActions: true,
        colorIndicator: COLOR_INDICATOR_TYPES.state,
        pbehaviorEnabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

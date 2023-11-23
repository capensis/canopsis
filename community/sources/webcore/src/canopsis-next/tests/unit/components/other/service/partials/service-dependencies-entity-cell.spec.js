import { generateRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES, ENTITY_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS, ENTITY_TYPES } from '@/constants';

import { getWidgetColumnLabel } from '@/helpers/entities/widget/list';

import ServiceDependenciesEntityCell from '@/components/other/service/partials/service-dependencies-entity-cell.vue';

const stubs = {
  'color-indicator-wrapper': true,
  'c-alarm-chip': true,
};

describe('service-dependencies', () => {
  const item = {
    entity: {
      _id: 'data-alarm-2-entity',
      name: 'Data alarm 2 entity',
      type: ENTITY_TYPES.service,
      state: 1,
      impact_level: 1,
      impact_state: 3,
      has_impacts: false,
    },
  };
  const column = {
    value: `entity.${ENTITY_FIELDS.name}`,
    text: getWidgetColumnLabel({ value: ENTITY_FIELDS.name }, ENTITY_FIELDS_TO_LABELS_KEYS),
  };
  const columnWithStateColorIndicator = {
    ...column,
    colorIndicator: COLOR_INDICATOR_TYPES.state,
  };
  const columnWithImpactStateColorIndicator = {
    ...column,
    colorIndicator: COLOR_INDICATOR_TYPES.impactState,
  };
  const stateColumn = {
    value: `entity.${ENTITY_FIELDS.state}`,
    text: getWidgetColumnLabel({ value: ENTITY_FIELDS.state }, ENTITY_FIELDS_TO_LABELS_KEYS),
  };
  const stateColumnWithStateColorIndicator = {
    ...stateColumn,
    colorIndicator: COLOR_INDICATOR_TYPES.state,
  };
  const stateColumnWithImpactStateColorIndicator = {
    ...stateColumn,
    colorIndicator: COLOR_INDICATOR_TYPES.impactState,
  };

  const snapshotFactory = generateRenderer(ServiceDependenciesEntityCell, { stubs });

  test('Renders `service-dependencies-entity-cell` with column', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        item,
        column,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-dependencies-entity-cell` with column with state color indicator', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        item,
        column: columnWithStateColorIndicator,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-dependencies-entity-cell` with column with impact state color indicator', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        item,
        column: columnWithImpactStateColorIndicator,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-dependencies-entity-cell` with state column', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        item,
        column: stateColumn,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-dependencies-entity-cell` with state column with state color indicator', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        item,
        column: stateColumnWithStateColorIndicator,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-dependencies-entity-cell` with state column with impact state color indicator', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        item,
        column: stateColumnWithImpactStateColorIndicator,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import { MAP_TYPES } from '@/constants';

import TreeOfDependenciesPreview from '@/components/other/map/partials/tree-of-dependencies-preview.vue';

const stubs = {
  'c-zoom-overlay': true,
  'c-help-icon': true,
};

const waitTicks = async ticksCount => range(ticksCount)
  .reduce(acc => acc.then(jest.runOnlyPendingTimers()), Promise.resolve());

describe('tree-of-dependencies-preview', () => {
  const map = {
    name: 'Map',
    type: MAP_TYPES.treeOfDependencies,
    parameters: {
      type: MAP_TYPES.treeOfDependencies,
      entities: [
        {
          entity: {
            _id: 'root-entity-id',
            name: 'root-entity-name',
            enabled: true,
            infos: {},
            type: 'service',
            impact_level: 1,
            category: {},
            ok_events: 0,
            ko_events: 1385055,
            state: 3,
            impact_state: 3,
            status: 1,
            alarm_last_update_date: 1705907028,
            depends_count: 606,
            entity_pattern: [],
          },
          pinned_entities: range(8).map(index => ({
            _id: `pinned-entity-id-${index}`,
            name: `pinned-entity-name-${index}`,
            enabled: true,
            infos: {},
            type: 'resource',
            impact_level: 1,
            category: null,
            last_event_date: 1701435382,
            connector: 'connector',
            component: 'component',
            ok_events: 0,
            ko_events: 0,
            state: 0,
            impact_state: 0,
            status: 0,
            depends_count: 0,
          })),
        },
      ],
    },
  };

  const snapshotFactory = generateRenderer(TreeOfDependenciesPreview, { stubs });

  test('Renders `tree-of-dependencies-preview` with required props', async () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation();

    jest.useFakeTimers();

    const rect = {
      width: 1000,
      height: 1000,
      left: 0,
      top: 0,
    };

    jest.spyOn(HTMLDivElement.prototype, 'getBoundingClientRect').mockReturnValue(rect);
    jest.spyOn(HTMLDivElement.prototype, 'clientHeight', 'get').mockReturnValue(1000);
    jest.spyOn(HTMLDivElement.prototype, 'clientWidth', 'get').mockReturnValue(1000);
    jest.spyOn(window, 'getComputedStyle').mockReturnValue({
      getPropertyValue: () => 0,
    });

    const wrapper = snapshotFactory({
      propsData: {
        map,
      },
    });

    await waitTicks(100);

    const nodeCanvas = wrapper.find('canvas[data-id="layer2-node"]');

    expect(wrapper).toMatchSnapshot();

    expect(nodeCanvas.element).toMatchCanvasSnapshot();

    jest.useRealTimers();

    consoleWarnSpy.mockReset();
  });
});

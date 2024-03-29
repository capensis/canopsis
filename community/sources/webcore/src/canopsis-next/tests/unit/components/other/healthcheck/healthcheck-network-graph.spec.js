import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import { HEALTHCHECK_ENGINES_NAMES, HEALTHCHECK_SERVICES_NAMES } from '@/constants';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/healthcheck-network-graph.vue';

const stubs = {
  'c-zoom-overlay': true,
  'c-help-icon': true,
};

const waitTicks = async ticksCount => range(ticksCount)
  .reduce(acc => acc.then(jest.runOnlyPendingTimers()), Promise.resolve());

describe('healthcheck-network-graph', () => {
  const snapshotFactory = generateRenderer(HealthcheckNetworkGraph, { stubs });

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

  test('Renders `healthcheck-network-graph` with default props', async () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation();

    jest.useFakeTimers();

    const wrapper = snapshotFactory({
      propsData: {
        // services: [],
        // enginesGraph: {},
        // enginesParameters: {},
      },
    });

    await waitTicks(100);

    const nodeCanvas = wrapper.find('canvas[data-id="layer2-node"]');

    expect(wrapper).toMatchSnapshot();

    expect(nodeCanvas.element).toMatchCanvasSnapshot();

    jest.useRealTimers();

    consoleWarnSpy.mockReset();
  });

  test('Renders `healthcheck-network-graph` with custom props', async () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation();

    jest.useFakeTimers();

    const defaultEngine = {
      instances: 1,
      min_instances: 1,
      optimal_instances: 1,
      queue_length: 0,
      time: 1711693029,
      is_running: true,
      is_queue_overflown: false,
      is_too_few_instances: false,
      is_diff_instances_config: false,
    };

    const wrapper = snapshotFactory({
      propsData: {
        services: [
          {
            name: HEALTHCHECK_SERVICES_NAMES.mongo,
            is_running: true,
          },
          {
            name: HEALTHCHECK_SERVICES_NAMES.redis,
            is_running: true,
          },
          {
            name: HEALTHCHECK_SERVICES_NAMES.rabbit,
            is_running: true,
          },
          {
            name: HEALTHCHECK_SERVICES_NAMES.timescaleDB,
            is_running: true,
          },
        ],
        enginesGraph: {
          nodes: [
            HEALTHCHECK_ENGINES_NAMES.fifo,
            HEALTHCHECK_ENGINES_NAMES.che,
            HEALTHCHECK_ENGINES_NAMES.axe,
            HEALTHCHECK_ENGINES_NAMES.correlation,
            HEALTHCHECK_ENGINES_NAMES.remediation,
            HEALTHCHECK_ENGINES_NAMES.pbehavior,
            HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
            HEALTHCHECK_ENGINES_NAMES.action,
            HEALTHCHECK_ENGINES_NAMES.webhook,
          ],
          edges: [
            {
              from: HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
              to: HEALTHCHECK_ENGINES_NAMES.action,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.che,
              to: HEALTHCHECK_ENGINES_NAMES.axe,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.fifo,
              to: HEALTHCHECK_ENGINES_NAMES.che,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.axe,
              to: HEALTHCHECK_ENGINES_NAMES.correlation,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.correlation,
              to: HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.axe,
              to: HEALTHCHECK_ENGINES_NAMES.pbehavior,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.axe,
              to: HEALTHCHECK_ENGINES_NAMES.remediation,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.action,
              to: HEALTHCHECK_ENGINES_NAMES.webhook,
            },
          ],
        },
        enginesParameters: {
          [HEALTHCHECK_ENGINES_NAMES.action]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.axe]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.che]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.correlation]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.dynamicInfos]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.fifo]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.pbehavior]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.remediation]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.webhook]: defaultEngine,
        },
      },
    });

    await waitTicks(100);

    const nodeCanvas = wrapper.find('canvas[data-id="layer2-node"]');

    expect(wrapper).toMatchSnapshot();

    expect(nodeCanvas.element).toMatchCanvasSnapshot();

    jest.useRealTimers();

    consoleWarnSpy.mockReset();
  });

  test('Renders `healthcheck-network-graph` with custom props and snmp node', async () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation();

    jest.useFakeTimers();

    const defaultEngine = {
      instances: 1,
      min_instances: 1,
      optimal_instances: 1,
      queue_length: 0,
      time: 1711693029,
      is_running: true,
      is_queue_overflown: false,
      is_too_few_instances: false,
      is_diff_instances_config: false,
    };

    const wrapper = snapshotFactory({
      propsData: {
        services: [
          {
            name: HEALTHCHECK_SERVICES_NAMES.mongo,
            is_running: true,
          },
          {
            name: HEALTHCHECK_SERVICES_NAMES.redis,
            is_running: true,
          },
          {
            name: HEALTHCHECK_SERVICES_NAMES.rabbit,
            is_running: true,
          },
          {
            name: HEALTHCHECK_SERVICES_NAMES.timescaleDB,
            is_running: true,
          },
        ],
        enginesGraph: {
          nodes: [
            HEALTHCHECK_ENGINES_NAMES.fifo,
            HEALTHCHECK_ENGINES_NAMES.che,
            HEALTHCHECK_ENGINES_NAMES.axe,
            HEALTHCHECK_ENGINES_NAMES.correlation,
            HEALTHCHECK_ENGINES_NAMES.remediation,
            HEALTHCHECK_ENGINES_NAMES.pbehavior,
            HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
            HEALTHCHECK_ENGINES_NAMES.action,
            HEALTHCHECK_ENGINES_NAMES.webhook,
            HEALTHCHECK_ENGINES_NAMES.snmp,
          ],
          edges: [
            {
              from: HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
              to: HEALTHCHECK_ENGINES_NAMES.action,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.che,
              to: HEALTHCHECK_ENGINES_NAMES.axe,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.fifo,
              to: HEALTHCHECK_ENGINES_NAMES.che,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.axe,
              to: HEALTHCHECK_ENGINES_NAMES.correlation,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.correlation,
              to: HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.axe,
              to: HEALTHCHECK_ENGINES_NAMES.pbehavior,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.axe,
              to: HEALTHCHECK_ENGINES_NAMES.remediation,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.action,
              to: HEALTHCHECK_ENGINES_NAMES.webhook,
            },
            {
              from: HEALTHCHECK_ENGINES_NAMES.snmp,
              to: HEALTHCHECK_ENGINES_NAMES.fifo,
            },
          ],
        },
        enginesParameters: {
          [HEALTHCHECK_ENGINES_NAMES.snmp]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.action]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.axe]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.che]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.correlation]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.dynamicInfos]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.fifo]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.pbehavior]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.remediation]: defaultEngine,
          [HEALTHCHECK_ENGINES_NAMES.webhook]: defaultEngine,
        },
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

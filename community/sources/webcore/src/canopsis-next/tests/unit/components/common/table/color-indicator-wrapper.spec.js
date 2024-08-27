import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES, ENTITIES_STATES } from '@/constants';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';

const stubs = {};

describe('color-indicator-wrapper', () => {
  const snapshotFactory = generateRenderer(ColorIndicatorWrapper, {

    stubs,
    attachTo: document.body,
  });

  it('Renders `color-indicator-wrapper` with state type and default slot', async () => {
    snapshotFactory({
      propsData: {
        entity: {},
        alarm: {
          v: {
            state: {
              val: ENTITIES_STATES.major,
            },
          },
        },
        type: COLOR_INDICATOR_TYPES.state,
      },
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` without type but with default slot', async () => {
    snapshotFactory({
      propsData: {
        entity: {},
        alarm: {
          v: {
            state: {
              val: ENTITIES_STATES.critical,
            },
          },
        },
        type: COLOR_INDICATOR_TYPES.state,
      },
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` with unresolved type', async () => {
    snapshotFactory({
      propsData: {
        entity: {},
        alarm: {},
        type: 'unresolved-type',
      },
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` with entity impact level and impact state', async () => {
    snapshotFactory({
      propsData: {
        entity: {
          impact_level: 2,
          impact_state: 12,
        },
        type: COLOR_INDICATOR_TYPES.impactState,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` with entity impact level and alarm state', async () => {
    snapshotFactory({
      propsData: {
        entity: {
          impact_level: 2,
        },
        alarm: {
          v: {
            state: {
              val: ENTITIES_STATES.minor,
            },
          },
        },
        type: COLOR_INDICATOR_TYPES.impactState,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});

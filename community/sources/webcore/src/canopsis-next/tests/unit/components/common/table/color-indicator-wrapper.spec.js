import { mount, createVueInstance } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES, ENTITIES_STATES } from '@/constants';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';

const localVue = createVueInstance();

const stubs = {};

const snapshotFactory = (options = {}) => mount(ColorIndicatorWrapper, {
  localVue,
  stubs,

  ...options,
});

describe('color-indicator-wrapper', () => {
  it('Renders `color-indicator-wrapper` with state type and default slot', () => {
    const wrapper = snapshotFactory({
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

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` without type but with default slot', () => {
    const wrapper = snapshotFactory({
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

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` with unresolved type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {},
        alarm: {},
        type: 'unresolved-type',
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` with entity impact level and impact state', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {
          impact_level: 2,
          impact_state: 12,
        },
        type: COLOR_INDICATOR_TYPES.impactState,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `color-indicator-wrapper` with entity impact level and alarm state', () => {
    const wrapper = snapshotFactory({
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

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});

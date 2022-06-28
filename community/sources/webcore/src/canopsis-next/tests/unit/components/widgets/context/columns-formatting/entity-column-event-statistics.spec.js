import { mount, createVueInstance } from '@unit/utils/vue';
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import EntityColumnEventStatistics from '@/components/widgets/context/columns-formatting/entity-column-event-statistics.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(EntityColumnEventStatistics, {
  localVue,
  attachTo: document.body,

  ...options,
});

describe('entity-column-event-statistics', () => {
  it('Renders `entity-column-event-statistics` with default entity', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {
          ok_events: 15,
          ko_events: 25,
        },
      },
    });

    const tooltip = wrapper.findTooltip();

    expect(tooltip.element).toMatchSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `entity-column-event-statistics` with pbehavior', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {
          ok_events: 30,
          ko_events: 45,
          pbehavior_info: {
            canonical_type: PBEHAVIOR_TYPE_TYPES.inactive,
          },
        },
      },
    });

    const tooltip = wrapper.findTooltip();

    expect(tooltip.element).toMatchSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });
});

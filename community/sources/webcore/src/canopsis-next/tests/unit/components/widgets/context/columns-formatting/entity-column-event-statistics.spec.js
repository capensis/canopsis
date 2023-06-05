import flushPromises from 'flush-promises';

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
  it('Renders `entity-column-event-statistics` with default entity', async () => {
    snapshotFactory({
      propsData: {
        entity: {
          ok_events: 15,
          ko_events: 25,
        },
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `entity-column-event-statistics` with pbehavior', async () => {
    snapshotFactory({
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

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});

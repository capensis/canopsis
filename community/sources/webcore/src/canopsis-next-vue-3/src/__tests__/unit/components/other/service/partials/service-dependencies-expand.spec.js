import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ENTITY_TYPES } from '@/constants';

import ServiceDependenciesExpand from '@/components/other/service/partials/service-dependencies-expand.vue';

const stubs = {
  'v-tooltip': createActivatorElementStub('v-tooltip'),
};

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('service-dependencies-expand', () => {
  const item = {
    entity: {
      _id: 'data-alarm-2-entity',
      name: 'Data alarm 2 entity',
      type: ENTITY_TYPES.service,
      state: 1,
      impact_level: 1,
      impact_state: 0,
    },
  };

  const itemWithLoadMore = {
    ...item,

    loadMore: true,
  };

  const factory = generateShallowRenderer(ServiceDependenciesExpand, { stubs });

  it('Show triggered after trigger click', async () => {
    const wrapper = factory({
      propsData: {
        item,
      },
    });

    const button = selectButton(wrapper);

    button.triggerCustomEvent('click');

    expect(wrapper).toEmit('show', item);
  });

  it('Load triggered after trigger click', async () => {
    const wrapper = factory({
      propsData: {
        item: itemWithLoadMore,
      },
    });

    const button = selectButton(wrapper);

    button.triggerCustomEvent('click');

    expect(wrapper).toEmit('load', itemWithLoadMore);
  });
});

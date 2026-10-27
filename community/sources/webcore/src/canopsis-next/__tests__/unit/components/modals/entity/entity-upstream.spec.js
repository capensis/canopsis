import { generateShallowRenderer } from '@unit/utils/vue';

import { MODALS } from '@/constants';

import modalComponents from '@/components/modals';
import EntityUpstreamModal from '@/components/modals/entity/entity-upstream.vue';

const stubs = {
  'c-modal-wrapper': {
    template: '<div><slot name="title" /><slot name="text" /></div>',
  },
  'entity-upstream-network-graph': {
    name: 'EntityUpstreamNetworkGraph',
    props: ['entity'],
    template: '<div />',
  },
};

const selectNetworkGraph = wrapper =>
  wrapper.findComponent({ name: 'EntityUpstreamNetworkGraph' });

describe('entity-upstream modal', () => {
  const factory = generateShallowRenderer(EntityUpstreamModal, { stubs });

  test('Modal is registered', () => {
    expect(modalComponents[MODALS.entityUpstream]).toBeDefined();
  });

  test('Network graph is rendered with modal configuration', () => {
    const entity = {
      _id: 'entity-id',
      name: 'Entity name',
      upstream: 'parent-id',
    };
    const title = 'Top level entities';
    const wrapper = factory({
      props: {
        modal: {
          config: { entity, title },
        },
      },
    });

    expect(wrapper.text()).toContain(title);
    expect(selectNetworkGraph(wrapper).props('entity')).toEqual(entity);
  });
});

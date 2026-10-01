import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ENTITY_TYPES, MODALS } from '@/constants';

import TreeOfDependenciesPreview from '@/components/other/map/partials/tree-of-dependencies-preview.vue';

const on = vi.fn();
const off = vi.fn();

const stubs = {
  'c-zoom-overlay': {
    template: '<div><slot /></div>',
  },
  'network-graph': {
    name: 'NetworkGraph',
    setup(props, { expose }) {
      expose({
        $cy: { on, off },
      });

      return {};
    },
    template: '<div />',
  },
  'c-help-icon': true,
};

describe('tree-of-dependencies-preview', () => {
  const $modals = mockModals();
  const entity = {
    _id: 'root-entity',
    name: 'Root entity',
    type: ENTITY_TYPES.service,
    depends_count: 2,
  };
  const map = {
    parameters: {
      entities: [
        {
          entity,
          pinned_entities: [
            {
              _id: 'dependency-entity',
              name: 'Dependency entity',
              type: ENTITY_TYPES.resource,
            },
          ],
        },
      ],
    },
  };

  const factory = generateShallowRenderer(TreeOfDependenciesPreview, { stubs });

  test('Entity label stays readable on the light graph background', () => {
    const wrapper = factory({
      props: { map, colorIndicator: 'state' },
    });
    const nodeTemplate = wrapper.vm.nodeHtmlLabelsOptions.find(({ query }) => query === 'node');
    const container = document.createElement('div');

    container.innerHTML = nodeTemplate.tpl({ entity });

    expect(container.querySelector('.position-absolute').style.color).toBe('rgb(0, 0, 0)');
  });

  test('Dependencies modal showed after trigger nested show all button content', () => {
    const wrapper = factory({
      props: { map },
      mocks: { $modals },
    });
    const showAllTemplate = wrapper.vm.nodeHtmlLabelsOptions.find(
      ({ query }) => query === 'node[showAll]',
    );
    const container = document.createElement('div');

    container.innerHTML = showAllTemplate.tpl({ entity });

    const showAllButton = container.querySelector('.tree-of-dependencies__show-all-btn');
    const buttonContent = container.querySelector('.v-btn__content');

    expect(showAllButton.classList.contains('bg-grey-lighten-3')).toBe(true);

    wrapper.vm.tapHandler({
      target: {
        data: () => ({}),
      },
      originalEvent: {
        target: buttonContent,
      },
    });

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.entityDependenciesList,
      config: {
        entity,
        impact: false,
      },
    });
  });
});

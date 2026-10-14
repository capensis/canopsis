import { generateShallowRenderer } from '@unit/utils/vue';

import { MODALS } from '@/constants';

import modalComponents from '@/components/modals';
import EntityDependenciesListModal from '@/components/modals/entity/entity-dependencies-list.vue';

const stubs = {
  'c-modal-wrapper': {
    template: '<div><slot name="title" /><slot name="text" /></div>',
  },
  'entity-dependencies-list-component': {
    name: 'EntityDependenciesListComponent',
    props: ['entityId', 'widget', 'columns', 'impact'],
    template: '<div />',
  },
};

const selectDependenciesList = wrapper =>
  wrapper.findComponent({ name: 'EntityDependenciesListComponent' });

describe('entity-dependencies-list modal', () => {
  const factory = generateShallowRenderer(EntityDependenciesListModal, { stubs });

  test('Modal is registered', () => {
    expect(modalComponents[MODALS.entityDependenciesList]).toBeDefined();
  });

  test('Dependencies list is rendered with modal configuration', () => {
    const entity = {
      _id: 'entity-id',
      name: 'Entity name',
    };
    const widget = {
      parameters: {
        widgetColumns: [{ value: 'name' }],
      },
    };
    const wrapper = factory({
      props: {
        modal: {
          config: {
            entity,
            widget,
            impact: true,
          },
        },
      },
    });

    expect(wrapper.text()).toContain(`${entity.name} impacted entities`);
    expect(selectDependenciesList(wrapper).props()).toEqual(
      expect.objectContaining({
        entityId: entity._id,
        widget,
        columns: widget.parameters.widgetColumns,
        impact: true,
      }),
    );
  });
});

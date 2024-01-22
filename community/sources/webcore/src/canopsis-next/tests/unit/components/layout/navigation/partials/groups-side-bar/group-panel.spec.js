import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import GroupsPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';

const stubs = {
  'v-expansion-panel-content': {
    template: `
      <div class="v-expansion-panel-content">
        <slot name="header" />
      </div>
    `,
  },
};

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('group-panel', () => {
  const factory = generateShallowRenderer(GroupsPanel, { stubs });
  const snapshotFactory = generateRenderer(GroupsPanel, {
    parentComponent: {
      provide: {
        expansionPanels: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
    },
  });

  it('Change event emitted after trigger button', () => {
    const wrapper = factory({
      propsData: {
        editable: true,
        group: {},
      },
    });

    selectButton(wrapper).triggerCustomEvent('click', new Event('click'));

    expect(wrapper).toEmit('change');
  });

  it('Renders `group-panel` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        group: {
          title: 'Group title',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `group-panel` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        group: {
          title: 'Custom group title',
        },
        editable: true,
        isEditing: true,
        orderChanged: true,
        hideActions: true,
      },
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

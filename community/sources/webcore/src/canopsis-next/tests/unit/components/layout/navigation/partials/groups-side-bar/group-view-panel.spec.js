import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import GroupsPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';

const selectButtons = wrapper => wrapper.findAll('v-btn-stub');
const selectEditButton = wrapper => selectButtons(wrapper).at(0);
const selectDuplicateButton = wrapper => selectButtons(wrapper).at(1);

describe('group-view-panel', () => {
  const factory = generateShallowRenderer(GroupsPanel);
  const snapshotFactory = generateRenderer(GroupsPanel);

  it('Change event emitted after trigger edit button', () => {
    const wrapper = factory({
      propsData: {
        view: {},
        editable: true,
      },
    });

    selectEditButton(wrapper).triggerCustomEvent('click', new Event('click'));

    expect(wrapper).toEmit('change');
  });

  it('Duplicate event emitted after trigger duplicate button', () => {
    const wrapper = factory({
      propsData: {
        view: {},
        editable: true,
        duplicable: true,
      },
    });

    selectDuplicateButton(wrapper).triggerCustomEvent('click', new Event('click'));

    expect(wrapper).toEmit('duplicate');
  });

  it('Renders `group-view-panel` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        view: {
          title: 'View title',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `group-view-panel` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        view: {
          title: 'Custom view title',
        },
        editable: true,
        duplicable: true,
        isOrderChanged: true,
        isViewActive: true,
      },
      slots: {
        title: '<div class="title-slot">TITLE</div>',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

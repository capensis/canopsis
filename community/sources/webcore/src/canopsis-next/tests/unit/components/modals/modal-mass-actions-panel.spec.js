import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { PORTALS_NAMES } from '@/constants';

import ModalMassActionsPanel from '@/components/modals/modal-mass-actions-panel.vue';

const stubs = {
  'portal-target': true,
};

const snapshotStubs = {
  'portal-target': {
    props: ['name'],
    template: `
      <div class="portal-target">
        Portal target name: {{ name }}
      </div>
    `,
  },
};

const selectPortalTarget = wrapper => wrapper.find('portal-target-stub');

describe('modal-mass-actions-panel', () => {
  const factory = generateShallowRenderer(ModalMassActionsPanel, { stubs });
  const snapshotFactory = generateRenderer(ModalMassActionsPanel, { stubs: snapshotStubs });

  test('Portal target receives correct name based on id prop', () => {
    const id = Faker.datatype.uuid();
    const wrapper = factory({
      propsData: {
        id,
      },
    });

    const portalTarget = selectPortalTarget(wrapper);
    const expectedName = `${PORTALS_NAMES.massActionsPanel}-${id}`;

    expect(portalTarget.attributes('name')).toBe(expectedName);
  });

  test('Portal target name updates when id prop changes', async () => {
    const id1 = Faker.datatype.uuid();
    const wrapper = factory({
      propsData: {
        id: id1,
      },
    });

    expect(selectPortalTarget(wrapper).attributes('name')).toBe(
      `${PORTALS_NAMES.massActionsPanel}-${id1}`,
    );

    const id2 = Faker.datatype.uuid();
    wrapper.setProps({ id: id2 });
    await wrapper.vm.$nextTick();

    expect(selectPortalTarget(wrapper).attributes('name')).toBe(
      `${PORTALS_NAMES.massActionsPanel}-${id2}`,
    );
  });

  test('Renders `modal-mass-actions-panel` with default props', () => {
    const id = 'modal-alarms-list';
    const wrapper = snapshotFactory({
      propsData: {
        id,
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper.text()).toContain(`Portal target name: ${PORTALS_NAMES.massActionsPanel}-${id}`);
  });

  test('Renders `modal-mass-actions-panel` with custom id', () => {
    const id = 'custom-modal-id';
    const wrapper = snapshotFactory({
      propsData: {
        id,
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper.text()).toContain(`Portal target name: ${PORTALS_NAMES.massActionsPanel}-${id}`);
  });
});

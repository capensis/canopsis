import { generateRenderer } from '@unit/utils/vue';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import PbehaviorsListExpandItem from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-list-expand-item.vue';

const stubs = {
  'pbehavior-patterns-form': true,
  'pbehavior-entities': true,
  'pbehavior-comments': true,
  'pbehavior-recurrence-rule': true,
};

const selectTabItems = wrapper => wrapper.findAll('.v-tab');
const selectTabItemByIndex = (wrapper, index) => selectTabItems(wrapper).at(index);

describe('pbehaviors-list-expand-item', () => {
  const snapshotFactory = generateRenderer(PbehaviorsListExpandItem, {

    stubs,
  });

  test('Renders `pbehaviors-list-expand-item` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          _id: 'pbehavior-id',
          comments: [{}, {}],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehaviors-list-expand-item` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          _id: 'pbehavior-id',
          rrule: 'pbehavior-rrule-id',
          entity_pattern: createEntityIdPatternByValue('entity-pattern'),
          comments: [{}, {}, {}],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `pbehaviors-list-expand-item` with opened entities tab', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          _id: 'pbehavior-id',
        },
      },
    });

    await selectTabItemByIndex(wrapper, 1).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehaviors-list-expand-item` with opened comments tab', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          _id: 'pbehavior-id',
          comments: [{}, {}, {}],
        },
      },
    });

    await selectTabItemByIndex(wrapper, 2).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehaviors-list-expand-item` with opened recurrence rule tab', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior: {
          _id: 'pbehavior-id',
          rrule: 'pbehavior-rrule-id',
        },
      },
    });

    await selectTabItemByIndex(wrapper, 3).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });
});

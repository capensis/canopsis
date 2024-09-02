import { flushPromises, generateRenderer } from '@unit/utils/vue';

import ExtraDetailsParents from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-parents.vue';

describe('extra-details-parents', () => {
  const total = 3;
  const rules = [
    {
      id: 'parent-rule-1-id',
      name: 'parent-rule-1-name',
    },
    {
      id: 'parent-rule-2-id',
      name: 'parent-rule-2-name',
    },
    {
      id: 'parent-rule-3-id',
      name: 'parent-rule-3-name',
    },
  ];

  const snapshotFactory = generateRenderer(ExtraDetailsParents, {

    attachTo: document.body,
  });

  it('Renders `extra-details-parents` with full parents', async () => {
    snapshotFactory({
      propsData: {
        total,
        rules,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-parents` without rules', async () => {
    snapshotFactory({
      propsData: {
        total,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});

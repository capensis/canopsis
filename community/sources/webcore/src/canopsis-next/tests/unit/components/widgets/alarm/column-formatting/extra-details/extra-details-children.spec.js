import { generateRenderer, flushPromises } from '@unit/utils/vue';

import ExtraDetailsChildren from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-children.vue';

describe('extra-details-children', () => {
  const total = 3;
  const opened = 2;
  const closed = 1;
  const rule = {
    name: 'rule-name',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsChildren, {

    attachTo: document.body,
  });

  it('Renders `extra-details-children` with full children and rule', async () => {
    snapshotFactory({
      propsData: {
        total,
        opened,
        closed,
        rule,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});

import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import flushPromises from 'flush-promises';
import ExtraDetailsChildren from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-children.vue';

const localVue = createVueInstance();

describe('extra-details-children', () => {
  const total = 3;
  const rule = {
    name: 'rule-name',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsChildren, {
    localVue,
    attachTo: document.body,
  });

  it('Renders `extra-details-children` with full children and rule', async () => {
    snapshotFactory({
      propsData: {
        total,
        rule,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});

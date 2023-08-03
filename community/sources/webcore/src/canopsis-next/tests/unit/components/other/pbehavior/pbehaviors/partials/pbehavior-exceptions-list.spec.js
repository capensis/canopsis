import { range } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import PbehaviorExceptionsList from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-exceptions-list.vue';

describe('pbehavior-exceptions-list', () => {
  const totalItems = 5;
  const exceptions = range(totalItems).map(index => ({
    name: `exception-${index}`,
    exdates: [
      {
        begin: 1000055522,
        end: 1000155522,
        type: {
          name: `exdate-${index}-1-name`,
        },
      },
      {
        begin: 1000256522,
        type: {
          name: `exdate-${index}-2-name`,
        },
      },
    ],
  }));

  const snapshotFactory = generateRenderer(PbehaviorExceptionsList);

  test('Renders `pbehavior-exceptions-list` without exceptions', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehavior-exceptions-list` with exceptions', () => {
    const wrapper = snapshotFactory({
      propsData: {
        exceptions,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});

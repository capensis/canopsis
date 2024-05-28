import { omit } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import ExtraDetailsPbehavior from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-pbehavior.vue';

const stubs = {
  'c-alarm-pbehavior-chip': true,
  'c-simple-tooltip': true,
};

describe('extra-details-pbehavior', () => {
  const pbehaviorInfo = {
    name: 'pbehavior-name',
    author: 'pbehavior-author',
    type_name: 'type-name',
    reason_name: 'pbehavior-reason-name',
    icon_name: 'icon-name',
    color: 'color',
    last_comment: {
      _id: 'pbehavior-comment-1-id',
      author: {
        display_name: 'pbehavior-comment-1-author',
      },
      message: 'pbehavior-comment-1-message',
    },
  };

  const snapshotFactory = generateRenderer(ExtraDetailsPbehavior, { stubs });

  it('Renders `extra-details-pbehavior` with full pbehavior', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviorInfo,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without reason', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviorInfo: omit(pbehaviorInfo, ['reason_name']),
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without comment', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviorInfo: omit(pbehaviorInfo, ['last_comment']),
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `extra-details-pbehavior` without comment author', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviorInfo: {
          ...pbehaviorInfo,
          last_comment: {
            author: null,
            ...pbehaviorInfo.last_comment,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

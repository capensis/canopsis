import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import AlarmHeaderTag from '@/components/widgets/alarm/headers-formatting/alarm-header-tag.vue';

const localVue = createVueInstance();

const stubs = {
  'c-alarm-action-chip': {
    template: '<span v-on="$listeners" class="c-alarm-action-chip"></span>',
  },
};

const factory = (options = {}) => shallowMount(AlarmHeaderTag, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmHeaderTag, {
  localVue,
  stubs,

  ...options,
});

const selectChip = wrapper => wrapper.find('.c-alarm-action-chip');

describe('alarm-header-tag', () => {
  const tags = [
    { value: 'tag1', color: 'color1' },
    { value: 'tag2', color: 'color2' },
    { value: 'tag3', color: 'color3' },
  ];
  const selectedTag = tags[0].value;
  const alarmTagModule = {
    name: 'alarmTag',
    getters: {
      items: () => tags,
    },
  };

  const store = createMockedStoreModules([alarmTagModule]);

  test('Should emit `clear` event', () => {
    const wrapper = factory({
      propsData: {
        selectedTag,
      },
      store,
    });

    const chip = selectChip(wrapper);

    chip.trigger('close');

    expect(wrapper).toEmit('clear');
  });

  it('Renders `alarm-header-tag` with selected tag', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selectedTag,
      },
      store,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-header-tag` without selected tag', () => {
    const wrapper = snapshotFactory({
      store,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-header-tag` with selected tag and default slot', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selectedTag,
      },
      slots: {
        default: 'Default text slot',
      },
      store,
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});

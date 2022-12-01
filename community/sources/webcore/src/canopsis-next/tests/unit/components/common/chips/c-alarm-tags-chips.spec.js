import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import AlarmColumnValueTags from '@/components/common/chips/c-alarm-tags-chips.vue';

const localVue = createVueInstance();

const stubs = {
  'c-alarm-tag-chip': {
    template: '<span v-on="$listeners" class="c-alarm-tag-chip"></span>',
  },
};

const factory = (options = {}) => shallowMount(AlarmColumnValueTags, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmColumnValueTags, {
  localVue,
  stubs,

  ...options,
});

const selectChip = wrapper => wrapper.find('.c-alarm-tag-chip');

describe('c-alarm-tags-chips', () => {
  const tags = [
    { value: 'tag1', color: 'color1' },
    { value: 'tag2', color: 'color2' },
    { value: 'tag3', color: 'color3' },
  ];
  const selectedTag = tags[2].value;
  const alarm = {
    tags: [tags[0].value],
  };
  const alarmTagModule = {
    name: 'alarmTag',
    getters: {
      items: () => tags,
    },
  };

  const store = createMockedStoreModules([alarmTagModule]);

  test('Should emit `select` event on `click`', () => {
    const wrapper = factory({
      propsData: {
        alarm,
      },
      store,
    });

    const chip = selectChip(wrapper);

    chip.trigger('click');

    expect(wrapper).toEmit('select', tags[0].value);
  });

  it('Renders `c-alarm-tags-chips` without selected tag and dropdown', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm,
      },
      store,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-tags-chips` without selected tag and with dropdown', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          tags: tags.map(({ value }) => value),
        },
      },
      store,
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-tags-chips` with selected tag and dropdown', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          tags: tags.map(({ value }) => value),
        },
        selectedTag,
      },
      store,
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });
});

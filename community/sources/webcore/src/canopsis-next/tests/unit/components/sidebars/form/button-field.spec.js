import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';

import ButtonField from '@/components/sidebars/form/fields/button-field.vue';

const stubs = {
  'v-btn': createButtonStub('v-btn'),
};

const selectCreateButton = wrapper => wrapper.findAll('button.v-btn').at(0);
const selectEditButton = wrapper => wrapper.findAll('button.v-btn').at(0);
const selectDeleteButton = wrapper => wrapper.findAll('button.v-btn').at(1);

describe('button-field', () => {
  const factory = generateShallowRenderer(ButtonField, { stubs });
  const snapshotFactory = generateRenderer(ButtonField);

  it('Create event emitted after click on the button', () => {
    const wrapper = factory({
      propsData: {
        title: 'Title',
        addable: true,
      },
    });

    selectCreateButton(wrapper).trigger('click');

    expect(wrapper).toEmit('create', expect.any(Event));
  });

  it('Edit event emitted after click on the button', () => {
    const wrapper = factory({
      propsData: {
        title: 'Title',
        isEmpty: false,
      },
    });

    selectEditButton(wrapper).trigger('click');

    expect(wrapper).toEmit('edit', expect.any(Event));
  });

  it('Delete event emitted after click on the button', async () => {
    const wrapper = factory({
      propsData: {
        title: 'Title',
        isEmpty: false,
        removable: true,
      },
    });

    await selectDeleteButton(wrapper).trigger('click');

    expect(wrapper).toEmit('delete', expect.any(Event));
  });

  it('Renders `button-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `button-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Title',
        isEmpty: true,
        addable: true,
        removable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

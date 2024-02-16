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
        addable: true,
      },
    });

    const createButton = selectCreateButton(wrapper);

    createButton.trigger('click');

    const createEvents = wrapper.emitted('create');

    expect(createEvents).toHaveLength(1);

    const [eventData] = createEvents[0];
    expect(eventData).toEqual(expect.any(Event));
  });

  it('Edit event emitted after click on the button', () => {
    const wrapper = factory({
      propsData: {
        isEmpty: false,
      },
    });

    const editButton = selectEditButton(wrapper);

    editButton.trigger('click');

    const editEvents = wrapper.emitted('edit');

    expect(editEvents).toHaveLength(1);

    const [eventData] = editEvents[0];
    expect(eventData).toEqual(expect.any(Event));
  });

  it('Delete event emitted after click on the button', () => {
    const wrapper = factory({
      propsData: {
        isEmpty: false,
        removable: true,
      },
    });

    const deleteButton = selectDeleteButton(wrapper);

    deleteButton.trigger('click');

    const deleteEvents = wrapper.emitted('delete');

    expect(deleteEvents).toHaveLength(1);

    const [eventData] = deleteEvents[0];
    expect(eventData).toEqual(expect.any(Event));
  });

  it('Renders `button-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: '',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});

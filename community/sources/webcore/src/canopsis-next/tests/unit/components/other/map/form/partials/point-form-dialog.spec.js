import Faker from 'faker';
import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { mockModals } from '@unit/utils/mock-hooks';
import PointFormDialog from '@/components/other/map/form/partials/point-form-dialog.vue';
import { MODALS } from '@/constants';

const localVue = createVueInstance();

const stubs = {
  'point-form': true,
  'v-form': createFormStub('v-form'),
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'point-form': true,
};

const factory = (options = {}) => shallowMount(PointFormDialog, {
  localVue,
  stubs,
  attachTo: document.body,

  ...options,
});

const snapshotFactory = (options = {}) => mount(PointFormDialog, {
  localVue,
  stubs: snapshotStubs,
  attachTo: document.body,

  ...options,
});

const selectPointForm = wrapper => wrapper.find('point-form-stub');
const selectCloseButton = wrapper => wrapper.findAll('button').at(0);
const selectCancelButton = wrapper => wrapper.findAll('button').at(1);
const selectDeleteButton = wrapper => wrapper.find('button.error');
const selectSubmitButton = wrapper => wrapper.find('button[type="submit"]');

describe('point-form-dialog', () => {
  const $modals = mockModals();

  test('Point dialog submitted with changes', async () => {
    const point = {
      _id: 'id',
      entity: 'entity',
    };
    const wrapper = factory({
      propsData: {
        point,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(wrapper).toEmit('submit', point);
  });

  test('Point dialog submitted without changes', async () => {
    const wrapper = factory({
      propsData: {
        point: {},
      },
    });

    const newPoint = {
      _id: 'id',
      entity: 'entity',
    };

    const pointForm = selectPointForm(wrapper);

    pointForm.vm.$emit('input', newPoint);

    const submitButton = selectSubmitButton(wrapper);
    submitButton.trigger('click');

    await flushPromises();

    expect(wrapper).toEmit('submit', newPoint);
  });

  test('Point dialog submitted with valid data after update prop', async () => {
    const point = {
      _id: 'id',
      entity: 'entity',
    };
    const wrapper = factory({
      propsData: {
        point,
      },
    });

    const newPoint = {
      ...point,
      map: 'map',
    };

    await wrapper.setProps({
      point: newPoint,
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(wrapper).toEmit('submit', newPoint);
  });

  test('Point dialog didn\'t submitted with errors', async () => {
    const point = {
      _id: 'id',
      entity: 'entity',
    };
    const wrapper = factory({
      propsData: {
        point,
      },
    });

    const validator = wrapper.getValidator();

    const pointForm = selectPointForm(wrapper);

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => pointForm.vm,
      vm: pointForm.vm,
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(wrapper).not.toEmit('submit');
  });

  test('Point dialog closed after trigger cancel button', () => {
    const wrapper = factory({
      propsData: {
        point: {},
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    expect(wrapper).toEmit('cancel');
  });

  test('Point dialog closed after trigger close button', () => {
    const wrapper = factory({
      propsData: {
        point: {},
      },
    });

    const closeButton = selectCloseButton(wrapper);

    closeButton.trigger('click');

    expect(wrapper).toEmit('cancel');
  });

  test('Point removed after trigger delete button', () => {
    const wrapper = factory({
      propsData: {
        point: {},
        editing: true,
      },
    });

    const deleteButton = selectDeleteButton(wrapper);

    deleteButton.trigger('click');

    expect(wrapper).toEmit('remove');
  });

  test('Cancel emitted after click outside and confirm close', async () => {
    const pointId = Faker.datatype.number();

    const wrapper = snapshotFactory({
      propsData: {
        point: {
          _id: pointId,
        },
        editing: true,
      },
      mocks: {
        $modals,
      },
    });

    wrapper.clickOutside();

    expect($modals.show).toBeCalledWith(
      {
        id: pointId,
        name: MODALS.clickOutsideConfirmation,
        dialogProps: {
          persistent: true,
        },
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(false);

    expect(wrapper).toEmit('cancel');
  });

  test('Submit emitted after click outside and confirm save', async () => {
    const pointId = Faker.datatype.number();
    const point = {
      _id: pointId,
    };

    const wrapper = snapshotFactory({
      propsData: {
        point,
        editing: true,
      },
      mocks: {
        $modals,
      },
    });

    wrapper.clickOutside();

    expect($modals.show).toBeCalled();

    const [modalArguments] = $modals.show.mock.calls[0];

    await modalArguments.config.action(true);

    expect(wrapper).toEmit('submit', point);
  });

  test('Renders `point-form-dialog` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        point: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `point-form-dialog` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        point: {},
        editing: true,
        coordinates: true,
        existsEntities: [{}],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});

import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { MODALS } from '@/constants';

import PointFormDialog from '@/components/other/map/form/fields/point-form-dialog.vue';

const stubs = {
  'point-form': true,
  'v-form': createFormStub('v-form'),
  'v-btn': createButtonStub('v-btn'),
};

describe('point-form-dialog', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(PointFormDialog, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  beforeAll(() => vi.useFakeTimers());

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

    await wrapper.vm.submit();

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

    wrapper.vm.updateForm(newPoint);
    await wrapper.vm.submit();

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

    await wrapper.vm.submit();

    expect(wrapper).toEmit('submit', newPoint);
  });

  test('Point dialog didn\'t submitted with incomplete data', async () => {
    const wrapper = factory({
      propsData: {
        point: {
          _id: 'id',
        },
      },
    });

    await wrapper.vm.submit();

    expect(wrapper).not.toHaveBeenEmit('submit');
  });

  test('Point dialog closed after trigger cancel button', () => {
    const wrapper = factory({
      propsData: {
        point: {},
      },
    });

    wrapper.vm.close();

    expect(wrapper).toHaveBeenEmit('cancel');
  });

  test('Point dialog closed after trigger close button', () => {
    const wrapper = factory({
      propsData: {
        point: {},
      },
    });

    wrapper.vm.close();

    expect(wrapper).toHaveBeenEmit('cancel');
  });

  test('Point removed after trigger delete button', () => {
    const wrapper = factory({
      propsData: {
        point: {},
        editing: true,
      },
    });

    wrapper.vm.emit('remove');

    expect(wrapper).toHaveBeenEmit('remove');
  });

  test('Cancel emitted after click outside and confirm close', async () => {
    const pointId = Faker.datatype.number();
    const point = {
      _id: pointId,
      entity: undefined,
    };

    const wrapper = factory({
      propsData: {
        point,
        editing: true,
      },
      mocks: {
        $modals,
      },
    });

    const newPoint = {
      ...point,
      entity: Faker.datatype.string(),
    };

    wrapper.vm.updateForm(newPoint);

    wrapper.vm.clickOutsideDirective.handler();

    expect($modals.show).toHaveBeenCalledWith({
      id: pointId,
      name: MODALS.clickOutsideConfirmation,
      dialogProps: {
        persistent: true,
      },
      config: {
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(false);

    expect(wrapper).toHaveBeenEmit('cancel');
  });

  test('Submit emitted after click outside and confirm save', async () => {
    const pointId = Faker.datatype.number();
    const point = {
      _id: pointId,
    };

    const wrapper = factory({
      propsData: {
        point,
        editing: true,
      },
      mocks: {
        $modals,
      },
    });

    const newPoint = {
      ...point,
      entity: Faker.datatype.string(),
    };

    wrapper.vm.updateForm(newPoint);

    wrapper.vm.clickOutsideDirective.handler();

    expect($modals.show).toBeCalled();

    const [modalArguments] = $modals.show.mock.calls[0];

    await modalArguments.config.action(true);

    await flushPromises();

    expect(wrapper).toEmit('submit', newPoint);
  });
});

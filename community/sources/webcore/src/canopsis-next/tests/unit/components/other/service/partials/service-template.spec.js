import flushPromises from 'flush-promises';
import Faker from 'faker';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { WEATHER_ACTIONS_TYPES } from '@/constants';

import ServiceTemplate from '@/components/other/service/partials/service-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';

const localVue = createVueInstance();

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'service-entities-list': true,
};

const selectEntitiesList = wrapper => wrapper.find('service-entities-list-stub');

describe('service-template', () => {
  const service = {
    _id: 'service-id',
  };
  const modalTemplate = `
    Id: {{ entity._id }}
    Entities: {{ entities name="entity._id" }}
    <br />
  `;

  const { authModule } = createAuthModule();

  const store = createMockedStoreModules([authModule]);

  const snapshotFactory = generateRenderer(ServiceTemplate, {
    store,
    localVue,
    stubs,
    propsData: {
      service,
      pagination: {},
    },
  });

  test('Action applied after triggers entities list', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate,
        },
      },
    });

    await flushPromises();

    const entitiesList = selectEntitiesList(wrapper);

    const action = {
      entities: [{}],
      payload: { output: Faker.datatype.string() },
      actionType: WEATHER_ACTIONS_TYPES.entityAck,
    };

    await entitiesList.vm.$emit('apply:action', action);

    expect(wrapper).toEmit('apply:action', action);
  });

  test('Refresh applied after triggers entities list', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate,
        },
      },
    });

    await flushPromises();

    const entitiesList = selectEntitiesList(wrapper);

    await entitiesList.vm.$emit('refresh');

    expect(wrapper).toEmit('refresh');
  });

  test('Pagination updated after triggers entities list', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate,
        },
      },
    });

    await flushPromises();

    const entitiesList = selectEntitiesList(wrapper);

    const newPagination = {
      page: Faker.datatype.number(),
    };

    await entitiesList.vm.$emit('update:pagination', newPagination);

    expect(wrapper).toEmit('update:pagination', newPagination);
  });

  test('Renders `service-template` with required props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-template` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate,
        },
        itemsPerPage: 10,
        totalItems: 20,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});

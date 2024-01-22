import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';

import ServiceTemplate from '@/components/other/service/partials/service-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
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

    stubs,
    propsData: {
      service,
      options: {},
    },
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

    await entitiesList.triggerCustomEvent('refresh');

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

    await entitiesList.triggerCustomEvent('update:options', newPagination);

    expect(wrapper).toEmit('update:options', newPagination);
  });

  test('Renders `service-template` with required props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
